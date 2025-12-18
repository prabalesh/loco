-- Add current_step column to problems table
ALTER TABLE problems
ADD COLUMN current_step INTEGER NOT NULL DEFAULT 1
CHECK (current_step BETWEEN 1 AND 4);

-- Create index for faster queries on current_step
CREATE INDEX idx_problems_current_step ON problems(current_step);
