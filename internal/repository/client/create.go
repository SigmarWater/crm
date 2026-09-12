package client

import (
	"context"

	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

func (r *repository) Create(
	ctx context.Context,
	info *repoModel.ClientInfo,
) (*repoModel.Client, error) {
	return &repoModel.Client{}, nil
}
