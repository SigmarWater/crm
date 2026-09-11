package client

import (
	"context"

	"github.com/SigmarWater/crm/internal/errors"
	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

func (r *repository) Get(_ context.Context, uuid string) (*repoModel.Client, error) {
	r.mu.RLock()
	client, ok := r.storage[uuid]
	r.mu.RUnlock()

	if !ok {
		return nil, errors.ErrClientNotFound
	}

	return &repoModel.Client{
		UUID:      client.UUID,
		Name:      client.Name,
		Phone:     client.Phone,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}
