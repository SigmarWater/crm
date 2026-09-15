package service

import (
	"context"

	"github.com/SigmarWater/crm/customer/internal/model"
)

type CustomerService interface {
	Create(ctx context.Context, info *model.CreateCustomerInfo) (*model.Customer, error)
	Get(ctx context.Context, uuid string) (*model.Customer, error)
	Update(ctx context.Context, uuid string, updateInfo *model.UpdateCustomerInfo) (*model.Customer, error)
	Delete(ctx context.Context, uuid string) error
}
