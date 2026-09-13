package client

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	appError "github.com/SigmarWater/crm/internal/errors"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Delete(ctx context.Context, uuid string) error {
	deleteBuilder := sq.Delete("clients").
		Where(sq.Eq{"uuid": uuid}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING uuid")

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		log.Printf("failed to build delete query: %v\n", err)
		return err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var deletedUUID string

	if err := row.Scan(&deletedUUID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appError.ErrClientNotFound
		}
		log.Printf("failed to scan deleted client uuid: %v\n", err)
		return err
	}

	return nil
}
