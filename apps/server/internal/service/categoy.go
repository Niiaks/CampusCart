package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
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

// func (cs *CategoryService) UpdateCategory(ctx context.Context, update *types.UpdateCategory) (model.Category, error) {

// }

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
