package v1

import (
	"context"

	"github.com/SigmarWater/crm/customer/internal/converter"
	customerV1 "github.com/SigmarWater/crm/shared/pkg/customer_service/v1"
)

func (a *api) CreateCustomer(ctx context.Context, req *customerV1.CreateCustomerRequest) (*customerV1.CreateCustomerResponse, error) {
	customer, err := a.customerService.Create(ctx, converter.CreateCustomerInfoFromProto(req))
	if err != nil {
		return nil, err
	}

	return &customerV1.CreateCustomerResponse{
		Customer: converter.CustomerToProto(customer),
	}, nil
}
