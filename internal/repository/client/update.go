package client

import (
	"context"
	"time"

	repoError "github.com/SigmarWater/crm/internal/model"
	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

func (r *repository) Update(_ context.Context, uuid string, updateInfo *repoModel.UpdateClientInfo) (*repoModel.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.storage[uuid]

	if !ok {
		return nil, repoError.ErrClientNotFound
	}

	if name := updateInfo.Name; name != nil {
		client.Name = *name
	}

	if phone := updateInfo.Phone; phone != nil {
		client.Phone = *phone
	}

	if email := updateInfo.Email; email != nil {
		client.Email = *email
	}

	client.UpdatedAt = time.Now()

	r.storage[uuid] = client

	return &client, nil
}
