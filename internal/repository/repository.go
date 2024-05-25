package repository

import (
	"bookmarks/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetProjectsByCategory(category string) ([]*models.Project, error)
	GetProjectResources(projectID int) ([]*models.Bookmark, error)
}
