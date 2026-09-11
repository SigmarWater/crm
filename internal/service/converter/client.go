package converter

import (
	serviceModel "github.com/SigmarWater/crm/internal/model"
	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

func CreateClientInfoToRepo(info *serviceModel.CreateClientInfo) *repoModel.ClientInfo {
	return &repoModel.ClientInfo{
		Name:  info.Name,
		Phone: info.Phone,
		Email: info.Email,
	}
}

func UpdateClientInfoToRepo(info *serviceModel.UpdateClientInfo) *repoModel.UpdateClientInfo {
	return &repoModel.UpdateClientInfo{
		Name:  info.Name,
		Phone: info.Phone,
		Email: info.Email,
	}
}

func ClientFromRepo(client *repoModel.Client) *serviceModel.Client {
	return &serviceModel.Client{
		UUID:      client.UUID,
		Name:      client.Name,
		Phone:     client.Phone,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}
