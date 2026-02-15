package service

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	errs "github.com/Niiaks/campusCart/internal/err"
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
		return nil, errs.NewUnauthorizedError("invalid credentials", false)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errs.NewUnauthorizedError("invalid credentials", false)
	}

	// check if email is verified before logging in
	if user.EmailVerified == false {
		return nil, errs.NewUnauthorizedError("email not verified, verify to continue", false)
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
		return nil, errs.NewBadRequestError("email already in use", false, nil, nil, nil)
	}

	if !isValidEmail(request.Email) {
		return nil, errs.NewBadRequestError("invalid email. Student email required", false, nil, nil, nil)
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
		EmailVerificationCode:          &emailVerificationCode,
		EmailVerificationCodeExpiresAt: &expiresAt,
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
		return nil, errs.NewBadRequestError("user not found", false, nil, nil, nil)
	}

	if user.EmailVerified {
		return nil, errs.NewBadRequestError("email already verified", false, nil, nil, nil)
	}

	if user.EmailVerificationCode == nil || *user.EmailVerificationCode != request.Code {
		return nil, errs.NewBadRequestError("invalid verification code", false, nil, nil, nil)
	}

	if user.EmailVerificationCodeExpiresAt == nil || time.Now().After(*user.EmailVerificationCodeExpiresAt) {
		return nil, errs.NewBadRequestError("verification code has expired", false, nil, nil, nil)
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
