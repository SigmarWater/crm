package client

import (
	"context"

	"github.com/SigmarWater/crm/internal/errors"
)

func (r *repository) Delete(_ context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[uuid]; !ok {
		return errors.ErrClientNotFound
	}

	delete(r.storage, uuid)

	return nil
}
