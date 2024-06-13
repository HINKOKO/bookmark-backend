package main

import (
	"context"
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

func (app *application) verifyTokenFromCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "no token found", http.StatusUnauthorized)
			return
		}

		// token := cookie.Value
		_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// authRequired - middleware that check that authentication is ok - add the userID to context (for availability to next handlers)
func (app *application) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			http.Error(w, "unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Add userID to Context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) verifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims *Claims
		var err error

		// Try to get the token from the cookie
		cookie, err := r.Cookie("refresh_token")
		if err == nil {
			_, claims, err = app.auth.GetTokenFromCookieAndVerify(cookie.Value)
		}

		// If there's no valid cookie token, try the Authorization header
		if err != nil {
			_, claims, err = app.auth.GetTokenFromHeaderAndVerify(w, r)
		}

		// If neither method succeeded, respond with Unauthorized
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := app.DB.GetUserByID(claims.UserID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store claims in the context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
