package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"poultry-management.com/internal/auth"
	"poultry-management.com/pkg/repo"
)

type Auth struct {
	Config     Config
	Repository *repo.Repository
	Router     *gin.Engine
}

type Config struct {
	JWTSecret string
}

func NewAuth(config Config, repo *repo.Repository, e *gin.Engine) {
	a := &Auth{
		Config:     config,
		Repository: repo,
		Router:     e,
	}

	setupRoutes(a)
}

func setupRoutes(Auth *Auth) {
	Auth.Router.POST("/signup", Auth.signupHandler)
	Auth.Router.POST("/login", Auth.loginHandler)

	protected := Auth.Router.Group("/api")
	protected.Use(authMiddleware([]byte(Auth.Config.JWTSecret)))
	{
		protected.GET("/ping", Auth.pingHandler)
	}
}

func (Auth *Auth) signupHandler(c *gin.Context) {
	var req auth.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Auth.Repository.CreateSuperAdmin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateJWT(user.ID, 0, user.Username, []byte(Auth.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (Auth *Auth) loginHandler(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := Auth.Repository.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, *user.TenantID, user.Username, []byte(Auth.Config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (Auth *Auth) pingHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	tenantID := c.GetInt("tenant_id")
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("pong user id: %d, tenant id: %d, username: %s", userID, tenantID, username),
	})
}

func authMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		claims, err := auth.ValidateJWT(tokenString, jwtKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
