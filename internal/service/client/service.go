package client

import (
	"github.com/SigmarWater/crm/internal/repository"
	srv "github.com/SigmarWater/crm/internal/service"
)

var _ srv.ClientService = (*service)(nil)

type service struct {
	clientRepository repository.ClientRepository
}

func NewClientService(clientRepository repository.ClientRepository) *service {
	return &service{
		clientRepository: clientRepository,
	}
}
