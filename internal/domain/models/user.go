package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"password" gorm:"not null"`
	Name     string `json:"name" gorm:"not null"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	EMAIL     string `json:"email"`
	NAME      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
