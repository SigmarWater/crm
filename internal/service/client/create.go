package client

import (
	"context"

	serviceModel "github.com/SigmarWater/crm/internal/model"
	"github.com/SigmarWater/crm/internal/service/converter"
)

func (s *service) Create(ctx context.Context, info *serviceModel.CreateClientInfo) (*serviceModel.Client, error) {
	client, err := s.clientRepository.Create(ctx, converter.CreateClientInfoToRepo(info))
	if err != nil {
		return nil, err
	}

	return converter.ClientFromRepo(client), nil
}
