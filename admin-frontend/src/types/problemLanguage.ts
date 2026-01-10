// types/problemLanguage.ts or add to types/index.ts

export type ProblemLanguage = {
  id: number;
  problem_id: number;
  language_id: number;
  language_name: string;
  language_version?: string;
  function_code: string;
  main_code: string;
  solution_code: string;
  is_validated: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateProblemLanguageRequest = {
  language_id: number;
  function_code: string;
  main_code: string;
  solution_code: string;
};

export type UpdateProblemLanguageRequest = {
  function_code: string;
  main_code: string;
  solution_code: string;
};

export type ValidationResult = {
  is_validated: boolean;
  test_results: TestResult[];
  updated_at: string;
};

export type TestResult = {
  test_case: number;
  passed: boolean;
  execution_time?: string;
  memory_used?: string;
  error?: string;
  actual_output?: string;
  expected_output?: string;
};
