package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

// 	app.writeJSON(w, http.StatusOK, u)
// }

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%+v\n", req)

	id, err := app.DB.InsertNewUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Println("failed to record new user")
		return
	}
	app.writeJSON(w, http.StatusAccepted, id)
}
