package models

import "time"

type Project struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CategoryID int       `json:"category_id"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
