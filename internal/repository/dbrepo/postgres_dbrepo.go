package dbrepo

import (
	"bookmarks/internal/models"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

/* Bookmarks functions - to retrieve, to modify, to insert */
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

func (m *PostgresDBRepo) GetResourcesByCategoryAndProject(category, project string) ([]*models.Bookmark, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var resources []*models.Bookmark

	query := `SELECT b.id, b.type, b.description, b.url FROM bookmarks b
		JOIN projects p ON b.project_id = p.id
		JOIN categories c ON p.category_id = c.id
		WHERE c.category = $1 AND p.name = $2`

	rows, err := m.DB.QueryContext(ctx, query, category, project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Bookmark
		err := rows.Scan(&r.ID, &r.Type, &r.Description, &r.Url)
		if err != nil {
			return nil, err
		}
		resources = append(resources, &r)
	}
	return resources, nil
}

func (m *PostgresDBRepo) InsertBookmark(bkm *models.Bookmark) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `
	INSERT INTO bookmarks (url, description, user_id, project_id, type)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt, bkm.Url, bkm.Description, bkm.UserID, bkm.ProjectID, bkm.Type)

	return err
}

func (m *PostgresDBRepo) GetContributors() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var conts []*models.User
	query := `SELECT id, username, email, coalesce(nickname, ''), coalesce(avatar_url, '') FROM users`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cont models.User
		err = rows.Scan(
			&cont.ID,
			&cont.UserName,
			&cont.Email,
			&cont.NickName,
			&cont.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		conts = append(conts, &cont)
	}
	return conts, nil
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

func (m *PostgresDBRepo) GetUserByID(userID int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var u models.User
	query := `SELECT id, username, avatar_url FROM users WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.AvatarURL,
	)
	if err != nil {
		return &u, err
	}
	return &u, nil
}

// InsertNewUser - Register a new 'classic' user - combination email + password
func (m *PostgresDBRepo) InsertNewUser(username, email, password, emailToken, defaultAvatar string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var userID int
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (username, email, password_hash, email_token, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = m.DB.QueryRowContext(ctx, stmt, username, email, hashedPassword, emailToken, defaultAvatar, time.Now(), time.Now()).Scan(&userID)

	if err != nil {
		return 0, err
	}
	return userID, nil
}

// StoreUserInDB - stores a new user who log/register with OAUTH (Github provider)
func (m *PostgresDBRepo) StoreUserInDB(userID string, user *goth.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	fakeMail := "any"
	fakePass := "pass"

	stmt := `INSERT INTO users (email, password_hash, username, avatar_url, jwt_token_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.ExecContext(ctx, stmt, fakeMail, fakePass, user.NickName, user.AvatarURL, userID)
	if err != nil {
		log.Println("StoreUserInDB:: failed to insert user", err)
		return err
	}
	return nil
}

func (m *PostgresDBRepo) StoreTokenPairs(userID int, accessToken, refreshToken string, expiry time.Time) error {
	stmt := `INSERT INTO tokens (user_id, access_token, refresh_token, expiry_date)
		VALUES ($1, $2, $3, $4)`
	_, err := m.DB.Exec(stmt, userID, accessToken, refreshToken, expiry)
	return err
}

// GetTokens - used to check if there are tokens
// func (app *application) GetTokens(userID int) (string, string, error) {
//     var accessToken, refreshToken string
//     err := app.DB.QueryRow("SELECT access_token, refresh_token FROM tokens WHERE user_id = $1 AND expiry_date > CURRENT_TIMESTAMP ORDER BY created_at DESC LIMIT 1", userID).Scan(&accessToken, &refreshToken)
//     if err != nil {
//         return "", "", err
//     }
//     return accessToken, refreshToken, nil
// }

// FetchUserFromDB - fetch a user by ID to give information to dashboard protected route
func (m *PostgresDBRepo) FetchUserFromDB(userID string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var u models.User
	log.Println("FetchUserFromDB:: userID just before querying db => ", userID)

	query := `SELECT username, email, avatar_url FROM users WHERE jwt_token_id = $1`

	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&u.UserName,
		&u.Email,
		&u.AvatarURL,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetUserByConfirmationToken - Get a user with by email_token confirmation (when registering for the first time)
func (m *PostgresDBRepo) GetUserByConfirmationToken(token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var user models.User
	query := `SELECT id, username, email, verified FROM users WHERE email_token = $1`

	row := m.DB.QueryRowContext(ctx, query, token)
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Verified,
	)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (m *PostgresDBRepo) VerifyUser(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `UPDATE users SET verified = TRUE, email_token = NULL WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostgresDBRepo) DeleteTokensPairOnLogOut(userID int) error {
	_, err := m.DB.Exec("DELETE FROM tokens WHERE user_id = $1", userID)
	return err
}

// ============== Unused ==========
// func (m *PostgresDBRepo) GetProjectResources(projectID int) ([]*models.Bookmark, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()

// 	var bookmarks []*models.Bookmark

// 	query := `SELECT id, url, title, description, user_id, project_id, created_at,
// 		updated_at FROM bookmarks WHERE project_id = $1`

// 	rows, err := m.DB.QueryContext(ctx, query, projectID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var b models.Bookmark
// 		err := rows.Scan(
// 			b.ID,
// 			b.Url,
// 			b.Title,
// 			b.Description,
// 			b.UserID,
// 			b.ProjectID,
// 			b.CreatedAt,
// 			b.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		bookmarks = append(bookmarks, &b)
// 	}
// 	return bookmarks, nil
// }

// func (db *DB) VerifyUser(userID int) error {
// 	query := `UPDATE users SET verified = TRUE, confirmation_token = NULL WHERE id = $1`
// 	_, err := db.Exec(query, userID)
// 	if err != nil {
// 		return fmt.Errorf("could not verify user: %w", err)
// 	}
// 	return nil
// }
