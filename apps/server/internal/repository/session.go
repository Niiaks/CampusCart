package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Niiaks/campusCart/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

const sessionDuration = 7 * 24 * time.Hour

type SessionRepo interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetUserBySession(ctx context.Context, sessionID string) (*model.User, error)
	DeleteSession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string) error
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		pool: pool,
	}
}

// CreateSession persists a new session row with computed expiry and returns the generated ID.
func (sr *SessionRepository) CreateSession(ctx context.Context, session *model.Session) error {

	expiresAt := time.Now().Add(sessionDuration)

	sql := `INSERT INTO sessions (user_id, ip_address, user_agent, refresh_token, expires_at) VALUES($1,$2,$3,$4,$5) returning id`

	err := sr.pool.QueryRow(ctx, sql, session.UserID, session.IPAddress, session.UserAgent, session.RefreshToken, expiresAt).Scan(&session.ID)

	if err != nil {
		return fmt.Errorf("create session error: %w", err)
	}

	return nil
}

// GetUserBySession fetches the user for a valid, non-expired refresh token.
func (sr *SessionRepository) GetUserBySession(ctx context.Context, token string) (*model.User, error) {

	sql := `SELECT u.id, u.email, u.role FROM sessions s INNER JOIN users u on s.user_id = u.id WHERE s.refresh_token = $1 AND s.expires_at > $2 AND s.deleted_at IS NULL`

	var user model.User

	row := sr.pool.QueryRow(ctx, sql, token, time.Now())

	err := row.Scan(&user.ID, &user.Email, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("session not found: %w", err)
		}
		return nil, err
	}
	return &user, nil
}

// DeleteSession removes a session by its id; returns an error if no row was deleted.
func (sr *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	sql := `UPDATE sessions SET deleted_at = $1 WHERE refresh_token = $2`

	result, err := sr.pool.Exec(ctx, sql, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("delete session error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// RefreshSession bumps last_activity and extends expiry for a valid refresh token.
func (sr *SessionRepository) RefreshSession(ctx context.Context, token string) error {
	now := time.Now()
	expiresAt := now.Add(sessionDuration)

	sql := `UPDATE sessions SET last_activity = $1, expires_at = $2 WHERE refresh_token = $3 AND expires_at > $4 AND deleted_at IS NULL`

	result, err := sr.pool.Exec(ctx, sql, now, expiresAt, token, now)
	if err != nil {
		return fmt.Errorf("refresh session error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found or expired")
	}

	return nil
}
