package model

import "time"

type Client struct {
	UUID      string
	Name      string
	Phone     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClientInfo struct {
	Name  string
	Phone string
	Email string
}

type UpdateClientInfo struct {
	Name  *string
	Phone *string
	Email *string
}
