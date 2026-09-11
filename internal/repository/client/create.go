package client

import (
	"context"
	"time"

	repoModel "github.com/SigmarWater/crm/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) Create(_ context.Context, info *repoModel.ClientInfo) (*repoModel.Client, error) {
	clientUUID := uuid.NewString()
	now := time.Now()

	client := repoModel.Client{
		UUID:      clientUUID,
		Name:      info.Name,
		Phone:     info.Phone,
		Email:     info.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[clientUUID] = client

	return &client, nil
}
