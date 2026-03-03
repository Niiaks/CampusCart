package handler

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
	"github.com/Niiaks/campusCart/pkg/types"
)

type BrandHandler struct {
	Handler
	brandService *service.BrandService
}

func NewBrandHandler(server *server.Server, brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{
		Handler:      NewHandler(server),
		brandService: brandService,
	}
}

const brandMaxUploadSize = 5 << 20 // 5 MB

// GetCurrent returns the authenticated seller's brand profile.
func (bh *BrandHandler) GetCurrent() http.HandlerFunc {
	return Handle(bh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) (*types.BrandResponse, error) {
		brandID := middleware.GetBrandID(r.Context())
		if brandID == "" {
			return nil, errs.NewUnauthorizedError("brand not found", false)
		}

		brand, err := bh.brandService.GetBrand(r.Context(), brandID)
		if err != nil {
			return nil, err
		}

		return toBrandResponse(brand), nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

// Update allows sellers to update brand profile fields and optional images.
func (bh *BrandHandler) Update() http.HandlerFunc {
	return Handle(bh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.UpdateBrand) (*types.BrandResponse, error) {
		brandID := middleware.GetBrandID(r.Context())
		if brandID == "" {
			return nil, errs.NewUnauthorizedError("brand not found", false)
		}

		update := req
		var profileFile interface{}
		var bannerFile interface{}

		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "multipart/") {
			if err := r.ParseMultipartForm(brandMaxUploadSize); err != nil && err != http.ErrNotMultipart {
				return nil, errs.NewBadRequestError("invalid multipart form", false, nil, nil, nil)
			}

			if _, ok := r.Form["name"]; ok {
				name := strings.TrimSpace(r.FormValue("name"))
				if name == "" {
					return nil, errs.NewBadRequestError("name cannot be empty", true, nil, []errs.FieldError{{Field: "name", Error: "cannot be empty"}}, nil)
				}
				update.Name = &name
			}

			if _, ok := r.Form["description"]; ok {
				desc := strings.TrimSpace(r.FormValue("description"))
				update.Description = &desc
			}

			if file, _, err := r.FormFile("profile_image"); err == nil {
				defer file.Close()
				fileBytes, err := io.ReadAll(file)
				if err != nil {
					return nil, errs.NewBadRequestError("could not read profile_image", false, nil, nil, nil)
				}
				if !strings.HasPrefix(http.DetectContentType(fileBytes), "image/") {
					return nil, errs.NewBadRequestError("profile_image must be an image", true, nil, []errs.FieldError{{Field: "profile_image", Error: "must be an image"}}, nil)
				}
				profileFile = bytes.NewReader(fileBytes)
			} else if err != http.ErrMissingFile {
				return nil, errs.NewBadRequestError("invalid profile_image", false, nil, nil, nil)
			}

			if file, _, err := r.FormFile("banner_image"); err == nil {
				defer file.Close()
				fileBytes, err := io.ReadAll(file)
				if err != nil {
					return nil, errs.NewBadRequestError("could not read banner_image", false, nil, nil, nil)
				}
				if !strings.HasPrefix(http.DetectContentType(fileBytes), "image/") {
					return nil, errs.NewBadRequestError("banner_image must be an image", true, nil, []errs.FieldError{{Field: "banner_image", Error: "must be an image"}}, nil)
				}
				bannerFile = bytes.NewReader(fileBytes)
			} else if err != http.ErrMissingFile {
				return nil, errs.NewBadRequestError("invalid banner_image", false, nil, nil, nil)
			}
		} else {
			if update.Name != nil {
				trimmed := strings.TrimSpace(*update.Name)
				if trimmed == "" {
					return nil, errs.NewBadRequestError("name cannot be empty", true, nil, []errs.FieldError{{Field: "name", Error: "cannot be empty"}}, nil)
				}
				update.Name = &trimmed
			}
			if update.Description != nil {
				trimmed := strings.TrimSpace(*update.Description)
				update.Description = &trimmed
			}
		}

		brand, err := bh.brandService.UpdateBrand(r.Context(), brandID, update, profileFile, bannerFile)
		if err != nil {
			return nil, err
		}

		return toBrandResponse(brand), nil
	}, http.StatusOK, func() *types.UpdateBrand { return &types.UpdateBrand{} })
}

func toBrandResponse(b *model.Brand) *types.BrandResponse {
	if b == nil {
		return nil
	}

	var descPtr *string
	if strings.TrimSpace(b.Description) != "" {
		desc := b.Description
		descPtr = &desc
	}

	var profilePtr *string
	if strings.TrimSpace(b.ProfileUrl) != "" {
		p := b.ProfileUrl
		profilePtr = &p
	}

	var bannerPtr *string
	if strings.TrimSpace(b.BannerUrl) != "" {
		p := b.BannerUrl
		bannerPtr = &p
	}

	return &types.BrandResponse{
		ID:          b.ID,
		SellerID:    b.SellerID,
		Name:        b.Name,
		Slug:        b.Slug,
		Description: descPtr,
		ProfileUrl:  profilePtr,
		BannerUrl:   bannerPtr,
		IsVerified:  b.IsVerified,
		CreatedAt:   b.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:   b.UpdatedAt.Format(http.TimeFormat),
	}
}
