package service

import (
	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/lib/job"
	"github.com/Niiaks/campusCart/internal/repository"
)

type Service struct {
	Auth     *AuthService
	Job      *job.JobService
	Category *CategoryService
	Listing  *ListingService
	Brand    *BrandService
}

func NewServices(repo *repository.Repository, job *job.JobService, file *file.Client) *Service {
	authService := NewAuthService(repo.User, repo.Session, job)
	categoryService := NewCategoryService(repo.Category, file)
	listingService := NewListingService(repo.Listing, file)
	brandService := NewBrandService(repo.Brand, file)
	return &Service{
		Auth:     authService,
		Job:      job,
		Category: categoryService,
		Listing:  listingService,
		Brand:    brandService,
	}
}
