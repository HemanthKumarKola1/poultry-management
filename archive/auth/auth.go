package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = "secret-key"

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) AuthenticateSuperAdmin(username, password string) (string, error) {
	var storedPasswordHash string
	var id int
	err := s.db.QueryRow("SELECT id, password_hash FROM public.super_admins WHERE username = $1", username).Scan(&id, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", fmt.Errorf("authenticating super admin: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"role":    "super_admin",
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	return tokenString, nil
}
