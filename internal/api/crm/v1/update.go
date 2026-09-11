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

func (a *api) UpdateClient(ctx context.Context, req *crmV1.UpdateClientRequest) (*crmV1.UpdateClientResponse, error) {
	client, err := a.clientService.Update(ctx, req.GetUuid(), converter.UpdateClientInfoFromProto(req))
	if err != nil {
		if errors.Is(err, appError.ErrClientNotFound) {
			return nil, status.Errorf(codes.NotFound, "client with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &crmV1.UpdateClientResponse{
		Client: converter.ClientToProto(client),
	}, nil
}
