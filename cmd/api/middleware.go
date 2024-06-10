package main

import (
	"net/http"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // URL de votre frontend
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// func (app *application) validateUser(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var req RegisterRequest
// 		err := json.NewDecoder(r.Body).Decode(&req)

// 		if err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}

// 		if req.Username == "" || req.Email == "" || req.Password == "" {
// 			http.Error(w, "one or several field are missing.", http.StatusBadRequest)
// 			return
// 		}
// 		// password complexity - email format here too ?
// 		next.ServeHTTP(w, r)
// 	})
// }

// func (app *application) enableCORS(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

// 		if r.Method == "OPTIONS" {
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
// 			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type,X-CSRF-Token, Authorization")
// 			return
// 		} else {
// 			h.ServeHTTP(w, r)
// 		}
// 	})
// }

// *** === Old Middleware === ***
// mux.Use(cors.Handler(cors.Options{
// 	AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 	AllowCredentials: true,
// 	MaxAge:           300,
// }))
