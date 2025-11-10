package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID       uint64         `json:"user_id"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"-"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}
