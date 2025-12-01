-- Create problem_languages junction table
CREATE TABLE IF NOT EXISTS problem_languages (
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language_id INTEGER NOT NULL REFERENCES languages(id) ON DELETE CASCADE,
    custom_template TEXT,
    is_enabled BOOLEAN DEFAULT TRUE,
    order_index INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (problem_id, language_id)
);

-- Indexes
CREATE INDEX idx_problem_languages_problem_id ON problem_languages(problem_id);
CREATE INDEX idx_problem_languages_language_id ON problem_languages(language_id);
CREATE INDEX idx_problem_languages_is_enabled ON problem_languages(problem_id, is_enabled);
