package model

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}
