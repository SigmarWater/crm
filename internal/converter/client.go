package converter

import (
	serviceModel "github.com/SigmarWater/crm/internal/model"
	crmV1 "github.com/SigmarWater/crm/pkg/crm_service/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateClientInfoFromProto(
	req *crmV1.CreateClientRequest,
) *serviceModel.CreateClientInfo {
	return &serviceModel.CreateClientInfo{
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Email: req.GetEmail(),
	}
}

func ClientToProto(client *serviceModel.Client) *crmV1.Client {
	return &crmV1.Client{
		Uuid:      client.UUID,
		Name:      client.Name,
		Phone:     client.Phone,
		Email:     client.Email,
		CreatedAt: timestamppb.New(client.CreatedAt),
		UpdatedAt: timestamppb.New(client.UpdatedAt),
	}
}

func UpdateClientInfoFromProto(
	req *crmV1.UpdateClientRequest,
) *serviceModel.UpdateClientInfo {
	return &serviceModel.UpdateClientInfo{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}
}
