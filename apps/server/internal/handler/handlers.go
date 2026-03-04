package handler

import (
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
)

type Handlers struct {
	Health   *HealthHandler
	OpenAPI  *OpenAPIHandler
	Auth     *AuthHandler
	Category *CategoryHandler
	Listing  *ListingHandler
	Brand    *BrandHandler
}

func NewHandlers(s *server.Server, authService *service.AuthService, categoryService *service.CategoryService, listingService *service.ListingService, brandService *service.BrandService) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(s),
		OpenAPI:  NewOpenAPIHandler(s),
		Auth:     NewAuthHandler(s, authService),
		Category: NewCategoryHandler(s, categoryService),
		Listing:  NewListingHandler(s, listingService),
		Brand:    NewBrandHandler(s, brandService),
	}
}
