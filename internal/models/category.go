package models

import "time"

type Category struct {
	ID        int       `json:"id"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
