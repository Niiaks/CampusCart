package service

import "errors"

// Domain-level sentinel errors returned by services.
// Handlers inspect these to build HTTP-specific responses (codes, field errors, actions).

// Auth errors
var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrEmailNotVerified        = errors.New("email not verified")
	ErrEmailAlreadyExists      = errors.New("email already in use")
	ErrInvalidStudentEmail     = errors.New("invalid email, student email required")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailAlreadyVerified    = errors.New("email already verified")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrVerificationCodeExpired = errors.New("verification code expired")
)

// Brand errors
var (
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)
