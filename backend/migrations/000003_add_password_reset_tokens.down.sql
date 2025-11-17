DROP INDEX IF EXISTS idx_users_password_reset_token;

ALTER TABLE users
    DROP COLUMN IF EXISTS password_reset_token,
    DROP COLUMN IF EXISTS password_reset_token_expires_at,
    DROP COLUMN IF EXISTS password_reset_sent_at;