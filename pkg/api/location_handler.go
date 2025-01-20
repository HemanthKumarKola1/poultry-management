package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "poultry-management.com/internal/db/sqlc"
	"poultry-management.com/pkg/repo"
)

type LocationHandler struct {
	locrepo repo.LocationRepo
	Router  *gin.Engine
	config  Config
}

func NewLocationHandler(l repo.LocationRepo, r *gin.Engine) *LocationHandler {
	lh := &LocationHandler{
		locrepo: l,
		Router:  r,
	}
	setupLocationRoutes(lh)
	return lh
}

func setupLocationRoutes(lh *LocationHandler) {

	protected := lh.Router.Group("/tenant-location")
	protected.Use(authLocationMiddelware([]byte(lh.config.JWTSecret)))
	{
		protected.POST("/create", lh.createLocation)
		protected.GET("/get", lh.getLocations)
	}
}

func (h *LocationHandler) getLocations(c *gin.Context) {
	locIDStr := c.GetHeader("location")
	intVal, err := strconv.Atoi(locIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	location, err := h.locrepo.GetLocation(c.Request.Context(), int32(intVal), int32(c.GetInt("tenant_id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"locations": location})
}

func (h *LocationHandler) createLocation(c *gin.Context) {
	var req db.CreateLocationParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loc, err := h.locrepo.CreateNewTenantLocation(c.Request.Context(), req, int32(c.GetInt("tenant_id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Locations created successfully", "location": loc})
}
