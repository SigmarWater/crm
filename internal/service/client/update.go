package client

import (
	"context"

	serviceModel "github.com/SigmarWater/crm/internal/model"
	"github.com/SigmarWater/crm/internal/service/converter"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo *serviceModel.UpdateClientInfo) (*serviceModel.Client, error) {
	client, err := s.clientRepository.Update(ctx, uuid, converter.UpdateClientInfoToRepo(updateInfo))
	if err != nil {
		return nil, err
	}
	return converter.ClientFromRepo(client), nil
}
