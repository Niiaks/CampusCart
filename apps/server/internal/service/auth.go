package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepo
	sessionRepo repository.SessionRepo
}

func NewAuthService(userRepo repository.UserRepo, sessionRepo repository.SessionRepo) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (auth *AuthService) Login(ctx context.Context, request *types.LoginUser) (*types.LoginResponse, error) {
	user, err := auth.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	session := &model.Session{
		UserID: user.ID,
	}

	err = auth.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	// Fetch the full user response (without sensitive fields)
	userResponse, err := auth.userRepo.SelectUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		SessionID: session.ID,
		User:      userResponse,
	}, nil
}

func (auth *AuthService) Register(ctx context.Context, request *types.RegisterUser) (*types.LoginResponse, error) {
	existing, err := auth.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashed),
		Phone:    request.Phone,
	}

	if err := auth.userRepo.InsertUser(ctx, user); err != nil {
		return nil, err
	}

	session := &model.Session{
		UserID: user.ID,
	}
	if err := auth.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	userResponse, err := auth.userRepo.SelectUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		SessionID: session.ID,
		User:      userResponse,
	}, nil
}

func (auth *AuthService) Logout(ctx context.Context, sessionID string) error {
	return auth.sessionRepo.DeleteSession(ctx, sessionID)
}
