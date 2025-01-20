package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"poultry-management.com/pkg/repo"
)

type LocationHandler struct {
	loc    repo.LocationRepo
	Router *gin.Engine
	config Config
}

func NewLocationHandler(l repo.LocationRepo, r *gin.Engine) *LocationHandler {
	lh := &LocationHandler{
		loc:    l,
		Router: r,
	}
	setupLocationRoutes(lh)
	return lh
}

func setupLocationRoutes(lh *LocationHandler) {

	protected := lh.Router.Group("/tenant-location")
	protected.Use(authMiddleware([]byte(lh.config.JWTSecret)))
	{
		protected.POST("/create", lh.createLocation)
		protected.GET("/get", lh.getLocations)
	}
}

func (h *LocationHandler) getLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Locations fetched successfully"})
}
func (h *LocationHandler) createLocation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Locations created successfully"})
}
