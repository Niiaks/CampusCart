package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/pkg/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

type UserRepo interface {
	InsertUser(ctx context.Context, user *model.User) error
	SelectUser(ctx context.Context, userID string) (*types.UserResponse, error)
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

// InsertUser creates a user and also creates a default brand with name for that user
// if they register.
func (ur *UserRepository) InsertUser(ctx context.Context, user *model.User) error {
	txn, err := ur.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txn.Rollback(ctx)

	userSql := `INSERT INTO users(username,email,password,phone) VALUES($1,$2,$3,$4) RETURNING id`
	brandSql := `INSERT INTO brands(name,seller_id) VALUES($1,$2)`

	err = txn.QueryRow(ctx, userSql, user.Username, user.Email, user.Password, user.Phone).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	_, err = txn.Exec(ctx, brandSql, user.Username, user.ID)
	if err != nil {
		return fmt.Errorf("error creating default brand: %w", err)
	}

	return txn.Commit(ctx)
}

func (ur *UserRepository) SelectUser(ctx context.Context, userID string) (*types.UserResponse, error) {
	sql := `SELECT id,username,email,phone,email_verified,last_active,is_active,created_at FROM users WHERE id = $1`

	var user types.UserResponse
	err := ur.pool.QueryRow(ctx, sql, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Phone, &user.EmailVerified, &user.LastActive, &user.IsActive, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found %s", err)
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	sql := `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int

	err := ur.pool.QueryRow(ctx, sql, email).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
