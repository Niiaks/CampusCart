package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/Niiaks/campusCart/internal/lib/job"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/pkg/types"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepo
	sessionRepo repository.SessionRepo
	jobService  *job.JobService
}

const emailDuration = 30 * time.Minute

func NewAuthService(userRepo repository.UserRepo, sessionRepo repository.SessionRepo, jobService *job.JobService) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jobService:  jobService,
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

	// check if email is verified before logging in
	if user.EmailVerified == false {
		return nil, errors.New("email not verified, verify to continue")
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

func (auth *AuthService) Register(ctx context.Context, request *types.RegisterUser) (*types.RegisterResponse, error) {
	existing, err := auth.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	if !isValidEmail(request.Email) {
		return nil, errors.New("invalid email. Student email required")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	emailVerificationCode := generateVerificationCode()
	now := time.Now()
	expiresAt := now.Add(emailDuration)

	user := &model.User{
		Username:                       request.Username,
		Email:                          request.Email,
		Password:                       string(hashed),
		Phone:                          request.Phone,
		EmailVerificationCode:          emailVerificationCode,
		EmailVerificationCodeExpiresAt: expiresAt,
	}

	if err := auth.userRepo.InsertUser(ctx, user); err != nil {
		return nil, err
	}

	// Enqueue verification email
	task, err := job.NewEmailVerificationTask(user.Email, user.Username, emailVerificationCode)
	if err != nil {
		return nil, fmt.Errorf("failed to create verification email task: %w", err)
	}
	// add a deduplication key to prevent same email from sending twice
	if err := auth.jobService.Enqueue(task, asynq.TaskID("verify:"+user.Email)); err != nil {
		return nil, fmt.Errorf("failed to enqueue verification email: %w", err)
	}

	return &types.RegisterResponse{
		Message: "registration successful, please check your email for a verification code",
	}, nil
}

func (auth *AuthService) VerifyEmail(ctx context.Context, request *types.VerifyEmailRequest) (*types.LoginResponse, error) {
	user, err := auth.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.EmailVerified {
		return nil, errors.New("email already verified")
	}

	if user.EmailVerificationCode != request.Code {
		return nil, errors.New("invalid verification code")
	}

	if time.Now().After(user.EmailVerificationCodeExpiresAt) {
		return nil, errors.New("verification code has expired")
	}

	// Mark email as verified
	if err := auth.userRepo.VerifyUserEmail(ctx, user.Email); err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	// Create session now that email is verified
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

	// Enqueue welcome email
	welcomeTask, err := job.NewEmailWelcomeTask(user.Email, user.Username)
	if err == nil {
		_ = auth.jobService.Enqueue(welcomeTask, asynq.TaskID("welcome:"+user.Email))
	}

	return &types.LoginResponse{
		SessionID: session.ID,
		User:      userResponse,
	}, nil
}

func (auth *AuthService) Logout(ctx context.Context, sessionID string) error {
	return auth.sessionRepo.DeleteSession(ctx, sessionID)
}

func (auth *AuthService) GetCurrentUser(ctx context.Context, userID string) (*types.UserResponse, error) {
	return auth.userRepo.SelectUser(ctx, userID)
}

func generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.IntN(1000000))
}

// check if email is valid student email by taking the part after @
// should match st.ug.edu.gh. any other is invalid
// could use student db api to check in the future if need be
func isValidEmail(email string) bool {
	parts := strings.Split(email, "@")

	if parts[1] == "st.ug.edu.gh" {
		return true
	}
	return false
}
