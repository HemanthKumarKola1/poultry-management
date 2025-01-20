package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "poultry-management.com/internal/db/sqlc"
)

type LocationRepo interface {
	CreateNewTenantLocation(ctx context.Context, req db.CreateLocationParams, tenantID int32) (db.Location, error)
	GetLocation(ctx context.Context, locID int32, tenantID int32) (db.Location, error)
}

func (r *Repository) CreateNewTenantLocation(ctx context.Context, req db.CreateLocationParams, tenantID int32) (db.Location, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return db.Location{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("SET search_path TO tenant_%dI", tenantID))
	if err != nil {
		return db.Location{}, err
	}

	loc, err := r.CreateLocation(ctx, req)
	if err != nil {
		return db.Location{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return db.Location{}, err
	}

	return loc, nil
}

func (r *Repository) GetLocation(ctx context.Context, locID int32, tenantID int32) (db.Location, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return db.Location{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("SET search_path TO tenant_%dI", tenantID))
	if err != nil {
		return db.Location{}, err
	}

	loc, err := r.GetLocationByID(ctx, locID)
	if err != nil {
		return db.Location{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return db.Location{}, err
	}

	return loc, nil
}
