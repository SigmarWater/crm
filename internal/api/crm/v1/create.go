package v1

import (
	"context"

	"github.com/SigmarWater/crm/internal/converter"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
)

func (a *api) CreateClient(ctx context.Context, req *crmV1.CreateClientRequest) (*crmV1.CreateClientResponse, error) {
	client, err := a.clientService.Create(ctx, converter.CreateClientInfoFromProto(req))
	if err != nil {
		return nil, err
	}

	return &crmV1.CreateClientResponse{
		Client: converter.ClientToProto(client),
	}, nil
}
