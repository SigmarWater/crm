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

func (a *api) UpdateCustomer(ctx context.Context, req *customerV1.UpdateCustomerRequest) (*customerV1.UpdateCustomerResponse, error) {
	client, err := a.customerService.Update(ctx, req.GetUuid(), converter.UpdateCustomerInfoFromProto(req))
	if err != nil {
		if errors.Is(err, appError.ErrCustomerNotFound) {
			return nil, status.Errorf(codes.NotFound, "customer with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &customerV1.UpdateCustomerResponse{
		Customer: converter.CustomerToProto(client),
	}, nil
}
