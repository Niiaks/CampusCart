package router

import (
	"net/http"
	"time"

	"github.com/Niiaks/campusCart/internal/handler"
	customMiddleware "github.com/Niiaks/campusCart/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func NewRouter(h *handler.Handlers, mw *customMiddleware.Middlewares) chi.Router {

	r := chi.NewRouter()

	//Global middleware
	r.Use(customMiddleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Global.CORS())
	r.Use(mw.Tracing.NewRelicMiddleware())
	r.Use(mw.ContextEnhancer.EnhanceContext)
	r.Use(mw.Global.RequestLogger())
	r.Use(mw.Global.Recover())
	r.Use(mw.Global.Secure())

	// Custom error handlers for consistent JSON responses
	r.NotFound(mw.Global.NotFoundHandler())
	r.MethodNotAllowed(mw.Global.MethodNotAllowedHandler())

	//rate limit
	r.Use(httprate.Limit(
		100,
		time.Minute,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error": "Rate-limited. Please, slow down."}`, http.StatusTooManyRequests)
		}),
	))

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		if err := h.OpenAPI.ServeOpenApiUI(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Serve static files (openapi.json, etc.)
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.Health.CheckHealth)

		// Public auth routes
		r.Post("/auth/login", h.Auth.Login())
		r.Post("/auth/register", h.Auth.Register())
		r.Post("/auth/verify-email", h.Auth.VerifyEmail())

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.Authenticate)
			r.Post("/auth/logout", h.Auth.Logout())
			r.Get("/auth/me", h.Auth.GetCurrentUser())
		})

		// Brand profile
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.Authenticate)
			r.Get("/brands/me", h.Brand.GetCurrent())
			r.Patch("/brands/me", h.Brand.Update())
		})

		// Public category routes
		r.Get("/categories", h.Category.GetAll())
		r.Get("/categories/{id}", h.Category.GetByID())
		r.Get("/categories/{id}/attributes", h.Category.GetAttributes())

		// Listings
		r.Get("/listings", h.Listing.List())
		r.Get("/listings/{id}", h.Listing.Get())
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.Authenticate)
			r.Post("/listings", h.Listing.Create())
			r.Patch("/listings/{id}", h.Listing.Update())
			r.Delete("/listings/{id}", h.Listing.Delete())
			r.Post("/listings/upload-signature", h.Listing.UploadSignature())
		})

		// Admin category routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.Authenticate)
			r.Use(mw.Authorization.Authorize)
			r.Post("/categories", h.Category.Create())
			r.Patch("/categories/{id}", h.Category.Update())
			r.Delete("/categories/{id}", h.Category.Delete())
		})
	})
	return r
}
