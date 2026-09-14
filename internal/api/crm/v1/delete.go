package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	appError "github.com/SigmarWater/crm/internal/errors"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
)

func (a *api) DeleteClient(ctx context.Context, req *crmV1.DeleteClientRequest) (*emptypb.Empty, error) {
	err := a.clientService.Delete(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, appError.ErrClientNotFound) {
			return nil, status.Errorf(codes.NotFound, "client with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
