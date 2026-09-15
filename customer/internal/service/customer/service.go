package customer

import (
	"github.com/SigmarWater/crm/customer/internal/repository"
	srv "github.com/SigmarWater/crm/customer/internal/service"
)

var _ srv.CustomerService = (*service)(nil)

type service struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) *service {
	return &service{
		customerRepository: customerRepository,
	}
}
