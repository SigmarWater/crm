package model

import "errors"

var (
	ErrClientNotFound = errors.New("client not found")
	ErrClientExists   = errors.New("client already exists")
)
