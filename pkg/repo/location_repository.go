package repo

import (
	"github.com/gin-gonic/gin"
	"poultry-management.com/internal/auth"
	"poultry-management.com/pkg/domain"
)

type LocationRepo interface {
	CreateLocation(c *gin.Context, req auth.SignupRequest) (int32, error)
	GetLocation(c *gin.Context, req auth.LoginRequest) (domain.User, error)
}

func (r *Repository) CreateLocation(c *gin.Context, req auth.SignupRequest) (int32, error)

func (r *Repository) GetLocation(c *gin.Context, req auth.LoginRequest) (domain.User, error)
