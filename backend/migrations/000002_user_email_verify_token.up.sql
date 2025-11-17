ALTER TABLE users 
	ADD COLUMN email_verification_token VARCHAR(64),
	ADD COLUMN email_verification_token_expires_at TIMESTAMP WITH TIME ZONE,
	ADD COLUMN email_verification_attempts INT DEFAULT 0,
	ADD COLUMN email_verification_last_sent_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_users_email_verification_token ON users(email_verification_token);