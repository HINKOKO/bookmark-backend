package models

import "time"

type User struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	NickName  string    `json:"nickname"`
	Password  string    `json:"password"`
	TokenHash string    `json:"tokenhash"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
