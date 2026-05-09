package models

import "gorm.io/gorm"

type Tasks struct {
	gorm.Model
	Title       string `json:"tilte" gorm:"not null"`
	Description string `json:"description"`
	Completed   bool   `json:"completed" gorm:"default:false"`
	UserID      uint   `json:"user_id" gorm:"not null"`
}

type TaskResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UserID      uint   `json:"user_id"`
	CreatedAt   string `json:"created_at"`
}
