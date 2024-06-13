package models

import "time"

type User struct {
	ID         int       `json:"id"`
	JwtTokenID string    `json:"jwt_token_id,omitempty"`
	UserName   string    `json:"username,omitempty"`
	Email      string    `json:"email"`
	NickName   string    `json:"nickname,omitempty"`
	Password   string    `json:"password,omitempty"`
	EmailToken string    `json:"email_token,omitempty"`
	TokenHash  string    `json:"token_hash,omitempty"`
	AvatarURL  string    `json:"avatar_url"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
