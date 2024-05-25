package main

import (
	"bookmarks/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) GetProjectsByCategory(w http.ResponseWriter, r *http.Request) {
	var projects []*models.Project
	category := chi.URLParam(r, "category")

	projects, err := app.DB.GetProjectsByCategory(category)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, projects)
}
