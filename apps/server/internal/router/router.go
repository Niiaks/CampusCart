package router

import (
	"net/http"
	"time"

	"github.com/Niiaks/campusCart/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func NewRouter(h *handler.Handlers) chi.Router {

	r := chi.NewRouter()

	//Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//rate limit
	r.Use(httprate.Limit(
		10,
		time.Minute,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error": "Rate-limited. Please, slow down."}`, http.StatusTooManyRequests)
		}),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.Health.CheckHealth)
	})
	return r
}
