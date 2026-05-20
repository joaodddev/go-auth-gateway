package handler

import (
	"encoding/json"
	"net/http"

	"github.com/joaodddev/go-auth-gateway/internal/auth"
	"github.com/joaodddev/go-auth-gateway/pkg/response"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	jwtManager *auth.JWTManager
}

func NewAuthHandler(jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		jwtManager: jwtManager,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate credentials against user service
	// For demo purposes, we'll accept any non-empty credentials
	if req.Email == "" || req.Password == "" {
		response.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.Generate("user-123", req.Email, "user")
	if err != nil {
		response.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response.Success(w, map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"email": req.Email,
			"role":  "user",
		},
	})
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{
		"status": "valid",
	})
}
