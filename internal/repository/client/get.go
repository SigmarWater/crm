package client

import (
	"context"

	repoModel "github.com/SigmarWater/crm/internal/repository/model"
)

func (r *repository) Get(
	ctx context.Context,
	uuid string,
) (*repoModel.Client, error) {
	return &repoModel.Client{}, nil
}
