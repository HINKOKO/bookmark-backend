package repository

import (
	"bookmarks/internal/models"
	"database/sql"

	"github.com/markbates/goth"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetProjectsByCategory(category string) ([]*models.Project, error)
	GetProjectResources(projectID int) ([]*models.Bookmark, error)

	GetUserByEmail(email string) (models.User, error)
	InsertNewUser(username, email, password string) (int, error)
	StoreUserInDB(userID string, user *goth.User) error
}
