package customer

import (
	"context"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	appError "github.com/SigmarWater/crm/customer/internal/errors"
	repoModel "github.com/SigmarWater/crm/customer/internal/repository/model"
)

func (r *repository) Update(
	ctx context.Context,
	uuid string,
	updateInfo *repoModel.UpdateCustomerInfo,
) (*repoModel.Customer, error) {
	updateBuilder := sq.Update("customers").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": uuid}).
		Suffix("RETURNING uuid, name, phone, email, created_at, updated_at")

	if updateInfo.Name != nil {
		updateBuilder = updateBuilder.Set("name", *updateInfo.Name)
	}

	if updateInfo.Phone != nil {
		updateBuilder = updateBuilder.Set("phone", *updateInfo.Phone)
	}

	if updateInfo.Email != nil {
		updateBuilder = updateBuilder.Set("email", *updateInfo.Email)
	}

	updateBuilder = updateBuilder.Set("updated_at", time.Now())

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build update query: %v\n", err)
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var customer repoModel.Customer

	if err := row.Scan(
		&customer.UUID,
		&customer.Name,
		&customer.Phone,
		&customer.Email,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrCustomerNotFound
		}

		log.Printf("failed to scan updated customer: %v\n", err)
		return nil, err
	}

	return &customer, nil
}
