package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SigmarWater/crm/customer/internal/converter"
	appError "github.com/SigmarWater/crm/customer/internal/errors"
	customerV1 "github.com/SigmarWater/crm/shared/pkg/customer_service/v1"
)

func (a *api) GetCustomer(ctx context.Context, req *customerV1.GetCustomerRequest) (*customerV1.GetCustomerResponse, error) {
	customer, err := a.customerService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, appError.ErrCustomerNotFound) {
			return nil, status.Errorf(codes.NotFound, "customer with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &customerV1.GetCustomerResponse{
		Customer: converter.CustomerToProto(customer),
	}, nil
}
