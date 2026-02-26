package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

type CategoryRepo interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategories(ctx context.Context) ([]model.Category, error)
	GetCategory(ctx context.Context, categoryID string) (*model.Category, error)
	UpdateCategory(ctx context.Context, categoryID string, updateCategory *types.UpdateCategory) error
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		pool: pool,
	}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, category *model.Category) error {
	sql := `INSERT INTO categories(name,image_url,public_id) VALUES($1,$2,$3) RETURNING id`

	err := cr.pool.QueryRow(ctx, sql, category.Name, category.ImageUrl, category.PublicID).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	return nil
}

func (cr *CategoryRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	sql := `SELECT id,name,image_url,public_id FROM categories ORDER BY name`

	rows, err := cr.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category

	for rows.Next() {
		var category model.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ImageUrl,
			&category.PublicID,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (cr *CategoryRepository) GetCategory(ctx context.Context, categoryID string) (*model.Category, error) {
	sql := `SELECT id,name,image_url,public_id FROM categories WHERE id =$1 ORDER BY name`

	var category model.Category

	err := cr.pool.QueryRow(ctx, sql, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.ImageUrl,
		&category.PublicID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, err
	}
	return &category, nil
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, categoryID string, updateCategory *types.UpdateCategory) error {
	sql := `UPDATE categories SET name = COALESCE($1,name), image_url = COALESCE($2,image_url), updated_at = $3 WHERE id = $4`

	_, err := cr.pool.Exec(ctx, sql, updateCategory.Name, updateCategory.ImageUrl, time.Now(), categoryID)

	if err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	return nil
}
