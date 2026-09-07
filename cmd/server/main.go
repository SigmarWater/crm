package main

import (
	"context"
	"fmt"
	"github.com/SigmarWater/crm/internal/interceptor"
	crmV1 "github.com/SigmarWater/crm/pkg/api/crm_service"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const grpcPort = 50051

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
	c.mu.RUnlock()
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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
	log.Println("Server stopped")
}
