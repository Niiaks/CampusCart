package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Niiaks/campusCart/pkg/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BrandRepository struct {
	pool *pgxpool.Pool
}

func NewBrandRepository(pool *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{
		pool: pool,
	}
}

func (br *BrandRepository) UpdateBrand(ctx context.Context, brandID string, brandUpdate *types.UpdateBrand) error {

	sql := `UPDATE brands SET name = COALESCE($1,name), description = COALESCE($2,description), profile_url = COALESCE($3,profile_url), banner_url = COALESCE($4,banner_url), updated_at = $5 WHERE id = $6`

	_, err := br.pool.Exec(ctx, sql, brandUpdate.Name, brandUpdate.Description, brandUpdate.ProfileUrl, brandUpdate.BannerUrl, time.Now(), brandID)

	if err != nil {
		return fmt.Errorf("error updating brand %s", err)
	}
	return nil
}
