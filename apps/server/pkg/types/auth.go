package types

import "time"

type RegisterUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10"`
}

type UserResponse struct {
	ID            string    `json:"id"`
	Username      string    `json:"username" validate:"required,min=3"`
	Email         string    `json:"email" validate:"required,email"`
	Phone         string    `json:"phone" validate:"required,min=10"`
	Role          string    `json:"role" validate:"required,oneof=seller buyer admin"`
	EmailVerified bool      `json:"email_verified,omitempty"`
	PhoneVerified bool      `json:"phone_verified,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	LastActive    time.Time `json:"last_active"`
	CreatedAt     time.Time `json:"created_at"`
}
