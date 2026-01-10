-- migrations/000xxx_create_problem_languages_table.up.sql

-- Create problem_languages junction table
CREATE TABLE IF NOT EXISTS problem_languages (
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language_id INTEGER NOT NULL REFERENCES languages(id) ON DELETE CASCADE,
    function_code TEXT NOT NULL,
    main_code TEXT NOT NULL,
    solution_code TEXT NOT NULL,
    is_validated BOOLEAN DEFAULT FALSE,
    validated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (problem_id, language_id)
);

-- Indexes
CREATE INDEX idx_problem_languages_problem_id ON problem_languages(problem_id);
CREATE INDEX idx_problem_languages_language_id ON problem_languages(language_id);
CREATE INDEX idx_problem_languages_is_validated ON problem_languages(problem_id, is_validated);

-- Trigger function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_problem_languages_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update updated_at
CREATE TRIGGER trigger_update_problem_languages_updated_at
BEFORE UPDATE ON problem_languages
FOR EACH ROW
EXECUTE FUNCTION update_problem_languages_updated_at();
