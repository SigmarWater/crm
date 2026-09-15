package customer

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	repoModel "github.com/SigmarWater/crm/customer/internal/repository/model"
)

func (r *repository) Create(
	ctx context.Context,
	info *repoModel.CustomerInfo,
) (*repoModel.Customer, error) {
	insertBuilder := sq.Insert("customers").
		PlaceholderFormat(sq.Dollar).
		Columns(
			"name",
			"phone",
			"email",
			"created_at",
			"updated_at",
		).
		Values(info.Name, info.Phone, info.Email, time.Now(), time.Now()).
		Suffix("RETURNING uuid, created_at, updated_at")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		log.Printf("failed insert query: %v\n", err)
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var uuid string
	var createdAt time.Time
	var updatedAt time.Time

	if err := row.Scan(&uuid, &createdAt, &updatedAt); err != nil {
		log.Printf("faled scan uuid: %v\n", err)
		return nil, err
	}

	return &repoModel.Customer{
		UUID:      uuid,
		Name:      info.Name,
		Phone:     info.Phone,
		Email:     info.Email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
