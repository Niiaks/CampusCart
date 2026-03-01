package types

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (r *RegisterUser) Validate() error {
	return validate.Struct(r)
}

func (l *LoginUser) Validate() error {
	return validate.Struct(l)
}

type RegisterUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10"`
}

type LoginUser struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

func (v *VerifyEmailRequest) Validate() error {
	return validate.Struct(v)
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	RefreshToken string        `json:"-"`
	User         *UserResponse `json:"user"`
}

type UserResponse struct {
	ID            string    `json:"id"`
	Username      string    `json:"username" validate:"required,min=3"`
	Email         string    `json:"email" validate:"required,email"`
	Phone         string    `json:"phone" validate:"required,min=10"`
	Role          string    `json:"role" validate:"required,oneof=user admin"`
	EmailVerified bool      `json:"email_verified,omitempty"`
	PhoneVerified bool      `json:"phone_verified,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	LastActive    time.Time `json:"last_active"`
	CreatedAt     time.Time `json:"created_at"`
}

// EmptyRequest is used for endpoints that don't require a request body.
type EmptyRequest struct{}

func (e *EmptyRequest) Validate() error {
	return nil
}
