package client

import (
	"sync"

	repo "github.com/SigmarWater/crm/internal/repository"
	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

var _ repo.ClientRepository = (*repository)(nil)

type repository struct {
	storage map[string]repoModel.Client
	mu      sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		storage: make(map[string]repoModel.Client),
	}
}
