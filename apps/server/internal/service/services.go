package service

import (
	"github.com/Niiaks/campusCart/internal/lib/job"
	"github.com/Niiaks/campusCart/internal/repository"
)

type Service struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(repo *repository.Repository, job *job.JobService) *Service {
	authService := NewAuthService(repo.User, repo.Session, job)
	return &Service{
		Auth: authService,
		Job:  job,
	}
}
