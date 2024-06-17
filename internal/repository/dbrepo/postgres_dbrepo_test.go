package dbrepo

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectsByCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer db.Close()

	repo := &PostgresDBRepo{DB: db}

	rows := sqlmock.NewRows([]string{"id", "name", "category_id", "category"}).
		AddRow(1, "Project 1", 1, "system-linux").
		AddRow(2, "Project 2", 1, "system-linux")

	mock.ExpectQuery(`select p.id, p.name, p.category_id, c.category FROM projects p JOIN categories c ON p.category_id = c.id where c.category = \$1`).
		WithArgs("system-linux").
		WillReturnRows(rows)

	projects, err := repo.GetProjectsByCategory("system-linux")

	assert.NoError(t, err)
	assert.Len(t, projects, 2)
	assert.Equal(t, "Project 1", projects[0].Name)
	assert.Equal(t, "Project 2", projects[1].Name)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
