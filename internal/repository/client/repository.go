package client

import (
	repo "github.com/SigmarWater/crm/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repo.ClientRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
