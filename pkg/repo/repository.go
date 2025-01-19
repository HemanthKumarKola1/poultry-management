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

func (r *Repository) CreateUser(ctx context.Context, req auth.SignupRequest) (int32, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("hashing password: %w", err)
	}
	var id int32
	if req.SuperAdmin {
		id, err = r.createSuperAdmin(ctx, req.Password, hashedPassword)
	} else {
		if req.TenantID == 0 {
			return 0, fmt.Errorf("please provide valid tenant id")
		}
		var tid int32 = int32(req.TenantID)
		id, err = r.createOtherUser(ctx, req.Username, hashedPassword, req.Role, &tid)
	}

	if err != nil {
		return 0, fmt.Errorf("error creating %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return id, nil
}

func (r *Repository) createSuperAdmin(ctx context.Context, userName string, hashedPassword string) (int32, error) {
	u, err := r.Queries.CreateSuperAdmin(ctx, db.CreateSuperAdminParams{
		Username:     userName,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return 0, fmt.Errorf("creating super admin: %w", err)
	}
	return u.ID, nil
}

func (r *Repository) createOtherUser(ctx context.Context, userName string, hashedPassword string, role string, tenantID *int32) (int32, error) {
	u, err := r.Queries.CreateUser(ctx, db.CreateUserParams{
		Username:     userName,
		PasswordHash: hashedPassword,
		Role:         role,
		TenantID:     tenantID,
	})
	return u.ID, err
}
