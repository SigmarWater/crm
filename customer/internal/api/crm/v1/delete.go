package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	appError "github.com/SigmarWater/crm/customer/internal/errors"
	customerV1 "github.com/SigmarWater/crm/shared/pkg/customer_service/v1"
)

func (a *api) DeleteCustomer(ctx context.Context, req *customerV1.DeleteCustomerRequest) (*emptypb.Empty, error) {
	err := a.customerService.Delete(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, appError.ErrCustomerNotFound) {
			return nil, status.Errorf(codes.NotFound, "customer with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
