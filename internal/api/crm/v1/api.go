package v1

import (
	"github.com/SigmarWater/crm/internal/service"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
)

type api struct {
	crmV1.UnimplementedCRMServiceServer
	clientService service.ClientService
}

func NewAPI(clientService service.ClientService) *api {
	return &api{
		clientService: clientService,
	}
}
