package client

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	appError "github.com/SigmarWater/crm/internal/errors"
	repoModel "github.com/SigmarWater/crm/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(
	ctx context.Context,
	uuid string,
) (*repoModel.Client, error) {
	selectBuilder := sq.Select("uuid",
		"name",
		"phone",
		"email",
		"created_at",
		"updated_at").
		PlaceholderFormat(sq.Dollar).
		From("clients").
		Where(sq.Eq{"uuid": uuid})

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		log.Printf("failed select query: %v\n", err)
		return nil, err
	}

	row := r.pool.QueryRow(ctx, sql, args...)

	var client repoModel.Client

	if err := row.Scan(
		&client.UUID,
		&client.Name,
		&client.Phone,
		&client.Email,
		&client.CreatedAt,
		&client.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrClientNotFound
		}

		return nil, err
	}

	return &client, nil
}
