package customer

import (
	"context"

	serviceModel "github.com/SigmarWater/crm/customer/internal/model"
	"github.com/SigmarWater/crm/customer/internal/service/converter"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo *serviceModel.UpdateCustomerInfo) (*serviceModel.Customer, error) {
	client, err := s.customerRepository.Update(ctx, uuid, converter.UpdateCustomerInfoToRepo(updateInfo))
	if err != nil {
		return nil, err
	}
	return converter.CustomerFromRepo(client), nil
}
