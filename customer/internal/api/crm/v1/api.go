package v1

import (
	"github.com/SigmarWater/crm/customer/internal/service"
	customerV1 "github.com/SigmarWater/crm/shared/pkg/customer_service/v1"
)

type api struct {
	customerV1.UnimplementedCustomerServiceServer
	customerService service.CustomerService
}

func NewAPI(customerService service.CustomerService) *api {
	return &api{
		customerService: customerService,
	}
}
