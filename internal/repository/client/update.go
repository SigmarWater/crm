package client

import (
	"context"

	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

func (r *repository) Update(
	ctx context.Context,
	uuid string,
	updateInfo *repoModel.UpdateClientInfo,
) (*repoModel.Client, error) {
	return &repoModel.Client{}, nil
}
