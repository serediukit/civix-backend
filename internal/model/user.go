package model

import "time"

type User struct {
	UserID       uint64    `json:"-"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	PhoneNumber  string    `json:"phone_number"`
	AvatarUrl    string    `json:"avatar_url"`
	RegCityID    string    `json:"-"`
	RegTime      time.Time `json:"reg_time"`
	UpdTime      time.Time `json:"upd_time"`
	DelTime      time.Time `json:"del_time"`
}
