package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
	"github.com/Niiaks/campusCart/pkg/types"
)

// ListingHandler wires HTTP endpoints to listing service.
type ListingHandler struct {
	Handler
	listingService *service.ListingService
}

func NewListingHandler(server *server.Server, listingService *service.ListingService) *ListingHandler {
	return &ListingHandler{
		Handler:        NewHandler(server),
		listingService: listingService,
	}
}

// Create a new listing (expects uploaded asset URLs from direct uploads).
func (lh *ListingHandler) Create() http.HandlerFunc {
	return Handle(lh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.CreateListingRequest) (*types.ListingResponse, error) {
		brandID := middleware.GetBrandID(r.Context())
		if brandID == "" {
			return nil, errs.NewUnauthorizedError("brand not found", false)
		}
		if req.CategoryID == "" {
			return nil, errs.NewBadRequestError("category_id is required", false, nil, nil, nil)
		}
		if len(req.ImageUrls) == 0 {
			return nil, errs.NewBadRequestError("at least one image is required", true, nil, nil, nil)
		}

		listing := &model.Listing{
			BrandID:     brandID,
			CategoryID:  req.CategoryID,
			Title:       req.Title,
			Description: req.Description,
			Price:       req.Price,
			Condition:   req.Condition,
			Negotiable:  req.Negotiable,
			Attributes:  req.Attributes,
			ImageUrls:   req.ImageUrls,
			VideoUrls:   req.VideoUrls,
			IsActive:    req.IsActive,
			IsPromoted:  false,
		}

		if err := lh.listingService.CreateListing(r.Context(), listing); err != nil {
			return nil, err
		}

		return toListingResponse(listing), nil
	}, http.StatusCreated, func() *types.CreateListingRequest { return &types.CreateListingRequest{} })
}

// Get (with view increment) a listing by ID.
func (lh *ListingHandler) Get() http.HandlerFunc {
	return Handle(lh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) (*types.ListingResponse, error) {
		id := chi.URLParam(r, "id")
		if id == "" {
			return nil, errs.NewBadRequestError("listing id is required", false, nil, nil, nil)
		}

		listing, err := lh.listingService.ViewListing(r.Context(), id)
		if err != nil {
			return nil, err
		}

		return toListingResponse(listing), nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

// List listings with simple filters.
func (lh *ListingHandler) List() http.HandlerFunc {
	return Handle(lh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) ([]types.ListingResponse, error) {
		filter := types.ListingFilter{}

		q := r.URL.Query()
		if v := q.Get("category_id"); v != "" {
			filter.CategoryID = v
		}
		if v := q.Get("brand_id"); v != "" {
			filter.BrandID = v
		}
		if v := q.Get("brand"); v != "" {
			filter.BrandName = v
		} else if v := q.Get("brand_name"); v != "" { // fallback alias
			filter.BrandName = v
		}
		if v := q.Get("include_descendants"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				filter.IncludeDescendants = b
			}
		}
		if v := q.Get("search"); v != "" {
			filter.Search = v
		}
		if v := q.Get("condition"); v != "" {
			filter.Condition = v
		}
		if v := q.Get("min_price"); v != "" {
			if p, err := strconv.ParseInt(v, 10, 64); err == nil {
				filter.MinPrice = &p
			}
		}
		if v := q.Get("max_price"); v != "" {
			if p, err := strconv.ParseInt(v, 10, 64); err == nil {
				filter.MaxPrice = &p
			}
		}
		if v := q.Get("limit"); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				filter.Limit = p
			}
		}
		if v := q.Get("offset"); v != "" {
			if p, err := strconv.Atoi(v); err == nil {
				filter.Offset = p
			}
		}

		listings, err := lh.listingService.List(r.Context(), filter)
		if err != nil {
			return nil, err
		}

		resp := make([]types.ListingResponse, len(listings))
		for i := range listings {
			resp[i] = *toListingResponse(&listings[i])
		}
		return resp, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

// Update listing (partial JSON).
func (lh *ListingHandler) Update() http.HandlerFunc {
	return Handle(lh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.UpdateListing) (*types.ListingResponse, error) {
		id := chi.URLParam(r, "id")
		if id == "" {
			return nil, errs.NewBadRequestError("listing id is required", false, nil, nil, nil)
		}

		if err := lh.listingService.UpdateListing(r.Context(), id, req); err != nil {
			return nil, err
		}

		listing, err := lh.listingService.GetListing(r.Context(), id)
		if err != nil {
			return nil, err
		}

		return toListingResponse(listing), nil
	}, http.StatusOK, func() *types.UpdateListing { return &types.UpdateListing{} })
}

// Delete listing (requires brand ownership).
func (lh *ListingHandler) Delete() http.HandlerFunc {
	return HandleNoContent(lh.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) error {
		id := chi.URLParam(r, "id")
		if id == "" {
			return errs.NewBadRequestError("listing id is required", false, nil, nil, nil)
		}

		brandID := middleware.GetBrandID(r.Context())
		if brandID == "" {
			return errs.NewUnauthorizedError("brand not found", false)
		}

		return lh.listingService.DeleteListing(r.Context(), id, brandID)
	}, http.StatusNoContent, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

// UploadSignature returns signed params for direct-to-Cloudinary uploads.
type uploadSignatureRequest struct {
	Folder       string `json:"folder"`
	ResourceType string `json:"resource_type"` // image | video | auto
}

func (r *uploadSignatureRequest) Validate() error { return nil }

type uploadSignatureResponse struct {
	UploadURL string            `json:"upload_url"`
	Params    map[string]string `json:"params"`
}

func (lh *ListingHandler) UploadSignature() http.HandlerFunc {
	return Handle(lh.Handler, func(w http.ResponseWriter, r *http.Request, req *uploadSignatureRequest) (*uploadSignatureResponse, error) {
		folder := req.Folder
		if folder == "" {
			folder = "listing"
		}

		payload, err := lh.listingService.GenerateDirectUpload(r.Context(), folder, req.ResourceType)
		if err != nil {
			return nil, err
		}

		return &uploadSignatureResponse{
			UploadURL: payload.UploadURL,
			Params:    payload.Params,
		}, nil
	}, http.StatusOK, func() *uploadSignatureRequest { return &uploadSignatureRequest{} })
}

func toListingResponse(l *model.Listing) *types.ListingResponse {
	return &types.ListingResponse{
		ID:          l.ID,
		BrandID:     l.BrandID,
		CategoryID:  l.CategoryID,
		Title:       l.Title,
		Description: l.Description,
		Price:       l.Price,
		Condition:   l.Condition,
		Negotiable:  l.Negotiable,
		Attributes:  l.Attributes,
		ImageUrls:   l.ImageUrls,
		VideoUrls:   l.VideoUrls,
		IsActive:    l.IsActive,
		IsPromoted:  l.IsPromoted,
		ViewsCount:  l.ViewsCount,
		CreatedAt:   l.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:   l.UpdatedAt.Format(http.TimeFormat),
	}
}
