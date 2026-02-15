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

func (sr *SessionRepository) CreateSession(ctx context.Context, session *model.Session) error {

	expiresAt := time.Now().Add(sessionDuration)

	sql := `INSERT INTO sessions (user_id, expires_at) VALUES($1,$2) returning id`

	err := sr.pool.QueryRow(ctx, sql, session.UserID, expiresAt).Scan(&session.ID)

	if err != nil {
		return fmt.Errorf("Create session error %s", err)
	}

	return nil
}

func (sr *SessionRepository) GetUserBySession(ctx context.Context, sessionID string) (*model.User, error) {

	sql := `SELECT u.id, u.email, u.role FROM sessions s INNER JOIN users u on s.user_id = u.id WHERE s.id = $1 AND expires_at > $2`

	var user model.User

	row := sr.pool.QueryRow(ctx, sql, sessionID, time.Now())

	err := row.Scan(&user.ID, &user.Email, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Session not found %s", err)
		}
		return nil, err
	}
	return &user, nil
}

func (sr *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	sql := `DELETE FROM sessions WHERE id = $1`

	result, err := sr.pool.Exec(ctx, sql, sessionID)
	if err != nil {
		return fmt.Errorf("delete session error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (sr *SessionRepository) RefreshSession(ctx context.Context, sessionID string) error {
	now := time.Now()
	expiresAt := now.Add(sessionDuration)

	sql := `UPDATE sessions SET last_activity = $1, expires_at = $2, updated_at = $1 WHERE id = $3 AND expires_at > $4`

	result, err := sr.pool.Exec(ctx, sql, now, expiresAt, sessionID, now)
	if err != nil {
		return fmt.Errorf("refresh session error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found or expired")
	}

	return nil
}
