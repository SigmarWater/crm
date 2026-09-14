package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/SigmarWater/crm/jwt/internal/api"
	"github.com/SigmarWater/crm/jwt/internal/service"
	jwtV1 "github.com/SigmarWater/crm/jwt/pkg/jwt/v1"
)

const (
	grpcPort = "localhost:50051"
)

func main() {
	// Создаем JWT сервис
	jwtService := service.NewJWTService()

	// Создаем gRPC хендлер
	jwtHandler := api.NewJWTHandler(jwtService)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис
	jwtV1.RegisterJWTServiceServer(grpcServer, jwtHandler)

	// Включаем reflection для удобства тестирования
	reflection.Register(grpcServer)

	// Создаем listener
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Printf("Failed to listen on port %s: %v\n", grpcPort, err)
		return
	}

	fmt.Printf("🚀 JWT gRPC server listening on %s\n", grpcPort)
	fmt.Println("📋 Available users:")
	fmt.Println("  - admin:admin123")
	fmt.Println("  - user1:password1")
	fmt.Println("  - user2:password2")
	fmt.Println("  - john:john123")
	fmt.Println("  - alice:alice456")

	// Запускаем сервер
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Printf("Failed to serve gRPC server: %v\n", err)
	}
}
