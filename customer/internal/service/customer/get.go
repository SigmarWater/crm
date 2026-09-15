package customer

import (
	"context"

	serviceModel "github.com/SigmarWater/crm/customer/internal/model"
	"github.com/SigmarWater/crm/customer/internal/service/converter"
)

func (s *service) Get(ctx context.Context, uuid string) (*serviceModel.Customer, error) {
	client, err := s.customerRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return converter.CustomerFromRepo(client), nil
}
