package v1

import (
	"context"
	"errors"

	"github.com/SigmarWater/crm/internal/converter"
	appError "github.com/SigmarWater/crm/internal/errors"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetClient(ctx context.Context, req *crmV1.GetClientRequest) (*crmV1.GetClientResponse, error) {
	client, err := a.clientService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, appError.ErrClientNotFound) {
			return nil, status.Errorf(codes.NotFound, "client with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &crmV1.GetClientResponse{
		Client: converter.ClientToProto(client),
	}, nil
}
