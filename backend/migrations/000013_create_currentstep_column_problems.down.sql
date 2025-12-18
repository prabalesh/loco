-- Remove index first
DROP INDEX IF EXISTS idx_problems_current_step;

-- Remove current_step column from problems table
ALTER TABLE problems
DROP COLUMN IF EXISTS current_step;
