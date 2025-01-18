package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"poultry-management.com/internal/auth"
	db "poultry-management.com/internal/db/sqlc"
)

type Repository struct {
	*db.Queries
	db *pgxpool.Pool
}

func NewRepository(dbConn *pgxpool.Pool) *Repository {
	return &Repository{
		Queries: db.New(dbConn),
		db:      dbConn,
	}
}

func (r *Repository) CreateSuperAdmin(ctx context.Context, userName string, password string) (db.CreateSuperAdminRow, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return db.CreateSuperAdminRow{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return db.CreateSuperAdminRow{}, fmt.Errorf("hashing password: %w", err)
	}

	user, err := r.Queries.CreateSuperAdmin(ctx, db.CreateSuperAdminParams{
		Username:     userName,
		PasswordHash: hashedPassword,
	})

	if err != nil {
		return db.CreateSuperAdminRow{}, fmt.Errorf("creating super admin: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return db.CreateSuperAdminRow{}, fmt.Errorf("commit tx: %w", err)
	}

	return user, nil
}
