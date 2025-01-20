package repo

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"poultry-management.com/internal/auth"
	db "poultry-management.com/internal/db/sqlc"
	"poultry-management.com/pkg/domain"
)

type Repository struct {
	*db.Queries
	db *pgxpool.Pool
}

type UserRepo interface {
	CreateUser(c *gin.Context, req auth.SignupRequest) (int32, error)
	GetUser(c *gin.Context, req auth.LoginRequest) (domain.User, error)
}

func NewRepository(dbConn *pgxpool.Pool) *Repository {
	repo := &Repository{
		Queries: db.New(dbConn),
		db:      dbConn,
	}
	admin, _ := auth.HashPassword("admin1")
	repo.createSuperAdmin(context.Background(), "admin1", admin)
	return repo
}

func (r *Repository) CreateUser(c *gin.Context, req auth.SignupRequest) (int32, error) {

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return 0, fmt.Errorf("hashing password: %w", err)
	}
	var id int32
	if req.SuperAdmin && c.GetString("role") == "super_admin" {
		id, err = r.createSuperAdmin(c, req.Password, hashedPassword)
	} else if c.GetString("role") == "admin" || c.GetString("role") == "master" || c.GetString("role") == "super_admin" {
		if req.TenantID == 0 {
			return 0, fmt.Errorf("please provide a valid tenant id")
		}
		var tid int32 = int32(req.TenantID)
		id, err = r.createOtherUser(c, req.Username, hashedPassword, req.Role, &tid)
	}

	if err != nil {
		return 0, fmt.Errorf("error creating %w", err)
	}

	return id, nil
}

func (r *Repository) GetUser(c *gin.Context, req auth.LoginRequest) (domain.User, error) {
	var user domain.User
	var hashedpwd string
	var err error
	if req.SuperAdmin {
		hashedpwd, user, err = r.getSuperAdminByUsername(c.Request.Context(), req.Username)
	} else {
		hashedpwd, user, err = r.getuserByUsername(c.Request.Context(), req.Username)
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("error creating %w", err)
	}

	if err = auth.CheckPasswordHash(req.Password, hashedpwd); err != nil {
		return domain.User{}, fmt.Errorf(
			"invalid credentials %w", err,
		)
	}

	return user, nil
}

func (r *Repository) getSuperAdminByUsername(ctx context.Context, username string) (string, domain.User, error) {
	user, err := r.Queries.GetSuperAdminByUsername(ctx, username)
	if err != nil {
		return "", domain.User{}, fmt.Errorf("error retrieving super admin: %w", err)
	}
	return user.PasswordHash, domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      "super_admin",
		TenantID:  0,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *Repository) getuserByUsername(ctx context.Context, username string) (string, domain.User, error) {
	user, err := r.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return "", domain.User{}, fmt.Errorf("error retrieving super admin: %w", err)
	}
	return user.PasswordHash, domain.User{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		TenantID:  *user.TenantID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *Repository) createSuperAdmin(ctx context.Context, userName string, hashedPassword string) (int32, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	u, err := r.Queries.CreateSuperAdmin(ctx, db.CreateSuperAdminParams{
		Username:     userName,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return 0, fmt.Errorf("creating super admin: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return u.ID, nil
}

func (r *Repository) createOtherUser(ctx context.Context, userName string, hashedPassword string, role string, tenantID *int32) (int32, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	u, err := r.Queries.CreateUser(ctx, db.CreateUserParams{
		Username:     userName,
		PasswordHash: hashedPassword,
		Role:         role,
		TenantID:     tenantID,
	})
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}
	return u.ID, err
}
