package lib

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Niiaks/campusCart/internal/config"
	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
)

type Client struct {
	client *resend.Client
	logger *zerolog.Logger
}

func NewClient(cfg *config.IntegrationConfig, logger *zerolog.Logger) *Client {
	return &Client{
		client: resend.NewClient(cfg.ResendApiKey),
		logger: logger,
	}
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {

	tmplPath := fmt.Sprintf("%s/%s.html", "templates/emails", templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse email template %s", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	params := &resend.SendEmailRequest{
		From: fmt.Sprintf("%s <%s>", "CampusCart", "onboarding@resend.dev"),
		// for testing purposes --FIX when domain available
		To:      []string{"juniorpappoe@gmail.com"},
		Subject: subject,
		Html:    body.String(),
	}

	sent, err := c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	c.logger.Info().Any("resend_email", sent).Msg("email sent successfully")
	return nil
}
