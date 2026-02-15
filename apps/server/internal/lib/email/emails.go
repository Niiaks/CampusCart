package lib

func (c *Client) SendWelcomeEmail(to, username string) error {
	data := map[string]string{
		"username": username,
	}

	return c.SendEmail(to, "Welcome to campusCart!", WelcomeTemplate, data)
}

func (c *Client) SendEmailVerificationCode(to, username, code string) error {
	data := map[string]string{
		"username": username,
		"code":     code,
	}

	return c.SendEmail(to, "Verify your email — campusCart", VerificationTemplate, data)
}
