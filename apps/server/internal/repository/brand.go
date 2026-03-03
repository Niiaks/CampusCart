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

type BrandRepo interface {
	GetBrandIDBySeller(ctx context.Context, sellerID string) (string, error)
	GetBrandByID(ctx context.Context, brandID string) (*model.Brand, error)
	UpdateBrand(ctx context.Context, brandID string, brandUpdate *types.UpdateBrand) error
}

type BrandRepository struct {
	pool *pgxpool.Pool
}

func NewBrandRepository(pool *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{
		pool: pool,
	}
}

// GetBrandIDBySeller returns the brand id for a given seller (user) if present.
func (br *BrandRepository) GetBrandIDBySeller(ctx context.Context, sellerID string) (string, error) {
	sql := `SELECT id FROM brands WHERE seller_id = $1 AND deleted_at IS NULL LIMIT 1`
	var id string
	err := br.pool.QueryRow(ctx, sql, sellerID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("brand not found for seller: %w", err)
	}
	return id, nil
}

// GetBrandByID returns the brand record for a given brand ID.
func (br *BrandRepository) GetBrandByID(ctx context.Context, brandID string) (*model.Brand, error) {
	sql := `
		SELECT id, seller_id, name, slug, description, profile_url, banner_url, is_verified, created_at, updated_at
		FROM brands
		WHERE id = $1 AND deleted_at IS NULL
	`

	var brand model.Brand
	err := br.pool.QueryRow(ctx, sql, brandID).Scan(
		&brand.ID,
		&brand.SellerID,
		&brand.Name,
		&brand.Slug,
		&brand.Description,
		&brand.ProfileUrl,
		&brand.BannerUrl,
		&brand.IsVerified,
		&brand.CreatedAt,
		&brand.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("brand not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching brand: %w", err)
	}

	return &brand, nil
}

func (br *BrandRepository) UpdateBrand(ctx context.Context, brandID string, brandUpdate *types.UpdateBrand) error {

	sql := `UPDATE brands SET name = COALESCE($1,name), description = COALESCE($2,description), profile_url = COALESCE($3,profile_url), banner_url = COALESCE($4,banner_url), updated_at = $5 WHERE id = $6`

	_, err := br.pool.Exec(ctx, sql, brandUpdate.Name, brandUpdate.Description, brandUpdate.ProfileUrl, brandUpdate.BannerUrl, time.Now(), brandID)

	if err != nil {
		return fmt.Errorf("error updating brand: %w", err)
	}
	return nil
}
