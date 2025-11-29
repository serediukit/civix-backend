package contracts

import "github.com/serediukit/civix-backend/internal/model"

type GetUserRequest struct {
	UserID uint64 `json:"userid"`
	Email  string `json:"email"`
}

type GetUserResponse struct {
	User *model.User `json:"user"`
}

type UpdateUserRequest struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateUserResponse struct {
	User *model.User `json:"user"`
}
