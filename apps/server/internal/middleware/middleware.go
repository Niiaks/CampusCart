package middleware

import (
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middlewares struct {
	Global          *GlobalMiddlewares
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
}

func NewMiddlewares(s *server.Server) *Middlewares {

	var nrApp *newrelic.Application

	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}

	return &Middlewares{
		Global:          NewGlobalMiddlewares(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracing(nrApp),
	}
}
