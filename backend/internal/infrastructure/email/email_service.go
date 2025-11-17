package email

import (
	"context"
	"fmt"

	"github.com/prabalesh/loco/backend/pkg/config"
	"github.com/resend/resend-go/v2"
	"go.uber.org/zap"
)

type EmailService struct {
	client *resend.Client
	config *config.Config
	logger *zap.Logger
}

func NewEmailService(cfg *config.Config, logger *zap.Logger) *EmailService {
	client := resend.NewClient(cfg.Email.ResendAPIKey)

	return &EmailService{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, email, username, token string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.config.Email.FromName, s.config.Email.FromEmail),
		To:      []string{email},
		Subject: "Verify Your Email - Loco Platform",
		Html:    s.getVerificationEmailHTML(s.config.Server.AppBaseUrl, username, token),
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		s.logger.Error("Failed to send verification email",
			zap.String("email", email),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Verification email sent",
		zap.String("email", email),
	)

	return nil
}

func (s *EmailService) getVerificationEmailHTML(appUrl, username, token string) string {
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", appUrl, token)

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6; 
            color: #333; 
            margin: 0;
            padding: 0;
        }
        .container { 
            max-width: 600px; 
            margin: 0 auto; 
            background: #ffffff;
        }
        .header { 
            background: linear-gradient(135deg, #3b82f6 0%%, #8b5cf6 100%%); 
            color: white; 
            padding: 40px 20px; 
            text-align: center; 
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .content { 
            padding: 40px 30px;
            background: #f9fafb;
        }
        .verification-link {
            display: inline-block;
            padding: 12px 24px;
            background: #3b82f6;
            color: white !important;
            text-decoration: none;
            border-radius: 6px;
            margin: 30px 0;
            font-weight: bold;
            font-size: 18px;
        }
        .warning {
            background: #fef3c7;
            border-left: 4px solid #f59e0b;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        .footer { 
            text-align: center; 
            padding: 30px 20px; 
            color: #6b7280; 
            font-size: 14px;
            border-top: 1px solid #e5e7eb;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Email Verification</h1>
        </div>
        <div class="content">
            <h2>Hello %s!</h2>
            <p>Thank you for registering with Loco Platform. To complete your registration, please verify your email by clicking the button below:</p>
            
            <a href="%s" class="verification-link" target="_blank" rel="noopener noreferrer">
                Verify My Email
            </a>
            
            <div class="warning">
                <strong>⚠️ Important:</strong>
                <ul style="margin: 10px 0; padding-left: 20px;">
                    <li>This link expires in <strong>%d hours</strong></li>
                    <li>It can only be used once</li>
                    <li>Don't share this link with anyone</li>
                </ul>
            </div>
            
            <p>If you didn't create an account with us, please ignore this email or contact support if you have concerns.</p>
        </div>
        <div class="footer">
            <p>© 2025 Loco Platform. All rights reserved.</p>
            <p style="font-size: 12px; color: #9ca3af;">This is an automated email, please do not reply.</p>
        </div>
    </div>
</body>
</html>
    `, username, verificationLink, int(s.config.Email.PasswordResetExpiryMinutes/60))
}

func (s *EmailService) SendPasswordResetEmail(ctx context.Context, email, username, token string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.config.Email.FromName, s.config.Email.FromEmail),
		To:      []string{email},
		Subject: "Reset Your Password - Loco Platform",
		Html:    s.getPasswordResetTokenEmailHTML(s.config.Server.AppBaseUrl, username, token),
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		s.logger.Error("Failed to send password reset email",
			zap.String("email", email),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Password reset email sent",
		zap.String("email", email),
	)

	return nil
}

func (s *EmailService) getPasswordResetTokenEmailHTML(appUrl, username, token string) string {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", appUrl, token)

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background-color: #f9fafb;
            color: #333;
            margin: 0; padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%);
            color: white;
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .content {
            padding: 30px 25px;
        }
        .greeting {
            font-size: 18px;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background: #f59e0b;
            color: white !important;
            text-decoration: none;
            font-weight: 700;
            padding: 15px 30px;
            border-radius: 8px;
            font-size: 18px;
            margin: 25px 0;
        }
        .instructions {
            font-size: 14px;
            color: #92400e;
            margin-bottom: 15px;
        }
        .footer {
            font-size: 12px;
            color: #9ca3af;
            text-align: center;
            padding: 20px 10px;
            border-top: 1px solid #e5e7eb;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p class="greeting">Hello %s,</p>
            <p class="instructions">You requested to reset your password for your Loco Platform account. Please click the button below to proceed:</p>
            <p style="text-align:center;">
                <a href="%s" class="button" target="_blank" rel="noopener noreferrer">Reset Password</a>
            </p>
            <p class="instructions">This link will expire in <strong>%d minutes</strong> and can only be used once. Please do not share this link with anyone.</p>
            <p>If you did not request this email, please ignore it or contact our support team.</p>
        </div>
        <div class="footer">
            &copy; 2025 Loco Platform. All rights reserved.<br/>
            This is an automated email, please do not reply.
        </div>
    </div>
</body>
</html>
    `, username, resetLink, s.config.Email.PasswordResetExpiryMinutes)
}
