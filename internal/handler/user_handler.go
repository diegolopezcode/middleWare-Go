package handler

import (
	"task-api/internal/config"
	"task-api/internal/domain/models"

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
