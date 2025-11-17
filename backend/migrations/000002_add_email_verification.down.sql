ALTER TABLE users 
DROP COLUMN email_verification_token,
DROP COLUMN email_verification_token_expires_at,
DROP COLUMN email_verification_attempts,
DROP COLUMN email_verification_last_sent_at;

DROP INDEX IF EXISTS idx_users_email_verification_token;