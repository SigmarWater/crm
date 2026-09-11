package client

import (
	"context"

	serviceModel "github.com/SigmarWater/crm/internal/model"
	"github.com/SigmarWater/crm/internal/service/converter"
)

func (s *service) Get(ctx context.Context, uuid string) (*serviceModel.Client, error) {
	client, err := s.clientRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return converter.ClientFromRepo(client), nil
}
