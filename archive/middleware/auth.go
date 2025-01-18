package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Get the secret key from environment variable

// AuthContextKey is a custom type for context keys to avoid collisions.
type AuthContextKey string

const (
	UserIDContextKey   AuthContextKey = "userID"
	RoleContextKey     AuthContextKey = "role"
	TenantIDContextKey AuthContextKey = "tenantID"
)

func Authorize(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtKey, nil
			})

			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userRole, ok := claims["role"].(string)
				if !ok {
					http.Error(w, "Invalid role in token", http.StatusUnauthorized)
					return
				}

				if userRole != requiredRole && requiredRole != "any" { //Allow any role if requiredRole is "any"
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				userIDFloat, ok := claims["user_id"].(float64)
				if !ok {
					http.Error(w, "Invalid user ID in token", http.StatusBadRequest)
					return
				}
				userID := int(userIDFloat)

				ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
				ctx = context.WithValue(ctx, RoleContextKey, userRole)

				tenantIDFloat, ok := claims["tenant_id"].(float64)
				if ok { //It means it's not a super admin
					tenantID := int(tenantIDFloat)
					ctx = context.WithValue(ctx, TenantIDContextKey, tenantID)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDContextKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func GetRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(RoleContextKey).(string)
	if !ok {
		return "", fmt.Errorf("role not found in context")
	}
	return role, nil
}

func GetTenantIDFromContext(ctx context.Context) (int, error) {
	tenantID, ok := ctx.Value(TenantIDContextKey).(int)
	if !ok {
		return 0, fmt.Errorf("tenant ID not found in context")
	}
	return tenantID, nil
}
