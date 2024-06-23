package dbrepo

import (
	"bookmarks/internal/models"
	"context"
	"database/sql"
	"log"
	"strconv"
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

func (m *PostgresDBRepo) CheckEmailConflict(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var u models.User

	query := `SELECT id, username, email WHERE = $1`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
	)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if u.UserName != "" {
		return true, nil
	}
	return false, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var u models.User

	query := `SELECT id, jwt_token_id, username, email, password_hash, COALESCE(email_token, ''), token_hash, avatar_url, verified, is_admin, created_at, updated_at FROM users WHERE email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&u.ID,
		&u.JwtTokenID,
		&u.UserName,
		&u.Email,
		&u.Password,
		&u.EmailToken,
		&u.TokenHash,
		&u.AvatarURL,
		&u.Verified,
		&u.IsAdmin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return u, err
	}

	return u, nil
}

func (m *PostgresDBRepo) GetUserByID(userID int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var u models.User
	query := `SELECT id, username, email, avatar_url, verified, is_admin FROM users WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
		&u.AvatarURL,
		&u.Verified,
		&u.IsAdmin,
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
// func (m *PostgresDBRepo) StoreUserInDB(userID string, user *goth.User) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()

// 	fakeMail := "any"
// 	fakePass := "pass"

// 	stmt := `INSERT INTO users (email, password_hash, username, avatar_url) VALUES ($1, $2, $3, $4)`
// 	_, err := m.DB.ExecContext(ctx, stmt, fakeMail, fakePass, user.NickName, user.AvatarURL)
// 	if err != nil {
// 		log.Println("StoreUserInDB:: failed to insert user", err)
// 		return err
// 	}
// 	return nil
// }

func (m *PostgresDBRepo) StoreUserInDB(userID string, user *goth.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// fakeMail := "any"
	fakePass := "nope"
	realID, _ := strconv.Atoi(userID)

	stmt := `INSERT INTO users (id, email, password_hash, username, avatar_url) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.DB.ExecContext(ctx, stmt, realID, user.Email, fakePass, user.NickName, user.AvatarURL)
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

// FetchUserFromDB - fetch a user by ID to give information to dashboard protected route
func (m *PostgresDBRepo) FetchUserFromDB(userID string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	var u models.User
	log.Println("FetchUserFromDB:: userID just before querying db => ", userID)

	query := `SELECT username, COALESCE(email, ''), COALESCE(nickname, ''), password_hash, COALESCE(email_token, ''), COALESCE(token_hash, ''),
	avatar_url, verified, is_admin FROM users WHERE id = $1`

	row := m.DB.QueryRowContext(ctx, query, userID)
	err := row.Scan(
		&u.UserName,
		&u.Email,
		&u.NickName,
		&u.Password,
		&u.EmailToken,
		&u.TokenHash,
		&u.AvatarURL,
		&u.Verified,
		&u.IsAdmin,
	)
	if err != nil {
		log.Println(err)
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

func (m *PostgresDBRepo) SaveAvatarURL(userID int, avatarURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	stmt := `UPDATE users SET avatar_url = $1 WHERE id = $2`
	_, err := m.DB.ExecContext(ctx, stmt, avatarURL, userID)

	return err
}

// Get the bookmarks by user id - all the bookmarks a user fetched
func (m *PostgresDBRepo) GetDashboardStats(userID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	query := `SELECT COUNT(b.id) AS bookmark_count
	FROM public.bookmarks b
	JOIN public.projects p ON b.project_id = p.id
	JOIN public.categories c ON p.category_id = c.id
	WHERE b.user_id = $1`

	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *PostgresDBRepo) GetBookmarksByUser(userID int) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var bkmsByUser []map[string]interface{}

	query := `SELECT b.id, b.url, b.type, b.description, p.name AS project_name
	FROM bookmarks b
	JOIN projects p ON b.project_id = p.id
	WHERE user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url, bkType, bkDesc, bkProjectName string

		err := rows.Scan(
			&id,
			&url,
			&bkType,
			&bkDesc,
			&bkProjectName,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		bkm := map[string]interface{}{
			"id":           id,
			"url":          url,
			"type":         bkType,
			"description":  bkDesc,
			"project_name": bkProjectName,
		}
		bkmsByUser = append(bkmsByUser, bkm)

	}
	return bkmsByUser, nil
}
