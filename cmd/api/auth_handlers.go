package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// var store = sessions.NewCookieStore([]byte("verysecret"))

func (app *application) ClassicLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUserByEmail(loginReq.Email)
	if err != nil {
		http.Error(w, "no such user in our dataabse", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(loginReq.Password), []byte(user.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// generate a token for new user
	classicUser := jwtUser{
		ID:       fmt.Sprintf(string(user.ID)),
		Username: user.UserName,
	}
	tokenString, _ := app.auth.GenerateTokenPair(&classicUser)

	// Optionally, store the refresh token in the database
	// err = app.DB.StoreRefreshToken(user.ID, tokenString.RefreshToken)
	// if err != nil {
	// 	http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
	// 	return
	// }

	http.SetCookie(w, app.auth.GetRefreshCookie(tokenString.RefreshToken))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"access_token": tokenString.Token})

}

func (app *application) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "token must be expired - please start over to register", http.StatusBadRequest)
		return
	}

	// Using confirmation token - we retrieve corresponding user (pre-registered)
	user, err := app.DB.GetUserByConfirmationToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid or expired token", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Evrything valid, we UPDATE the user as verified - Register is complete !
	err = app.DB.VerifyUser(user.ID)
	if err != nil {
		// If an error occurred while updating the user's verification status, return a server error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "http://localhost:5173/email-confirmed", http.StatusAccepted)
}

func (app *application) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check in log the request
	// log.Printf("%+v\n", req)

	// Request is properly formatted - pretending new user deserves an email confirmation
	// generate a random token
	randomString := generateRandomString(32)
	log.Println(randomString)

	defaultAvatar := fmt.Sprintf("https://api.dicebear.com/8.x/pixel-art/svg?seed=%s", req.Username)

	err = app.sendConfirmationEmail(req.Email, randomString)
	if err != nil {
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		return
	}

	id, err := app.DB.InsertNewUser(req.Username, req.Email, req.Password, randomString, defaultAvatar)
	if err != nil {
		log.Println("Failed to register that new user")
		return
	}
	// Optionally, you can redirect the user to a success page
	http.Redirect(w, r, "http://localhost:5173/email-confirmation?redirect=login", http.StatusAccepted)
	app.writeJSON(w, http.StatusAccepted, id)
}

func (app *application) HandleAuth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", chi.URLParam(r, "provider"))
	r.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(w, r)
}

func (app *application) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("Error completing user auth: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := uuid.New().String()

	// generate a token for new user
	u := jwtUser{
		ID:       userID,
		Username: user.NickName,
	}
	tokenString, _ := app.auth.GenerateTokenPair(&u)
	refreshCookie := app.auth.GetRefreshCookie(tokenString.RefreshToken)
	http.SetCookie(w, refreshCookie)

	// store that new user in DB
	err = app.DB.StoreUserInDB(userID, &user)
	if err != nil {
		log.Printf("Error storing user in database: %v", err)
		http.Error(w, "Error storing user in database", http.StatusInternalServerError)
		return
	}

	// JSONify the user data fetched from oauth provider
	userData, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		log.Println("error marshaling user data", err)
		http.Error(w, "error encoding user data", http.StatusInternalServerError)
		return
	}
	log.Printf("user data from github %s", string(userData))

	redirectURL := fmt.Sprintf("http://localhost:5173/dashboard?accessToken=%s&user=%s", tokenString.Token, url.QueryEscape(string(userData)))
	http.Redirect(w, r, redirectURL, http.StatusFound)
	app.writeJSON(w, http.StatusOK, user)
}

// Dashboard - handler
func (app *application) Dashboard(w http.ResponseWriter, r *http.Request) {

}

// func (app *application) Login(w http.ResponseWriter, r *http.Request) {

// 	u, err := app.DB.GetUserByEmail(email)
// 	if err != nil {
// 		log.Println("no such user apprently")

// 	}

// 	err = app.CheckPassword(u, password)
// 	if err != nil {
// 		app.writeJSON(w, http.StatusBadRequest, nil)
// 	}

// 		app.writeJSON(w, http.StatusOK, u)
// 	}
// func (app *application) InitOauth(w http.ResponseWriter, r *http.Request) {
// 	provider := chi.URLParam(r, "provider")
// 	log.Println("provider from URLParam is => ", provider)

// 	q := r.URL.Query()
// 	q.Add("provider", provider)
// 	r.URL.RawQuery = q.Encode()
// 	gothic.BeginAuthHandler(w, r)
// }

// func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
// 	provider := chi.URLParam(r, "provider")
// 	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

// 	user, err := gothic.CompleteUserAuth(w, r)
// 	if err != nil {
// 		log.Printf("Error completing user auth: %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)

// 	// generate a token for new user
// 	u := jwtUser{
// 		// ID:       user.UserID,
// 		// Username: user.Name,
// 		ID:       "1",
// 		Username: "patrick_cohen",
// 	}
// 	tokenString, _ := app.auth.GenerateTokenPair(&u)
// 	refreshCookie := app.auth.GetRefreshCookie(tokenString.RefreshToken)

// 	// http.SetCookie(w, refreshCookie)
// 	// session, _ := store.Get(r, "session-name")
// 	// session.Values["user"] = user
// 	// session.Save(r, w)

// 	// w.Header().Set("Content-Type", "application/json")
// 	// json.NewEncoder(w).Encode(user)
// 	// app.writeJSON(w, http.StatusAccepted, tokenString)
// 	http.SetCookie(w, refreshCookie)

// 	redirectURL := fmt.Sprintf("http://localhost:5173/?token=%s", tokenString.Token)
// 	http.Redirect(w, r, redirectURL, http.StatusFound)
// }
