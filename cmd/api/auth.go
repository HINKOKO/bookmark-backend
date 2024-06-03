package main

import (
	"bookmarks/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthenticating struct {
	Username    string
	Email       string
	Password    string
	Verified    bool
	VerifyToken string
}

func (app *application) CheckPassword(u models.User, plainPass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPass))
}
