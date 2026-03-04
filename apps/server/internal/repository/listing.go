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

type ListingRepo interface {
	CreateListing(ctx context.Context, listing *model.Listing) error
	GetListingByID(ctx context.Context, id string) (*model.Listing, error)
	List(ctx context.Context, filter types.ListingFilter) ([]model.Listing, error)
	UpdateListing(ctx context.Context, id string, updateData *types.UpdateListing) error
	DeleteListing(ctx context.Context, id string, brandID string) error
	IncrementViews(ctx context.Context, id string) error
}

type ListingRepository struct {
	pool *pgxpool.Pool
}

func NewListingRepository(pool *pgxpool.Pool) *ListingRepository {
	return &ListingRepository{
		pool: pool,
	}
}

// CreateListing inserts a new listing into the database.
func (r *ListingRepository) CreateListing(ctx context.Context, listing *model.Listing) error {
	sql := `
		INSERT INTO listings (
			brand_id, category_id, title, description, price, condition, 
			negotiable, attributes, image_urls, video_urls, is_active
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, created_at, updated_at, views_count
	`
	err := r.pool.QueryRow(ctx, sql,
		listing.BrandID,
		listing.CategoryID,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.Condition,
		listing.Negotiable,
		listing.Attributes,
		listing.ImageUrls,
		listing.VideoUrls,
		listing.IsActive,
	).Scan(
		&listing.ID,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&listing.ViewsCount,
	)

	if err != nil {
		return fmt.Errorf("repository: failed to create listing: %w", err)
	}

	return nil
}

// GetListingByID retrieves a single active listing by ID.
func (r *ListingRepository) GetListingByID(ctx context.Context, id string) (*model.Listing, error) {
	sql := `
		SELECT 
			id, brand_id, category_id, title, description, price, condition, 
			negotiable, attributes, image_urls, video_urls, is_active, is_promoted, 
			views_count, created_at, updated_at
		FROM listings
		WHERE id = $1 AND deleted_at IS NULL
	`

	var listing model.Listing
	err := r.pool.QueryRow(ctx, sql, id).Scan(
		&listing.ID,
		&listing.BrandID,
		&listing.CategoryID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Condition,
		&listing.Negotiable,
		&listing.Attributes,
		&listing.ImageUrls,
		&listing.VideoUrls,
		&listing.IsActive,
		&listing.IsPromoted,
		&listing.ViewsCount,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository: listing not found: %w", err)
		}
		return nil, fmt.Errorf("repository: failed to get listing: %w", err)
	}

	return &listing, nil
}

// UpdateListing applies partial updates to a listing
func (r *ListingRepository) UpdateListing(ctx context.Context, id string, update *types.UpdateListing) error {
	sql := `
		UPDATE listings
		SET
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			category_id = COALESCE($3, category_id),
			price = COALESCE($4, price),
			condition = CAST(COALESCE($5, condition::text) AS listing_condition),
			negotiable = COALESCE($6, negotiable),
			attributes = COALESCE($7, attributes),
			image_urls = COALESCE($8, image_urls),
			video_urls = COALESCE($9, video_urls),
			is_active = COALESCE($10, is_active),
			is_promoted = COALESCE($11, is_promoted),
			updated_at = $12
		WHERE id = $13 AND deleted_at IS NULL
	`
	cmdTag, err := r.pool.Exec(ctx, sql,
		update.Title,
		update.Description,
		update.CategoryID,
		update.Price,
		update.Condition,
		update.Negotiable,
		update.Attributes,
		update.ImageUrls,
		update.VideoUrls,
		update.IsActive,
		update.IsPromoted,
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("repository: failed to update listing: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("repository: listing not found or no rows updated")
	}

	return nil
}

// DeleteListing marks a listing as deleted (soft delete).
func (r *ListingRepository) DeleteListing(ctx context.Context, id string, brandID string) error {
	sql := `
		UPDATE listings 
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND brand_id = $2 AND deleted_at IS NULL
	`
	cmdTag, err := r.pool.Exec(ctx, sql, id, brandID)
	if err != nil {
		return fmt.Errorf("repository: failed to delete listing: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("repository: listing not found or unauthorized deletion")
	}

	return nil
}

// IncrementViews atomically increments the views_count for a listing.
func (r *ListingRepository) IncrementViews(ctx context.Context, id string) error {
	sql := `
		UPDATE listings 
		SET views_count = views_count + 1 
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("repository: failed to increment listing views: %w", err)
	}
	return nil
}

// List returns a slice of active listings with optional filtering and pagination.
func (r *ListingRepository) List(ctx context.Context, filter types.ListingFilter) ([]model.Listing, error) {
	prefix := ""
	query := `
		SELECT 
			l.id, l.brand_id, l.category_id, l.title, l.description, l.price, l.condition, 
			l.negotiable, l.attributes, l.image_urls, l.video_urls, l.is_active, l.is_promoted, 
			l.views_count, l.created_at, l.updated_at
		FROM listings l
	`

	joins := ""
	where := "WHERE l.deleted_at IS NULL AND l.is_active = TRUE"
	args := []interface{}{}
	argIndex := 1

	if filter.IncludeDescendants && filter.CategoryID != "" {
		prefix = fmt.Sprintf(`
		WITH cat_tree AS (
			SELECT id FROM categories WHERE id = $%d
			UNION ALL
			SELECT c.id FROM categories c JOIN cat_tree ct ON c.parent_id = ct.id
		)
		`, argIndex)
		args = append(args, filter.CategoryID)
		argIndex++
	}

	if filter.BrandName != "" {
		joins += " JOIN brands b ON b.id = l.brand_id AND b.deleted_at IS NULL"
		where += fmt.Sprintf(" AND (b.name ILIKE $%d OR b.slug ILIKE $%d)", argIndex, argIndex+1)
		pattern := "%" + filter.BrandName + "%"
		args = append(args, pattern, pattern)
		argIndex += 2
	}

	if filter.CategoryID != "" {
		if filter.IncludeDescendants {
			where += " AND l.category_id IN (SELECT id FROM cat_tree)"
		} else {
			where += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
			args = append(args, filter.CategoryID)
			argIndex++
		}
	}

	if filter.BrandID != "" {
		where += fmt.Sprintf(" AND l.brand_id = $%d", argIndex)
		args = append(args, filter.BrandID)
		argIndex++
	}

	if filter.MinPrice != nil {
		where += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		args = append(args, *filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice != nil {
		where += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		args = append(args, *filter.MaxPrice)
		argIndex++
	}

	if filter.Condition != "" {
		where += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		args = append(args, filter.Condition)
		argIndex++
	}

	if filter.Search != "" {
		where += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	query = prefix + query + joins + " " + where + " ORDER BY l.is_promoted DESC, l.created_at DESC"

	limit := 20
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	query += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit)
	argIndex++

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list listings: %w", err)
	}
	defer rows.Close()

	var listings []model.Listing
	for rows.Next() {
		var l model.Listing
		if err := rows.Scan(
			&l.ID,
			&l.BrandID,
			&l.CategoryID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.Condition,
			&l.Negotiable,
			&l.Attributes,
			&l.ImageUrls,
			&l.VideoUrls,
			&l.IsActive,
			&l.IsPromoted,
			&l.ViewsCount,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: failed to scan listing row: %w", err)
		}
		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows iteration error: %w", err)
	}

	return listings, nil
}
