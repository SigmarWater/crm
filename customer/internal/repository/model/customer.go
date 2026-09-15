package model

import "time"

type Customer struct {
	UUID      string
	Name      string
	Phone     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CustomerInfo struct {
	Name  string
	Phone string
	Email string
}

type UpdateCustomerInfo struct {
	Name  *string
	Phone *string
	Email *string
}
