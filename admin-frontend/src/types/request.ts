export interface CreateOrUpdateLanguageRequest {
  language_id: string
  name: string
  version: string
  extension: string
  default_template: string
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
  tag_ids: number[];
  category_ids: number[];
}

export interface CreateTestCaseRequest {
  problem_id: number;
  input: string;
  expected_output: string;
  is_hidden?: boolean;
  is_sample?: boolean;
  order?: number;
}