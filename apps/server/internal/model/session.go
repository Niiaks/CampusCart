package model

import "time"

// Session represents an authentication session for a user.
// It encapsulates session metadata and lifecycle information used for issuing,
// validating and managing refresh tokens and active sessions.
type Session struct {
	ID     string `json:"id"`                          // server-side session id
	UserID string `json:"user_id" validate:"required"` // owning user id

	// RefreshToken holds the hashed refresh token and is not marshaled to JSON; it
	// must be stored and compared securely
	RefreshToken string `json:"refresh_token" validate:"required"` // hashed refresh token value

	// UserAgent and IPAddress capture client context (browser/device and source IP)
	// that can be used for logging, auditing or additional verification.
	UserAgent string `json:"user_agent" validate:"required"` // client user agent string
	IPAddress string `json:"ip_address" validate:"required"` // client IP at issuance

	// ExpiresAt denotes when the session becomes invalid and should be revoked/cleaned up.
	ExpiresAt time.Time `json:"expires_at"` // absolute expiry of session

	// LastActivity is updated on use to support idle timeouts or activity tracking.
	LastActivity time.Time `json:"last_activity"` // last touch timestamp
	Model
}
