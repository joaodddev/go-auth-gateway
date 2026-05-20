package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/joaodddev/go-auth-gateway/pkg/response"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type AuthMiddleware struct {
	jwtManager  *JWTManager
	publicPaths map[string]bool
}

func NewAuthMiddleware(jwtManager *JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		publicPaths: map[string]bool{
			"/health":   true,
			"/login":    true,
			"/register": true,
			"/public":   true,
		},
	}
}

func (am *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if path is public
		if am.publicPaths[r.URL.Path] {
			next(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := am.jwtManager.Validate(token)
		if err != nil {
			response.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func GetUserFromContext(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(UserContextKey).(*Claims); ok {
		return claims
	}
	return nil
}
