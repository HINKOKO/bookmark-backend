package dbrepo

import (
	"bookmarks/internal/models"
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/markbates/goth"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) GetProjectsByCategory(category string) ([]*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var projects []*models.Project

	query := `select p.id, p.name, p.category_id, c.category FROM projects p JOIN categories c ON p.category_id = c.id where c.category = $1`
	rows, err := m.DB.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.CategoryID,
			&p.Category,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &p)
	}
	return projects, nil
}

func (m *PostgresDBRepo) GetProjectResources(projectID int) ([]*models.Bookmark, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var bookmarks []*models.Bookmark

	query := `SELECT id, url, title, description, user_id, project_id, created_at,
		updated_at FROM bookmarks WHERE project_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(
			b.ID,
			b.Url,
			b.Title,
			b.Description,
			b.UserID,
			b.ProjectID,
			b.CreatedAt,
			b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, &b)
	}
	return bookmarks, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var u models.User

	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (m *PostgresDBRepo) InsertNewUser(username, email, hashedPassword string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var userID int

	stmt := `INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := m.DB.QueryRowContext(ctx, stmt, username, email, hashedPassword, time.Now(), time.Now()).Scan(&userID)

	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (m *PostgresDBRepo) StoreUserInDB(userID string, user *goth.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	realID, _ := strconv.Atoi(userID)
	fakeMail := "any"
	fakePass := "pass"

	stmt := `INSERT INTO users (id, email, password_hash, username, avatar_url) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.ExecContext(ctx, stmt, realID, fakeMail, fakePass, user.NickName, user.AvatarURL)
	if err != nil {
		log.Println("failed to insert user", err)
		return err
	}
	return nil
}
