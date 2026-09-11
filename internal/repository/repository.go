package repository

import (
	"context"

	"github.com/SigmarWater/crm/internal/repository/model"
)

type ClientRepository interface {
	Create(ctx context.Context, info *model.ClientInfo) (*model.Client, error)
	Get(ctx context.Context, uuid string) (*model.Client, error)
	Update(ctx context.Context, uuid string, updateInfo *model.UpdateClientInfo) (*model.Client, error)
	Delete(ctx context.Context, uuid string) error
}
