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
	config *config.EmailConfig
	logger *zap.Logger
}

func NewEmailService(cfg *config.EmailConfig, logger *zap.Logger) *EmailService {
	client := resend.NewClient(cfg.ResendAPIKey)

	return &EmailService{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, email, username, otp string) error {
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail),
		To:      []string{email},
		Subject: "Verify Your Email - Loco Platform",
		Html:    s.getVerificationEmailHTML(username, otp),
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

func (s *EmailService) getVerificationEmailHTML(username, otp string) string {
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
        .otp-box { 
            background: white; 
            border: 3px solid #3b82f6; 
            border-radius: 12px; 
            padding: 30px; 
            text-align: center; 
            margin: 30px 0;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .otp-code { 
            font-size: 42px; 
            font-weight: bold; 
            letter-spacing: 12px; 
            color: #3b82f6;
            font-family: 'Courier New', monospace;
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
        .button {
            display: inline-block;
            padding: 12px 24px;
            background: #3b82f6;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            margin: 10px 0;
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
            <p>Thank you for registering with Loco Platform. To complete your registration, please verify your email address using the code below:</p>
            
            <div class="otp-box">
                <p style="margin: 0 0 10px 0; font-size: 14px; color: #6b7280;">Your Verification Code</p>
                <div class="otp-code">%s</div>
            </div>
            
            <div class="warning">
                <strong>⚠️ Important:</strong>
                <ul style="margin: 10px 0; padding-left: 20px;">
                    <li>This code expires in <strong>%d minutes</strong></li>
                    <li>You have <strong>5 attempts</strong> to enter the correct code</li>
                    <li>Don't share this code with anyone</li>
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
    `, username, otp, s.config.OTPExpirationMinutes)
}
