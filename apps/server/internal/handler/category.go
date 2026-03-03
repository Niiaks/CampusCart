package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
	"github.com/Niiaks/campusCart/pkg/types"
)

type CategoryHandler struct {
	Handler
	categoryService *service.CategoryService
}

func NewCategoryHandler(server *server.Server, categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		Handler:         NewHandler(server),
		categoryService: categoryService,
	}
}

const maxUploadSize = 2 << 20 // 2 MB

func (ch *CategoryHandler) Create() http.HandlerFunc {
	return Handle(ch.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) (*types.CategoryResponse, error) {
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			return nil, errs.NewBadRequestError(
				fmt.Sprintf("file too large or invalid form: %v", err), false, nil, nil, nil,
			)
		}

		name := strings.TrimSpace(r.FormValue("name"))
		if name == "" {
			return nil, errs.NewBadRequestError(
				"category name is required", true, nil,
				[]errs.FieldError{{Field: "name", Error: "is required"}}, nil,
			)
		}

		// 3. Extract the uploaded image file
		file, _, err := r.FormFile("image")
		if err != nil {
			return nil, errs.NewBadRequestError(
				"image file is required", true, nil,
				[]errs.FieldError{{Field: "image", Error: "is required"}}, nil,
			)
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, errs.NewBadRequestError("could not read file", false, nil, nil, nil)
		}

		reader := bytes.NewReader(fileBytes)

		contentType := http.DetectContentType(fileBytes)
		if !strings.HasPrefix(contentType, "image/") {
			return nil, errs.NewBadRequestError(
				"invalid file type: only images are allowed", true, nil,
				[]errs.FieldError{{Field: "image", Error: "must be an image file"}}, nil,
			)
		}

		category := &model.Category{
			Name: name,
		}

		if err := ch.categoryService.CreateCategory(r.Context(), category, reader); err != nil {
			return nil, err
		}

		return &types.CategoryResponse{
			ID:        category.ID,
			ParentID:  category.ParentID,
			Name:      category.Name,
			Slug:      category.Slug,
			Icon:      category.Icon,
			PublicID:  category.PublicID,
			IsActive:  category.IsActive,
			SortOrder: category.SortOrder,
		}, nil
	}, http.StatusCreated, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

func (ch *CategoryHandler) GetAll() http.HandlerFunc {
	return Handle(ch.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) ([]types.CategoryResponse, error) {
		categories, err := ch.categoryService.GetCategories(r.Context())
		if err != nil {
			return nil, err
		}

		resp := make([]types.CategoryResponse, len(categories))
		for i, c := range categories {
			resp[i] = types.CategoryResponse{
				ID:        c.ID,
				ParentID:  c.ParentID,
				Name:      c.Name,
				Slug:      c.Slug,
				Icon:      c.Icon,
				PublicID:  c.PublicID,
				IsActive:  c.IsActive,
				SortOrder: c.SortOrder,
			}
		}

		return resp, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

func (ch *CategoryHandler) GetByID() http.HandlerFunc {
	return Handle(ch.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) (*types.CategoryResponse, error) {
		id := chi.URLParam(r, "id")
		if id == "" {
			return nil, errs.NewBadRequestError("category ID is required", false, nil, nil, nil)
		}

		category, err := ch.categoryService.GetCategoryByID(r.Context(), id)
		if err != nil {
			return nil, err
		}

		return &types.CategoryResponse{
			ID:        category.ID,
			ParentID:  category.ParentID,
			Name:      category.Name,
			Slug:      category.Slug,
			Icon:      category.Icon,
			PublicID:  category.PublicID,
			IsActive:  category.IsActive,
			SortOrder: category.SortOrder,
		}, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

// GetAttributes returns dynamic attributes for a category (optionally merged with parents).
func (ch *CategoryHandler) GetAttributes() http.HandlerFunc {
	return Handle(ch.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) ([]types.CategoryAttributeResponse, error) {
		id := chi.URLParam(r, "id")
		if id == "" {
			return nil, errs.NewBadRequestError("category ID is required", false, nil, nil, nil)
		}

		includeParents := true
		if v := r.URL.Query().Get("include_parents"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				includeParents = b
			}
		}

		attrs, err := ch.categoryService.GetCategoryAttributes(r.Context(), id, includeParents)
		if err != nil {
			return nil, err
		}

		return attrs, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

func (ch *CategoryHandler) Update() http.HandlerFunc {
	return Handle(ch.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) (*types.CategoryResponse, error) {
		id := chi.URLParam(r, "id")
		if id == "" {
			return nil, errs.NewBadRequestError("category ID is required", false, nil, nil, nil)
		}

		if err := r.ParseMultipartForm(maxUploadSize); err != nil && err != http.ErrNotMultipart {
			return nil, errs.NewBadRequestError(fmt.Sprintf("file too large or invalid form: %v", err), false, nil, nil, nil)
		}

		update := &types.UpdateCategory{}

		name := strings.TrimSpace(r.FormValue("name"))
		if name != "" {
			update.Name = &name
		}

		parentID := strings.TrimSpace(r.FormValue("parent_id"))
		if parentID != "" {
			update.ParentID = &parentID
		}

		var reader *bytes.Reader
		file, _, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			fileBytes, err := io.ReadAll(file)
			if err != nil {
				return nil, errs.NewBadRequestError("could not read file", false, nil, nil, nil)
			}
			contentType := http.DetectContentType(fileBytes)
			if !strings.HasPrefix(contentType, "image/") {
				return nil, errs.NewBadRequestError("invalid file type", true, nil, nil, nil)
			}
			reader = bytes.NewReader(fileBytes)
		}

		var fileInterface interface{} = nil
		if reader != nil {
			fileInterface = reader
		}

		category, err := ch.categoryService.UpdateCategory(r.Context(), id, update, fileInterface)
		if err != nil {
			return nil, err
		}

		return &types.CategoryResponse{
			ID:        category.ID,
			ParentID:  category.ParentID,
			Name:      category.Name,
			Slug:      category.Slug,
			Icon:      category.Icon,
			PublicID:  category.PublicID,
			IsActive:  category.IsActive,
			SortOrder: category.SortOrder,
		}, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

func (ch *CategoryHandler) Delete() http.HandlerFunc {
	return HandleNoContent(ch.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) error {
		id := chi.URLParam(r, "id")
		if id == "" {
			return errs.NewBadRequestError("category ID is required", false, nil, nil, nil)
		}

		err := ch.categoryService.DeleteCategory(r.Context(), id)
		if err != nil {
			return err
		}
		return nil
	}, http.StatusNoContent, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}
