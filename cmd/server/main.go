package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	clientV1API "github.com/SigmarWater/crm/internal/api/crm/v1"
	"github.com/SigmarWater/crm/internal/interceptor"
	clientRepository "github.com/SigmarWater/crm/internal/repository/client"
	clientService "github.com/SigmarWater/crm/internal/service/client"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if closeErr := lis.Close(); closeErr != nil {
			log.Printf("failed to close listener: %v\n", closeErr)
		}
	}()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.ValidatorInterceptor(),
		),
	)

	repo := clientRepository.NewRepository()
	service := clientService.NewClientService(repo)
	api := clientV1API.NewAPI(service)

	crmV1.RegisterCRMServiceServer(server, api)

	reflection.Register(server)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)

		if err := server.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			return
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем мультиплексор для HTTP запросов
	mux := runtime.NewServeMux()

	// Настраиваем опции для соединения с gRPC сервером
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем gRPC-gateway хендлеры
	if err := crmV1.RegisterCRMServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", grpcPort),
		opts,
	); err != nil {
		log.Printf("failed to register gateway: %v\n", err)
		return
	}

	// OpenAPI JSON
	fileServer := http.FileServer(http.Dir("docs/openapi"))

	// Создаем HTTP маршрутизатор
	httpMux := http.NewServeMux()

	// Регистрируем API эндпоинты
	httpMux.Handle("/", mux)

	// Swagger UI эндпоинты
	// Swagger UI: /swagger/ → docs/openapi/
	httpMux.Handle("/swagger/",
		http.StripPrefix("/swagger/", fileServer),
	)

	httpMux.Handle("/crm.swagger.json", fileServer)

	// Создаем HTTP сервер
	gwServer := http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           httpMux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		// Запускаем HTTP сервер
		log.Printf("HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err := gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Starting graceful shutdown...")

	// Сначала аккуратно останавливаем HTTP сервер
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := gwServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Println("HTTP server stopped")

	// В конце останавливаем gRPC сервер
	server.GracefulStop()

	log.Println("gRPC server stopped")
}
