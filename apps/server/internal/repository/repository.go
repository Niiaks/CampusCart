package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Session  *SessionRepository
	User     *UserRepository
	Brand    *BrandRepository
	Category *CategoryRepository
	Listing  *ListingRepository
	Saved    *SavedRepository
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Session:  NewSessionRepository(pool),
		User:     NewUserRepository(pool),
		Brand:    NewBrandRepository(pool),
		Category: NewCategoryRepository(pool),
		Listing:  NewListingRepository(pool),
		Saved:    NewSavedRepository(pool),
	}
}
