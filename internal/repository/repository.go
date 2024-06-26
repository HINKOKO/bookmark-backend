package repository

import (
	"bookmarks/internal/models"
	"database/sql"
	"time"

	"github.com/markbates/goth"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetProjectsByCategory(category string) ([]*models.Project, error)
	// GetProjectResources(projectID int) ([]*models.Bookmark, error)
	GetResourcesByCategoryAndProject(category, project string) ([]*models.Bookmark, error)
	InsertBookmark(bkm *models.Bookmark) error

	GetUserByEmail(email string) (models.User, error)
	GetUserByID(userID int) (*models.User, error)
	StoreUserInDB(userID string, user *goth.User) error

	// Tokens related functions
	StoreTokenPairs(userID int, accessToken, refreshToken string, expiry time.Time) error
	DeleteTokensPairOnLogOut(userID int) error

	FetchUserFromDB(userID string) (models.User, error)
	// email confirmation && classic authentication function
	// mail confirmation related function
	GetUserByConfirmationToken(token string) (*models.User, error)
	VerifyUser(userID int) error
	CheckEmailConflict(email string) (bool, error)
	InsertNewUser(username, email, password, emailToken, defaultAvatar string) (int, error)

	// Contributors functions
	GetContributors() ([]*models.User, error)

	SaveAvatarURL(userID int, avatarURL string) error
	GetBookmarksByUser(userID int) ([]map[string]interface{}, error)
}
