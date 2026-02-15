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

type UserRepository struct {
	pool *pgxpool.Pool
}

type UserRepo interface {
	InsertUser(ctx context.Context, user *model.User) error
	SelectUser(ctx context.Context, userID string) (*types.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	VerifyUserEmail(ctx context.Context, email string) error
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

	userSql := `INSERT INTO users(username,email,password,phone,email_verification_code,email_verification_expires_at) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	brandSql := `INSERT INTO brands(name,seller_id) VALUES($1,$2)`

	err = txn.QueryRow(ctx, userSql, user.Username, user.Email, user.Password, user.Phone, user.EmailVerificationCode, user.EmailVerificationCodeExpiresAt).Scan(&user.ID)
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

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	sql := `SELECT id,username,email,password,phone,email_verified,email_verification_code,email_verification_expires_at FROM users WHERE email = $1`

	var user model.User
	err := ur.pool.QueryRow(ctx, sql, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Phone, &user.EmailVerified, &user.EmailVerificationCode, &user.EmailVerificationCodeExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) VerifyUserEmail(ctx context.Context, email string) error {
	sql := `UPDATE users SET email_verified = TRUE, email_verification_code = NULL, email_verification_expires_at = NULL, updated_at = $1 WHERE email = $2`
	_, err := ur.pool.Exec(ctx, sql, time.Now(), email)
	return err
}
