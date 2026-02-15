package handler

import (
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Auth    *AuthHandler
}

func NewHandlers(s *server.Server, authService *service.AuthService) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Auth:    NewAuthHandler(s, authService),
	}
}
