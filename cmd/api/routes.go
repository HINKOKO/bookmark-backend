package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routes - declares all the routes and their respectives protection
func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	// Health route for CI/CD
	mux.Get("/", app.checkHealth)

	// Public routes
	mux.Handle("/", app.verifyToken(http.HandlerFunc(app.Home)))
	mux.Get("/bookmarks/{category}", app.GetProjectsByCategory)
	mux.Get("/bookmarks/{category}/{project}", app.GetResourcesForProject)
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
	mux.Get("/dashboard/{userID}", app.GetDashboardStats)

	fileServer := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/uploads/*", http.StripPrefix("/uploads", fileServer))

	// protected route section - now we are not kidding anymore
	mux.Route("/dashboard", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Post("/upload-avatar", app.UploadAvatar)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.adminRequired)
		mux.Get("/dashboard-panel", app.AdminDashboard)
		mux.Get("/list-users", app.ListUsers)
		mux.Get("/list-users/{userID}/bookmarks", app.ListBookmarksByUser)
	})
	return mux
}
