package models

import "time"

type Rating struct {
	ID         int       `json:"id"`
	UserID     int       `json:"id"`
	BookmarkID int       `json:"bookmark_id"`
	Rating     int       `json:"rating"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
