-- Drop trigger first
DROP TRIGGER IF EXISTS trigger_update_problem_languages_updated_at ON problem_languages;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_problem_languages_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_problem_languages_is_validated;
DROP INDEX IF EXISTS idx_problem_languages_language_id;
DROP INDEX IF EXISTS idx_problem_languages_problem_id;

-- Drop table
DROP TABLE IF EXISTS problem_languages;
