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
	// Handle CORS preflight requests
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	token, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
	if err != nil || token == "" {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID := claims.Subject

	userInfo, err := app.DB.FetchUserFromDB(userID)
	if err != nil {
		http.Error(w, "Error fetching user info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func (app *application) GetContributors(w http.ResponseWriter, r *http.Request) {

}
