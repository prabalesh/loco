-- Create user_problem_stats table (tracks user progress per problem)
CREATE TABLE IF NOT EXISTS user_problem_stats (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'unsolved' CHECK (status IN ('solved', 'attempted', 'unsolved')),
    attempts INTEGER DEFAULT 0 CHECK (attempts >= 0),
    first_solved_at TIMESTAMP,
    best_submission_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, problem_id)
);

-- Indexes
CREATE INDEX idx_user_problem_stats_user_id ON user_problem_stats(user_id);
CREATE INDEX idx_user_problem_stats_problem_id ON user_problem_stats(problem_id);
CREATE INDEX idx_user_problem_stats_status ON user_problem_stats(user_id, status);

-- Updated_at trigger
CREATE TRIGGER update_user_problem_stats_updated_at
    BEFORE UPDATE ON user_problem_stats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
