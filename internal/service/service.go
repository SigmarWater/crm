package service

import (
	"context"

	"github.com/SigmarWater/crm/internal/model"
)

type ClientService interface {
	Create(ctx context.Context, info *model.CreateClientInfo) (*model.Client, error)
	Get(ctx context.Context, uuid string) (*model.Client, error)
	Update(ctx context.Context, uuid string, updateInfo *model.UpdateClientInfo) (*model.Client, error)
	Delete(ctx context.Context, uuid string) error
}
