package service

import (
	"context"
	"fmt"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/pkg/types"
)

// BrandService handles brand profile operations for sellers.
type BrandService struct {
	repo repository.BrandRepo
	file *file.Client
}

func NewBrandService(repo repository.BrandRepo, fileClient *file.Client) *BrandService {
	return &BrandService{repo: repo, file: fileClient}
}

// GetBrand fetches a brand by ID.
func (s *BrandService) GetBrand(ctx context.Context, brandID string) (*model.Brand, error) {
	if brandID == "" {
		return nil, fmt.Errorf("brand_id is required")
	}
	return s.repo.GetBrandByID(ctx, brandID)
}

// UpdateBrand applies partial updates and uploads optional profile/banner images.
func (s *BrandService) UpdateBrand(ctx context.Context, brandID string, update *types.UpdateBrand, profileFile interface{}, bannerFile interface{}) (*model.Brand, error) {
	if brandID == "" {
		return nil, fmt.Errorf("brand_id is required")
	}
	if update == nil {
		return nil, fmt.Errorf("update payload is required")
	}

	changed := false

	if profileFile != nil {
		url, _, err := s.file.UploadImage(ctx, profileFile, fmt.Sprintf("brand/%s/profile", brandID))
		if err != nil {
			return nil, err
		}
		update.ProfileUrl = &url
		changed = true
	}

	if bannerFile != nil {
		url, _, err := s.file.UploadImage(ctx, bannerFile, fmt.Sprintf("brand/%s/banner", brandID))
		if err != nil {
			return nil, err
		}
		update.BannerUrl = &url
		changed = true
	}

	if update.Name != nil || update.Description != nil || update.ProfileUrl != nil || update.BannerUrl != nil {
		changed = true
	}

	if !changed {
		return nil, errs.NewBadRequestError("no fields to update", true, nil, nil, nil)
	}

	if err := s.repo.UpdateBrand(ctx, brandID, update); err != nil {
		return nil, err
	}

	return s.repo.GetBrandByID(ctx, brandID)
}
