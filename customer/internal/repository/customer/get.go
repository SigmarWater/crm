package customer

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	appError "github.com/SigmarWater/crm/customer/internal/errors"
	repoModel "github.com/SigmarWater/crm/customer/internal/repository/model"
)

func (r *repository) Get(
	ctx context.Context,
	uuid string,
) (*repoModel.Customer, error) {
	selectBuilder := sq.Select("uuid",
		"name",
		"phone",
		"email",
		"created_at",
		"updated_at").
		PlaceholderFormat(sq.Dollar).
		From("customers").
		Where(sq.Eq{"uuid": uuid})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		log.Printf("failed select query: %v\n", err)
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var client repoModel.Customer

	if err := row.Scan(
		&client.UUID,
		&client.Name,
		&client.Phone,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrCustomerNotFound
		}

		return nil, err
	}

	return &client, nil
}
