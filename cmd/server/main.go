package main

import (
	"context"
	"database/sql"
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
	"github.com/SigmarWater/crm/internal/migrator"
	clientRepository "github.com/SigmarWater/crm/internal/repository/client"
	clientService "github.com/SigmarWater/crm/internal/service/client"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load .env: %v\n", err)
		return
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Println("DB_URI is not set")
		return
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	sqlDB, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Printf("failed to open database for migrations: %v\n", err)
		return
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		log.Println("MIGRATIONS_DIR is not set")
		return
	}

	migratorRunner := migrator.NewMigrator(sqlDB, migrationsDir)
	if err := migratorRunner.Up(); err != nil {
		log.Printf("failed to run migrations: %v\n", err)
		return
	}

	repo := clientRepository.NewRepository(pool)
	service := clientService.NewClientService(repo)
	api := clientV1API.NewAPI(service)

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

	crmV1.RegisterCRMServiceServer(server, api)

	reflection.Register(server)

	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)

		if err := server.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

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
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := gwServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Println("HTTP server stopped")

	// В конце останавливаем gRPC сервер
	server.GracefulStop()

	log.Println("gRPC server stopped")
}
