package job

import (
	"context"
	"encoding/json"
	"fmt"

	lib "github.com/Niiaks/campusCart/internal/lib/email"
	"github.com/hibiken/asynq"
)

const (
	TaskEmailWelcome      = "email:welcome"
	TaskEmailVerification = "email:verification"
)

type EmailWelcomePayload struct {
	To       string `json:"to"`
	Username string `json:"username"`
}

type EmailVerificationPayload struct {
	To       string `json:"to"`
	Username string `json:"username"`
	Code     string `json:"code"`
}

func NewEmailWelcomeTask(to, username string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailWelcomePayload{To: to, Username: username})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email:welcome payload: %w", err)
	}
	return asynq.NewTask(TaskEmailWelcome, payload), nil
}

func NewEmailVerificationTask(to, username, code string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailVerificationPayload{To: to, Username: username, Code: code})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email:verification payload: %w", err)
	}
	return asynq.NewTask(TaskEmailVerification, payload), nil
}

func (j *JobService) HandleEmailWelcome(emailClient *lib.Client) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p EmailWelcomePayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("failed to unmarshal email:welcome payload: %w", err)
		}

		j.logger.Info().
			Str("task", TaskEmailWelcome).
			Str("to", p.To).
			Str("username", p.Username).
			Msg("sending welcome email")

		if err := emailClient.SendWelcomeEmail(p.To, p.Username); err != nil {
			j.logger.Error().Err(err).
				Str("task", TaskEmailWelcome).
				Str("to", p.To).
				Msg("failed to send welcome email")
			return fmt.Errorf("failed to send welcome email: %w", err)
		}

		j.logger.Info().
			Str("task", TaskEmailWelcome).
			Str("to", p.To).
			Msg("welcome email sent successfully")
		return nil
	}
}

func (j *JobService) HandleEmailVerification(emailClient *lib.Client) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p EmailVerificationPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("failed to unmarshal email:verification payload: %w", err)
		}

		j.logger.Info().
			Str("task", TaskEmailVerification).
			Str("to", p.To).
			Str("username", p.Username).
			Msg("sending verification email")

		if err := emailClient.SendEmailVerificationCode(p.To, p.Username, p.Code); err != nil {
			j.logger.Error().Err(err).
				Str("task", TaskEmailVerification).
				Str("to", p.To).
				Msg("failed to send verification email")
			return fmt.Errorf("failed to send verification email: %w", err)
		}

		j.logger.Info().
			Str("task", TaskEmailVerification).
			Str("to", p.To).
			Msg("verification email sent successfully")
		return nil
	}
}
