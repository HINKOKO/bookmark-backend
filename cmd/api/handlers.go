package main

import (
	"bookmarks/internal/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*models.User)
	authenticated := ok && user != nil

	var payload = struct {
		Status        string       `json:"status"`
		Message       string       `json:"message"`
		Version       string       `json:"version"`
		Authenticated bool         `json:"authenticated"`
		User          *models.User `json:"user,omitempty"`
	}{
		Status:        "active",
		Message:       "Go movies up and running",
		Version:       "1.0.0",
		Authenticated: authenticated,
		User:          user,
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

func (app *application) GetResourcesForProject(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	project := chi.URLParam(r, "project")

	resources, err := app.DB.GetResourcesByCategoryAndProject(category, project)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resources)

}

func (app *application) InsertNewBookmark(w http.ResponseWriter, r *http.Request) {
	var bkm models.Bookmark

	err := json.NewDecoder(r.Body).Decode(&bkm)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println(bkm)

	err = app.DB.InsertBookmark(&bkm)
	if err != nil {
		http.Error(w, "Failed to insert bookmark", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Bookmark added successfully"})
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
	var contributors []*models.User

	contributors, err := app.DB.GetContributors()
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, contributors)
}
