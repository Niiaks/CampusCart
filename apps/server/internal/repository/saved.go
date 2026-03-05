package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Niiaks/campusCart/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SavedRepo interface {
	Save(ctx context.Context, data *model.Saved) error
	GetSaved(ctx context.Context, userID string) ([]model.Saved, error)
	Remove(ctx context.Context, id string) error
}
type SavedRepository struct {
	pool *pgxpool.Pool
}

func NewSavedRepository(pool *pgxpool.Pool) *SavedRepository {
	return &SavedRepository{pool: pool}
}

func (sr *SavedRepository) Save(ctx context.Context, data *model.Saved) error {
	sql := `INSERT INTO saved(user_id,listing_id) VALUES($1,$2) RETURNING id`

	err := sr.pool.QueryRow(ctx, sql, data.UserID, data.ListingID).Scan(&data.ID)

	if err != nil {
		return fmt.Errorf("error saving listing: %w", err)
	}
	return nil
}

func (sr *SavedRepository) GetSaved(ctx context.Context, userID string) ([]model.Saved, error) {
	sql := `
		SELECT
			s.id, s.user_id, s.listing_id, s.created_at,
			l.id, l.brand_id, l.category_id, l.title, l.description,
			l.price, l.condition, l.negotiable, l.attributes,
			l.image_urls, l.video_urls, l.is_active, l.is_promoted,
			l.views_count, l.created_at, l.updated_at
		FROM saved s
		JOIN listings l ON l.id = s.listing_id
		WHERE s.user_id = $1
		  AND l.deleted_at IS NULL
		ORDER BY s.created_at DESC`

	rows, err := sr.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching saved listings: %w", err)
	}
	defer rows.Close()

	var results []model.Saved
	for rows.Next() {
		var s model.Saved
		var l model.Listing
		var rawAttrs []byte

		if err := rows.Scan(
			&s.ID, &s.UserID, &s.ListingID, &s.CreatedAt,
			&l.ID, &l.BrandID, &l.CategoryID, &l.Title, &l.Description,
			&l.Price, &l.Condition, &l.Negotiable, &rawAttrs,
			&l.ImageUrls, &l.VideoUrls, &l.IsActive, &l.IsPromoted,
			&l.ViewsCount, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning saved listing: %w", err)
		}

		if len(rawAttrs) > 0 {
			_ = json.Unmarshal(rawAttrs, &l.Attributes)
		}

		s.Listing = &l
		results = append(results, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (sr *SavedRepository) Remove(ctx context.Context, ID string) error {
	sql := `
		UPDATE saved 
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
		`
	cmdTag, err := sr.pool.Exec(ctx, sql, ID)
	if err != nil {
		return fmt.Errorf("repository: failed to remove saved: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("repository: saved listing not found or unauthorized deletion")
	}

	return nil
}
