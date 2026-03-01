package model

import "time"

// User represents a user entity in the database, encapsulating authentication,
// verification, and profile information. It includes fields for unique identification,
// credentials, contact details, verification codes and their expiration, status flags,
// and activity tracking. Sensitive fields such as password and verification codes are
// excluded from JSON serialization for security. Validation tags ensure data integrity
// for required fields. Embedded Model provides common metadata like timestamps.
type User struct {
	ID                             string     `json:"id"`                                                 // unique user id
	Username                       string     `json:"username" validate:"required,min=3,max=30,alphanum"` // public handle
	Email                          string     `json:"email" validate:"required,email"`                    // login email
	Password                       string     `json:"-" validate:"required,min=8"`                        // hashed password
	Phone                          string     `json:"phone,omitempty" validate:"required,min=10"`         // user phone
	Role                           string     `json:"role"`                                               // role e.g. user|admin
	EmailVerificationCode          *string    `json:"-"`                                                  // one-time email code
	PhoneVerificationCode          *string    `json:"-"`                                                  // one-time phone code
	EmailVerificationCodeExpiresAt *time.Time `json:"-"`                                                  // email code expiry
	PhoneVerificationCodeExpiresAt *time.Time `json:"-"`                                                  // phone code expiry
	EmailVerified                  bool       `json:"email_verified"`                                     // email confirmed flag
	PhoneVerified                  bool       `json:"phone_verified"`                                     // phone confirmed flag
	IsActive                       bool       `json:"is_active"`                                          // allowed to sign in
	IsBanned                       bool       `json:"is_banned"`                                          // globally blocked
	LastActive                     time.Time  `json:"last_active"`                                        // last seen timestamp
	Model
}
