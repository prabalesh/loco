-- Create languages table (global language registry)
CREATE TABLE IF NOT EXISTS languages (
    id SERIAL PRIMARY KEY,
    language_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(50),
    extension VARCHAR(10),
    default_template TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    executor_config JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_languages_language_id ON languages(language_id);
CREATE INDEX idx_languages_is_active ON languages(is_active);

-- Updated_at trigger function (create once, reuse for all tables)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for languages table
CREATE TRIGGER update_languages_updated_at
    BEFORE UPDATE ON languages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
