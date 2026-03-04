package service

import (
	"context"
	"fmt"

	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/pkg/types"
)

// ListingService wraps listing use-cases on top of the repository layer.
type ListingService struct {
	repo repository.ListingRepo
	file *file.Client
}

func NewListingService(repo repository.ListingRepo, fileClient *file.Client) *ListingService {
	return &ListingService{repo: repo, file: fileClient}
}

// CreateListing persists a new listing. Callers should set BrandID/CategoryID before invoking.
func (s *ListingService) CreateListing(ctx context.Context, listing *model.Listing) error {
	if listing == nil {
		return fmt.Errorf("listing is required")
	}
	if listing.BrandID == "" {
		return fmt.Errorf("brand_id is required")
	}
	if listing.CategoryID == "" {
		return fmt.Errorf("category_id is required")
	}
	return s.repo.CreateListing(ctx, listing)
}

// GetListing fetches a listing by ID without mutating counters.
func (s *ListingService) GetListing(ctx context.Context, id string) (*model.Listing, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return s.repo.GetListingByID(ctx, id)
}

// ViewListing fetches a listing and increments its views count.
func (s *ListingService) ViewListing(ctx context.Context, id string) (*model.Listing, error) {
	listing, err := s.GetListing(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.IncrementViews(ctx, id); err != nil {
		return nil, err
	}

	// Reflect the increment locally for immediate responses
	listing.ViewsCount++
	return listing, nil
}

// List returns listings filtered and paginated.
func (s *ListingService) List(ctx context.Context, filter types.ListingFilter) ([]model.Listing, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return s.repo.List(ctx, filter)
}

// UpdateListing applies partial updates to a listing by ID.
func (s *ListingService) UpdateListing(ctx context.Context, id string, update *types.UpdateListing) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if update == nil {
		return fmt.Errorf("update payload is required")
	}
	return s.repo.UpdateListing(ctx, id, update)
}

// DeleteListing soft-deletes a listing, verifying ownership via brand ID.
func (s *ListingService) DeleteListing(ctx context.Context, id string, brandID string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if brandID == "" {
		return fmt.Errorf("brand_id is required")
	}
	return s.repo.DeleteListing(ctx, id, brandID)
}

// GenerateDirectUpload proxies to the file client for presigned upload params.
func (s *ListingService) GenerateDirectUpload(ctx context.Context, folder string, resourceType string) (*file.DirectUploadPayload, error) {
	return s.file.GenerateDirectUpload(ctx, folder, resourceType)
}
