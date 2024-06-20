package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// define a max constant - 2MB - to limit avatar image size
const maxUploadSize = 2 * 1024 * 1024
const uploadPath = "./uploads"

func (app *application) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// ensure the directory to store avatarURL about users exists
	// this function will be called by the first who change his default avatar basically
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		http.Error(w, "unable to create directory to save avatarsURL", http.StatusBadRequest)
		return
	}
	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		app.errorJSON(w, err, http.StatusSeeOther)
		return
	}
	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create temp file in the within avatar upload directory
	tmpFile, err := os.CreateTemp(uploadPath, "upload-*.png")
	if err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	// copy the actual file to temporary file
	fileSize, err := io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}

	// Check size limitation
	if fileSize > maxUploadSize {
		http.Error(w, "avatar file too large", http.StatusInternalServerError)
		return
	}

	// Construct the URL path
	avatarURL := "/uploads/" + filepath.Base(tmpFile.Name())

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "unable to extract cookie", http.StatusInternalServerError)
		return
	}

	tokenStr := cookie.Value

	// here we save avatarURL to the user's profile in the database
	_, claims, _ := app.auth.GetTokenFromCookieAndVerify(tokenStr)
	userID := claims.UserID

	err = app.DB.SaveAvatarURL(userID, avatarURL)
	if err != nil {
		http.Error(w, "unable to save avatar in database", http.StatusInternalServerError)
		return
	}

	fullAvatarURL := "http://localhost:8080" + avatarURL
	// return avatar URL to frontend
	response := map[string]string{"avatar_url": fullAvatarURL}
	app.writeJSON(w, http.StatusOK, response)
}

func (app *application) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("user_id")
	log.Println(userIDstr)

}
