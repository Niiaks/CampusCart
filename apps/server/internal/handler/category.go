package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

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
			ID:       category.ID,
			Name:     category.Name,
			ImageUrl: category.ImageUrl,
			PublicID: category.PublicID,
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
				ID:       c.ID,
				Name:     c.Name,
				ImageUrl: c.ImageUrl,
				PublicID: c.PublicID,
			}
		}

		return resp, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}
