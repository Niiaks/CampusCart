package model

import "time"

// Session defines all fields for the session table in db.
type Session struct {
	ID           string    `json:"id" validate:"required"`
	UserID       string    `json:"user_id" validate:"required"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActivity time.Time `json:"last_activity"`
	Model
}
