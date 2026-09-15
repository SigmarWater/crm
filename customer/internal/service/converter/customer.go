package converter

import (
	serviceModel "github.com/SigmarWater/crm/customer/internal/model"
	repoModel "github.com/SigmarWater/crm/customer/internal/repository/model"
)

func CreateCustomerInfoToRepo(info *serviceModel.CreateCustomerInfo) *repoModel.CustomerInfo {
	return &repoModel.CustomerInfo{
		Name:  info.Name,
		Phone: info.Phone,
		Email: info.Email,
	}
}

func UpdateCustomerInfoToRepo(info *serviceModel.UpdateCustomerInfo) *repoModel.UpdateCustomerInfo {
	return &repoModel.UpdateCustomerInfo{
		Name:  info.Name,
		Phone: info.Phone,
		Email: info.Email,
	}
}

func CustomerFromRepo(client *repoModel.Customer) *serviceModel.Customer {
	return &serviceModel.Customer{
		UUID:      client.UUID,
		Name:      client.Name,
		Phone:     client.Phone,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}
