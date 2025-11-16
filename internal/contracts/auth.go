package contracts

import (
	"time"

	"github.com/serediukit/civix-backend/internal/model"
)

type RegisterRequest struct {
	Email       string         `json:"email" binding:"required,email"`
	Password    string         `json:"password" binding:"required,min=6"`
	Name        string         `json:"name" binding:"required"`
	Surname     string         `json:"surname"`
	PhoneNumber string         `json:"phone_number" binding:"min=10,max=13"`
	Location    model.Location `json:"location"`
}

type RegisterResponse struct {
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	RegTime time.Time `json:"reg_time"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  model.Token `json:"access_token"`
	RefreshToken model.Token `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  model.Token `json:"access_token"`
	RefreshToken model.Token `json:"refresh_token"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
