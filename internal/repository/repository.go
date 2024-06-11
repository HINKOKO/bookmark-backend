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
	StoreUserInDB(userID string, user *goth.User) error

	FetchUserFromDB(userID string) (models.User, error)
	// email confirmation && classic authentication function
	// mail confirmation related function
	GetUserByConfirmationToken(token string) (*models.User, error)
	VerifyUser(userID int) error
	InsertNewUser(username, email, password, emailToken, defaultAvatar string) (int, error)

	// Contributors functions
	GetContributors() ([]*models.User, error)
}
