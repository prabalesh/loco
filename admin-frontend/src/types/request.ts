export interface CreateOrUpdateLanguageRequest {
    language_id: String
    name: String,
    version: String,
    extension: String,
    default_template: String
}

export interface CreateOrUpdateProblemRequest {
  title: string;
  slug: string;
  description: string;
  difficulty: "easy" | "medium" | "hard";
  time_limit: number;
  memory_limit: number;
  validator_type: "exact_match";
  input_format: string;
  output_format: string;
  constraints: string;
  status: "draft" | "published";
  is_active: boolean;
}

export interface CreateTestCaseRequest {
  problem_id: number;
  input: string;
  expected_output: string;
  is_hidden?: boolean;
  is_sample?: boolean;
  order?: number;
}