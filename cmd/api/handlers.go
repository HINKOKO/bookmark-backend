package main

import (
	"bookmarks/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/microcosm-cc/bluemonday"
)

// Special handler for endpoint used for sanity check about CI/CD
func (app *application) checkHealth(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, []byte("healthy app"))
}

// Home - Handler for Homepage - rather used for backlog information
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

// GetProjectsByCategory - Handler to retrieve & serve the projects according to category
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

// GetResourcesForProject -  Handler to retrieve & serve the resources for a given project
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

// InsertNewBookmark - Handler to insert a new bookmark in the DB
func (app *application) InsertNewBookmark(w http.ResponseWriter, r *http.Request) {
	var bookmark models.Bookmark

	err := json.NewDecoder(r.Body).Decode(&bookmark)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(bookmark.Url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		log.Println(err)
		http.Error(w, "Invalid URL provided", http.StatusBadRequest)
		return
	}

	// Sanitize the text field 'description' from Bookmark model
	policy := bluemonday.UGCPolicy()
	bookmark.Description = policy.Sanitize(bookmark.Description)

	// Insert Sanitized bookmark into database
	err = app.DB.InsertBookmark(&bookmark)
	if err != nil {
		http.Error(w, "Failed to insert bookmark", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Bookmark added successfully"})
}

// GetUserInfo - Handler to retrieve user info (used accross the screens in FrontEnd - via useAuth context)
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
	log.Println("userid is\t", claims.UserID)

	userID := strconv.Itoa(claims.UserID)

	userInfo, err := app.DB.FetchUserFromDB(userID)
	if err != nil {
		http.Error(w, "Error fetching user info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// GetContributors - Handler to retrieve all contributors
func (app *application) GetContributors(w http.ResponseWriter, r *http.Request) {
	var contributors []*models.User

	contributors, err := app.DB.GetContributors()
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, contributors)
}

/*
========== WARNING ========
DOUBLE FUNCTION WHICH ARE DOING THE SAME - REFACTOR !!!
=============
*/
// ListUsers - Handler to list all users
func (app *application) ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []*models.User

	users, _ = app.DB.GetContributors()

	app.writeJSON(w, http.StatusAccepted, users)
}

// ListBookmarksByUser - Handler to fetch the bookmarks according to a selected User
func (app *application) ListBookmarksByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr) // must match the placeholder in route definition
	if err != nil {
		http.Error(w, "No such user - something really bad here", http.StatusBadRequest)
		return
	}

	bookmarks, err := app.DB.GetBookmarksByUser(userID)
	if err != nil {
		http.Error(w, "error fetching bookmarks", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, bookmarks)
}
