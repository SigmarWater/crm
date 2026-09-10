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
	"sync"
	"syscall"
	"time"

	"github.com/SigmarWater/crm/internal/interceptor"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
	uuid2 "github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

type crmService struct {
	crmV1.UnimplementedCRMServiceServer
	mu      sync.RWMutex
	clients map[string]*crmV1.Client
}

// Добавление клиента
func (c *crmService) CreateClient(_ context.Context, req *crmV1.CreateClientRequest) (*crmV1.CreateClientResponse, error) {
	newUUID := uuid2.NewString()

	client := &crmV1.Client{
		Uuid:  newUUID,
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Email: req.GetEmail(),
	}

	c.mu.Lock()
	c.clients[newUUID] = client
	c.mu.Unlock()

	log.Printf("Создан клиент с UUID %s\n", newUUID)

	return &crmV1.CreateClientResponse{
		Client: client,
	}, nil
}

// Получение информации о клиенте
func (c *crmService) GetClient(_ context.Context, req *crmV1.GetClientRequest) (*crmV1.GetClientResponse, error) {
	uuid := req.GetUuid()

	c.mu.RLock()
	defer c.mu.RUnlock()
	clientInfo, ok := c.clients[uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "client with UUID %s not found", req.GetUuid())
	}

	return &crmV1.GetClientResponse{
		Client: clientInfo,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
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

	service := &crmService{
		clients: make(map[string]*crmV1.Client),
	}

	crmV1.RegisterCRMServiceServer(server, service)

	reflection.Register(server)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)

		if err := server.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для HTTP запросов
		mux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC сервером
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err := crmV1.RegisterCRMServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
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
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Запускаем HTTP сервер
		log.Printf("HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")

	// Сначала аккуратно останавливаем HTTP сервер
	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("HTTP server stopped")
	}

	// В конце останавливаем gRPC сервер
	server.GracefulStop()
	log.Println("gRPC server stopped")
}
