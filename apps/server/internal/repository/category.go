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
	sql := `INSERT INTO categories(parent_id, name, slug, icon, public_id, is_active, sort_order) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`

	err := cr.pool.QueryRow(ctx, sql, category.ParentID, category.Name, category.Slug, category.Icon, category.PublicID, category.IsActive, category.SortOrder).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	return nil
}

func (cr *CategoryRepository) GetCategories(ctx context.Context) ([]model.Category, error) {
	sql := `SELECT id, parent_id, name, slug, icon, public_id, is_active, sort_order FROM categories WHERE deleted_at IS NULL ORDER BY sort_order, name`

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
			&category.ParentID,
			&category.Name,
			&category.Slug,
			&category.Icon,
			&category.PublicID,
			&category.IsActive,
			&category.SortOrder,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (cr *CategoryRepository) GetCategory(ctx context.Context, categoryID string) (*model.Category, error) {
	sql := `SELECT id, parent_id, name, slug, icon, public_id, is_active, sort_order FROM categories WHERE id =$1 AND deleted_at IS NULL`

	var category model.Category

	err := cr.pool.QueryRow(ctx, sql, categoryID).Scan(
		&category.ID,
		&category.ParentID,
		&category.Name,
		&category.Slug,
		&category.Icon,
		&category.PublicID,
		&category.IsActive,
		&category.SortOrder,
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
	sql := `UPDATE categories
SET
	name = COALESCE($1, name),
	slug = COALESCE($2, slug),
	icon = COALESCE($3, icon),
	public_id = COALESCE($4, public_id),
	parent_id = COALESCE($5, parent_id),
	is_active = COALESCE($6, is_active),
	sort_order = COALESCE($7, sort_order),
	updated_at = $8
WHERE id = $9 AND deleted_at IS NULL`

	_, err := cr.pool.Exec(ctx, sql, updateCategory.Name, updateCategory.Slug, updateCategory.Icon, updateCategory.PublicID, updateCategory.ParentID, updateCategory.IsActive, updateCategory.SortOrder, time.Now(), categoryID)

	if err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	return nil
}
