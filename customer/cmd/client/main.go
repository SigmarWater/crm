package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customer_service "github.com/SigmarWater/crm/shared/pkg/customer_service/v1"
)

const serverAddress = "localhost:50051"

func createCustomer(ctx context.Context, client customer_service.CustomerServiceClient) (*customer_service.Customer, error) {
	customerInfo := &customer_service.CreateCustomerRequest{
		Name:  gofakeit.Name(),
		Phone: fmt.Sprintf("+%s", gofakeit.Phone()),
		Email: gofakeit.Email(),
	}

	response, err := client.CreateCustomer(ctx, customerInfo)
	if err != nil {
		return &customer_service.Customer{}, err
	}

	return response.Customer, nil
}

func getCustomer(ctx context.Context, client customer_service.CustomerServiceClient, uuid string) (*customer_service.GetCustomerResponse, error) {
	request := &customer_service.GetCustomerRequest{
		Uuid: uuid,
	}

	response, err := client.GetCustomer(ctx, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func updateCustomer(
	ctx context.Context,
	client customer_service.CustomerServiceClient,
	uuid string,
) (*customer_service.Customer, error) {
	name := gofakeit.Name()
	phone := fmt.Sprintf("+%s", gofakeit.Phone())
	email := gofakeit.Email()

	request := &customer_service.UpdateCustomerRequest{
		Uuid:  uuid,
		Name:  &name,
		Phone: &phone,
		Email: &email,
	}

	response, err := client.UpdateCustomer(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Customer, nil
}

func deleteCustomer(
	ctx context.Context,
	client customer_service.CustomerServiceClient,
	uuid string,
) error {
	request := &customer_service.DeleteCustomerRequest{
		Uuid: uuid,
	}

	_, err := client.DeleteCustomer(ctx, request)
	return err
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

	client := customer_service.NewCustomerServiceClient(conn)

	log.Println("=== Тестирование API для работы с покупателями===")
	log.Println()

	// 1. Создаем покупателя
	log.Println("Создание покупателя")
	log.Println("===========================")
	customerInfo, err := createCustomer(ctx, client)
	if err != nil {
		log.Printf("Ошибка при создании покупателя: %v\n", err)
		return
	}

	// Выводим информацию о созданном покупателе
	log.Printf("Создан покупателя: UUID=%s\n", customerInfo.Uuid)

	// 2. Получаем информацию о покупателе
	log.Println("Получение информации о покупателе")
	log.Println("==================================")
	customerInfoResp, err := getCustomer(ctx, client, customerInfo.Uuid)
	if err != nil {
		log.Printf("Ошибка при получении информации о покупателе: %v\n", err)
		return
	}

	// Выводим информацию о полученном покупателе
	log.Printf("Получен покупатель: UUID=%s", customerInfoResp.Customer.Uuid)
	log.Printf("%v\n", customerInfoResp)

	// 3. Обновляем информацию о покупателе
	log.Println("Обновление информации о покупателе")
	log.Println("==================================")

	updatedCustomer, err := updateCustomer(ctx, client, customerInfo.Uuid)
	if err != nil {
		log.Printf("Ошибка при обновлении покупателя: %v\n", err)
		return
	}

	log.Printf("Обновлён покупатель: UUID=%s", updatedCustomer.Uuid)
	log.Printf("%v\n", updatedCustomer)

	// 4. Удаляем клиента
	log.Println("Удаление покупателя")
	log.Println("==================================")

	err = deleteCustomer(ctx, client, customerInfo.Uuid)
	if err != nil {
		log.Printf("Ошибка при удалении покупателя: %v\n", err)
		return
	}

	log.Printf("Клиент удалён: UUID=%s\n", customerInfo.Uuid)

	log.Println("Тестирование завершено!")
}
