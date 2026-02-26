package service

import (
	"context"

	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/middleware"
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

	logger := middleware.GetLogger(ctx)

	logger.Info().Msgf("the url is %s", url)
	logger.Info().Msgf("the id is %s", id)

	category.PublicID = id
	category.ImageUrl = url

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
