package service

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/pkg/types"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepo
	file         *file.Client
}

func NewCategoryService(categoryRepo repository.CategoryRepo, file *file.Client) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		file:         file,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, category *model.Category, file interface{}) error {
	url, id, err := cs.file.UploadImage(ctx, file, "category")
	if err != nil {
		return err
	}

	category.PublicID = id
	category.Icon = url
	if category.Slug == "" {
		category.Slug = slugify(category.Name)
	}
	if category.IsActive == false {
		category.IsActive = true
	}

	err = cs.categoryRepo.CreateCategory(ctx, category)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CategoryService) GetCategories(ctx context.Context) ([]model.Category, error) {
	return cs.categoryRepo.GetCategories(ctx)
}

func (cs *CategoryService) GetCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	return cs.categoryRepo.GetCategory(ctx, id)
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, id string, update *types.UpdateCategory, file interface{}) (*model.Category, error) {
	existingCategory, err := cs.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	if file != nil {
		url, publicID, err := cs.file.UploadImage(ctx, file, "category")
		if err != nil {
			return nil, err
		}

		update.Icon = &url
		update.PublicID = &publicID

		if existingCategory.PublicID != "" {
			_ = cs.file.DeleteFile(ctx, existingCategory.PublicID, "image")
		}
	}

	if update.Name != nil && (update.Slug == nil || *update.Slug == "") {
		newSlug := slugify(*update.Name)
		update.Slug = &newSlug
	}

	err = cs.categoryRepo.UpdateCategory(ctx, id, update)
	if err != nil {
		return nil, err
	}

	return cs.categoryRepo.GetCategory(ctx, id)
}

func (cs *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	cat, err := cs.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		return err
	}

	err = cs.file.DeleteFile(ctx, cat.PublicID,
		"image")
	if err != nil {
		return err
	}
	return cs.categoryRepo.DeleteCategory(ctx, id)
}

// GetCategoryAttributes returns attributes for a category, optionally merged with parent categories.
func (cs *CategoryService) GetCategoryAttributes(ctx context.Context, categoryID string, includeParents bool) ([]types.CategoryAttributeResponse, error) {
	attrs, err := cs.categoryRepo.GetCategoryAttributes(ctx, categoryID, includeParents)
	if err != nil {
		return nil, err
	}

	resp := make([]types.CategoryAttributeResponse, 0, len(attrs))
	for _, a := range attrs {
		var opts []string
		if len(a.OptionsRaw) > 0 {
			_ = json.Unmarshal(a.OptionsRaw, &opts)
		}
		resp = append(resp, types.CategoryAttributeResponse{
			Name:      a.Name,
			Label:     a.Label,
			Type:      a.Type,
			Options:   opts,
			Required:  a.Required,
			SortOrder: a.SortOrder,
		})
	}

	return resp, nil
}

var slugRegexp = regexp.MustCompile(`[^a-z0-9-]+`)

// slugify creates a URL-friendly slug from a name.
func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRegexp.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "category"
	}
	return s
}
