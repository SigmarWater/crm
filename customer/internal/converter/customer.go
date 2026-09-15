package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	serviceModel "github.com/SigmarWater/crm/customer/internal/model"
	customerV1 "github.com/SigmarWater/crm/shared/pkg/customer_service/v1"
)

func CreateCustomerInfoFromProto(
	req *customerV1.CreateCustomerRequest,
) *serviceModel.CreateCustomerInfo {
	return &serviceModel.CreateCustomerInfo{
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Email: req.GetEmail(),
	}
}

func CustomerToProto(customer *serviceModel.Customer) *customerV1.Customer {
	return &customerV1.Customer{
		Uuid:      customer.UUID,
		Name:      customer.Name,
		Phone:     customer.Phone,
		Email:     customer.Email,
		CreatedAt: timestamppb.New(customer.CreatedAt),
		UpdatedAt: timestamppb.New(customer.UpdatedAt),
	}
}

func UpdateCustomerInfoFromProto(
	req *customerV1.UpdateCustomerRequest,
) *serviceModel.UpdateCustomerInfo {
	return &serviceModel.UpdateCustomerInfo{
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	}
}
