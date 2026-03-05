package service

import (
	"context"

	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
)

type SavedService struct {
	repo repository.SavedRepo
}

func NewSavedService(repo repository.SavedRepo) *SavedService {
	return &SavedService{repo: repo}
}

func (s *SavedService) Save(ctx context.Context, data *model.Saved) error {
	return s.repo.Save(ctx, data)
}

func (s *SavedService) GetSaved(ctx context.Context, userID string) ([]model.Saved, error) {
	return s.repo.GetSaved(ctx, userID)
}

func (s *SavedService) Remove(ctx context.Context, ID string) error {
	return s.repo.Remove(ctx, ID)
}
