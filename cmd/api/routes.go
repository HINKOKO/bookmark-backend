package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	// Public routes
	mux.Handle("/", app.verifyToken(http.HandlerFunc(app.Home)))
	mux.Get("/bookmarks/{category}", app.GetProjectsByCategory)
	mux.Get("/bookmarks/{category}/{project}", app.GetResourcesForProject)
	mux.Get("/dashboard", app.Dashboard)
	mux.Get("/auth/{provider}", app.HandleAuth)
	mux.Get("/auth/{provider}/callback", app.HandleCallback)
	mux.Post("/register", app.RegisterNewUser)
	mux.Post("/login", app.ClassicLogin)
	mux.Get("/confirm-email", app.ConfirmEmail)
	mux.Get("/contributors", app.GetContributors)

	// USer information - Feed Dashboard && related screen with user data - Hybrid by now
	mux.Get("/user-info", app.GetUserInfo)
	mux.Handle("/logout", app.verifyToken(http.HandlerFunc(app.Logout)))

	mux.Post("/contributors/insert-bookmark", app.InsertNewBookmark)

	// Posting new resources
	// mux.Get("/contributors/categories", app.GetCategories)
	mux.Get("/contributors/{category}", app.GetProjectsByCategory)
	// mux.Post("/contributors/bookmarks", app.PostNewBookmarkByCategory)

	// protected route section - now we are not kidding anymore
	mux.Route("/contributor", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Get("/dashboard", app.Dashboard)
		// mux.Post("/{category}/{project}/new-resource", app.InsertNewBookmark)
	})

	return mux
}
