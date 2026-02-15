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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Global.Recover())
	r.Use(middleware.Logger)

	//rate limit
	r.Use(httprate.Limit(
		10,
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

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(mw.Auth.Authenticate)
			r.Post("/auth/logout", h.Auth.Logout())
			r.Get("/auth/me", h.Auth.GetCurrentUser())
		})
	})
	return r
}
