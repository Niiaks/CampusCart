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
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategoryAttributes(ctx context.Context, categoryID string, includeParents bool) ([]model.CategoryAttribute, error)
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

	cmdTag, err := cr.pool.Exec(ctx, sql, updateCategory.Name, updateCategory.Slug, updateCategory.Icon, updateCategory.PublicID, updateCategory.ParentID, updateCategory.IsActive, updateCategory.SortOrder, time.Now(), categoryID)

	if err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("category not found or no rows updated")
	}
	return nil
}

// GetCategoryAttributes returns attributes for a category, optionally including ancestors (child wins on name conflicts).
func (cr *CategoryRepository) GetCategoryAttributes(ctx context.Context, categoryID string, includeParents bool) ([]model.CategoryAttribute, error) {
	if categoryID == "" {
		return nil, fmt.Errorf("categoryID is required")
	}

	var rows pgx.Rows
	var err error

	if includeParents {
		query := `
			WITH RECURSIVE cat AS (
				SELECT id, parent_id, 0 AS depth FROM categories WHERE id = $1
				UNION ALL
				SELECT c.id, c.parent_id, cat.depth + 1 FROM categories c JOIN cat ON c.id = cat.parent_id
			), attrs AS (
				SELECT ca.*, cat.depth
				FROM category_attributes ca
				JOIN cat ON ca.category_id = cat.id
			), dedup AS (
				SELECT DISTINCT ON (name)
					id, category_id, name, label, type, options, required, sort_order, created_at
				FROM attrs
				ORDER BY name, depth ASC, sort_order ASC
			)
			SELECT id, category_id, name, label, type, options, required, sort_order, created_at
			FROM dedup
			ORDER BY sort_order, name
		`
		rows, err = cr.pool.Query(ctx, query, categoryID)
	} else {
		query := `
			SELECT id, category_id, name, label, type, options, required, sort_order, created_at
			FROM category_attributes
			WHERE category_id = $1
			ORDER BY sort_order, name
		`
		rows, err = cr.pool.Query(ctx, query, categoryID)
	}

	if err != nil {
		return nil, fmt.Errorf("error fetching category attributes: %w", err)
	}
	defer rows.Close()

	var attrs []model.CategoryAttribute
	for rows.Next() {
		var a model.CategoryAttribute
		var optionsRaw []byte
		if err := rows.Scan(&a.ID, &a.CategoryID, &a.Name, &a.Label, &a.Type, &optionsRaw, &a.Required, &a.SortOrder, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning category attribute: %w", err)
		}
		a.OptionsRaw = optionsRaw
		attrs = append(attrs, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attrs, nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, categoryID string) error {
	sql := `UPDATE categories SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	cmdTag, err := cr.pool.Exec(ctx, sql, time.Now(), categoryID)
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("category not found or already deleted")
	}
	return nil
}
