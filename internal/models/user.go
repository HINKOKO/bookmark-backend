package models

import "time"

type User struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
