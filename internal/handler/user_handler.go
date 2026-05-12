package handler

import (
	"encoding/json"
	"net/http"
	"task-api/internal/config"
	"task-api/internal/domain/models"
	"task-api/pkg/utils"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type AuthHandler struct {
	db *config.DB
}

func NewAuthHanlder(db *config.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required, min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.db.Create(&user).Error; err != nil {
		http.Error(w, "Email already exist", http.StatusConflict)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	})

	tokenString, _ := token.SignedString([]byte(config.Load().JWTSecret))

	response := AuthResponse{
		Token: tokenString,
		User: models.UserResponse{
			ID:        user.ID,
			EMAIL:     user.Email,
			NAME:      user.Name,
			CreatedAt: user.CreatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "applicaation/json")
	json.NewEncoder(w).Encode(response)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
	})
	tokenString, _ := token.SignedString([]byte(config.Load().JWTSecret))

	response := AuthResponse{
		Token: tokenString,
		User: models.UserResponse{
			ID:        user.ID,
			NAME:      user.Name,
			EMAIL:     user.Email,
			CreatedAt: user.CreatedAt.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
