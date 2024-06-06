package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var indexTemplate = `
<p><a href="/auth/google">Log in with google</a></p>`

var userTemplate = `
<h2>Welcome</h2>
<p>Name: {{.Name}}</p>`

var store = sessions.NewCookieStore([]byte("verysecret"))

func (app *application) HandleAuth(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	q.Add("provider", chi.URLParam(r, "provider"))
	r.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(w, r)
}

func (app *application) HandleCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	log.Println(provider)
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

	// Store minimal info about this new user with Oauth - for feeding Dashboard && Contributors page
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

//		app.writeJSON(w, http.StatusOK, u)
//	}
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
