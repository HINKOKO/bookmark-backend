package main

import (
	"bookmarks/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go movies up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

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

func (app *application) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		http.Error(w, "Access token is required", http.StatusBadRequest)
		return
	}

	// Here you would implement the logic to fetch user information based on the access token
	// For example, querying a database or making an API call to an authentication server

	// For demonstration purposes, let's return a mock user info
	userInfo := map[string]string{"username": "patrick_cohen"}

	// Serialize user info to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
