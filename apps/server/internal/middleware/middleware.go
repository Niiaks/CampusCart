package middleware

import (
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middlewares struct {
	Global          *GlobalMiddlewares
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	Auth            *AuthMiddleware
	Authorization   *AuthorizationMiddleware
}

func NewMiddlewares(s *server.Server, sessionRepo repository.SessionRepo) *Middlewares {

	var nrApp *newrelic.Application

	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}

	isProd := s.Config.Primary.Env != "development"

	return &Middlewares{
		Global:          NewGlobalMiddlewares(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracing(nrApp),
		Auth:            NewAuthMiddleware(sessionRepo, isProd),
		Authorization:   NewAuthorizationMiddleware(),
	}
}
