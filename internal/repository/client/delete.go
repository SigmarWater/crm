package client

import (
	"context"

	repoError "github.com/SigmarWater/crm/internal/model"
)

func (r *repository) Delete(_ context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[uuid]; !ok {
		return repoError.ErrClientNotFound
	}

	delete(r.storage, uuid)

	return nil
}
