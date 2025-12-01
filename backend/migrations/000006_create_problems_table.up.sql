-- Create problems table
CREATE TABLE IF NOT EXISTS problems (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) UNIQUE NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    time_limit INTEGER DEFAULT 2000 CHECK (time_limit > 0),
    memory_limit INTEGER DEFAULT 256 CHECK (memory_limit > 0),
    validator_type VARCHAR(50) DEFAULT 'exact_match' CHECK (
        validator_type IN ('exact_match', 'ignore_whitespace', 'unordered_array', 'float_tolerance', 'permutation')
    ),
    input_format TEXT,
    output_format TEXT,
    constraints TEXT,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    visibility VARCHAR(20) DEFAULT 'public' CHECK (visibility IN ('public', 'private')),
    is_active BOOLEAN DEFAULT TRUE,
    acceptance_rate DECIMAL(5,2) DEFAULT 0.00 CHECK (acceptance_rate >= 0 AND acceptance_rate <= 100),
    total_submissions INTEGER DEFAULT 0 CHECK (total_submissions >= 0),
    total_accepted INTEGER DEFAULT 0 CHECK (total_accepted >= 0),
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_problems_slug ON problems(slug);
CREATE INDEX idx_problems_difficulty ON problems(difficulty);
CREATE INDEX idx_problems_status ON problems(status);
CREATE INDEX idx_problems_visibility ON problems(visibility);
CREATE INDEX idx_problems_is_active ON problems(is_active);
CREATE INDEX idx_problems_created_by ON problems(created_by);
CREATE INDEX idx_problems_acceptance_rate ON problems(acceptance_rate);

-- Updated_at trigger
CREATE TRIGGER update_problems_updated_at
    BEFORE UPDATE ON problems
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
