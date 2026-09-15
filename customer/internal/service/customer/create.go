package customer

import (
	"context"

	serviceModel "github.com/SigmarWater/crm/customer/internal/model"
	"github.com/SigmarWater/crm/customer/internal/service/converter"
)

func (s *service) Create(ctx context.Context, info *serviceModel.CreateCustomerInfo) (*serviceModel.Customer, error) {
	client, err := s.customerRepository.Create(ctx, converter.CreateCustomerInfoToRepo(info))
	if err != nil {
		return nil, err
	}

	return converter.CustomerFromRepo(client), nil
}
