package models

import "time"

type Bookmark struct {
	ID          int       `json:"id"`
	Url         string    `json:"url"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	ProjectID   int       `json:"project_id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
