package service

import "github.com/Niiaks/campusCart/internal/repository"

type Service struct {
	Auth *AuthService
}

func NewServices(repo *repository.Repository) *Service {
	authService := NewAuthService(repo.User, repo.Session)
	return &Service{
		Auth: authService,
	}
}
