package email

import (
	"fmt"
	"net/smtp"

	"github.com/coderz-space/coderz.space/internal/config"
)

type Service interface {
	SendPasswordResetEmail(to string, resetToken string) error
}

type smtpService struct {
	config *config.Config
}

func NewService(cfg *config.Config) Service {
	return &smtpService{config: cfg}
}

func (s *smtpService) SendPasswordResetEmail(to string, resetToken string) error {
	// If SMTP is not fully configured, log and return (development mode fallback)
	if s.config.SMTPHost == "" || s.config.SMTPPort == 0 || s.config.SMTPFrom == "" {
		fmt.Printf("SMTP not fully configured. Would have sent reset email to %s with token %s\n", to, resetToken)
		return nil
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.config.FrontendOrigin, resetToken)

	subject := "Password Reset Request"
	body := fmt.Sprintf("Hello,\n\nYou requested a password reset. Click the link below to reset your password:\n\n%s\n\nIf you did not request this, please ignore this email.\n", resetLink)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + s.config.SMTPFrom + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPass, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	return smtp.SendMail(addr, auth, s.config.SMTPFrom, []string{to}, msg)
}
