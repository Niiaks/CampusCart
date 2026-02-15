package model

import "time"

// User defines all fields for the user table in db.
type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username" validate:"required,min=3"`
	Email         string    `json:"email" validate:"required,email"`
	Password      string    `json:"password" validate:"required,min=8"`
	Phone         string    `json:"phone" validate:"required,min=10"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified,omitempty"`
	PhoneVerified bool      `json:"phone_verified,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	LastActive    time.Time `json:"last_active"`
	Model
}
