package client

import (
	"context"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	appError "github.com/SigmarWater/crm/internal/errors"
	repoModel "github.com/SigmarWater/crm/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Update(
	ctx context.Context,
	uuid string,
	updateInfo *repoModel.UpdateClientInfo,
) (*repoModel.Client, error) {
	updateBuilder := sq.Update("clients").
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

		log.Printf("failed to scan updated client: %v\n", err)
		return nil, err
	}

	return &client, nil
}
