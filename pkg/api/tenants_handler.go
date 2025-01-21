package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"poultry-management.com/pkg/repo"
)

type TenantHandler struct {
	repo *repo.Repository
}

func NewTenantHandler(r *repo.Repository, router *gin.Engine) {
	handler := &TenantHandler{repo: r}
	router.POST("/tenants", handler.CreateTenant)
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	role := c.GetString("role")
	if role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only super_admin can create tenants"})
		return
	}

	var req struct {
		Name       string `json:"name" binding:"required"`
		LicenseKey string `json:"license_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.repo.CreateTenant(c.Request.Context(), req.Name, req.LicenseKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
