package main

import (
	"context"
	"fmt"
	"log"

	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddress = "localhost:50051"

func createClient(ctx context.Context, client crmV1.CRMServiceClient) (*crmV1.Client, error) {
	clientInfo := &crmV1.CreateClientRequest{
		Name:  gofakeit.Name(),
		Phone: fmt.Sprintf("+%s", gofakeit.Phone()),
		Email: gofakeit.Email(),
	}

	response, err := client.CreateClient(ctx, clientInfo)
	if err != nil {
		return &crmV1.Client{}, err
	}

	return response.Client, nil
}

func getClient(ctx context.Context, client crmV1.CRMServiceClient, uuid string) (*crmV1.GetClientResponse, error) {
	request := &crmV1.GetClientRequest{
		Uuid: uuid,
	}

	response, err := client.GetClient(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func main() {
	ctx := context.Background()

	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			log.Printf("failed to close connect: %v", closeErr)
		}
	}()

	client := crmV1.NewCRMServiceClient(conn)

	log.Println("=== Тестирование API для работы с клиентами===")
	log.Println()

	// 1. Создаем несколько наблюдений
	log.Println("Создание клиентов")
	log.Println("===========================")
	clientInfo, err := createClient(ctx, client)
	if err != nil {
		log.Printf("Ошибка при создании клиента: %v\n", err)
		return
	}

	// Выводим информацию о созданном наблюдении
	log.Printf("Создан клиент: UUID=%s\n", clientInfo.Uuid)

	// 2. Получаем информацию о наблюдении
	log.Println("Получение информации о клиента")
	log.Println("==================================")
	clientInfoResp, err := getClient(ctx, client, clientInfo.Uuid)
	if err != nil {
		log.Printf("Ошибка при получении информации о клиенте: %v\n", err)
		return
	}

	// Выводим информацию о полученном наблюдении
	log.Printf("Получен клиент: UUID=%s", clientInfoResp.Client.Uuid)
	log.Printf("%v\n", clientInfoResp)

	log.Println("Тестирование завершено!")
}
