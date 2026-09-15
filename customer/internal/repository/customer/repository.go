package customer

import (
	"github.com/jackc/pgx/v5/pgxpool"

	repo "github.com/SigmarWater/crm/customer/internal/repository"
)

var _ repo.CustomerRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
