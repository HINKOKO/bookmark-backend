package dbrepo

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestGetProjectsByCategory - testing that getting project by category behaves correctly
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

func TestGetResourcesByCategoryAndProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer db.Close()

	repo := &PostgresDBRepo{DB: db}
	// test data
	category := "system-linux"
	project := "libasm"

	// Mock the expected results
	rows := sqlmock.NewRows([]string{"id", "type", "description", "url"}).AddRow(1, "tutorial", "Assembly little project", "https://assemblyDesmystified.com")

	// expected query
	mock.ExpectQuery(`SELECT b.id, b.type, b.description, b.url FROM bookmarks b
	JOIN projects p ON b.project_id = p.id
	JOIN categories c ON p.category_id = c.id
	WHERE c.category = \$1 AND p.name = \$2`).WithArgs(category, project).WillReturnRows(rows)

	resources, err := repo.GetResourcesByCategoryAndProject(category, project)
	if err != nil {
		t.Fatalf("error calling GetResourcesByCategoryAndProject: %v", err)
	}

	// Check results
	assert.Equal(t, 1, len(resources), "expected one resources")
	assert.Equal(t, "tutorial", resources[0].Type, "expected resource type to match")
	assert.Equal(t, "Assembly little project", resources[0].Description, "expected description to match")
	assert.Equal(t, "https://assemblyDesmystified.com", resources[0].Url, "expected url links to match")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there was unfulfilled expectations: %s", err)
	}
}
