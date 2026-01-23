export interface User {
  id: number
  email: string
  username: string
  role: string
  is_active: boolean
  email_verified: boolean
  created_at: string
  updated_at: string
}

export interface ExecutorConfig {
  docker_image: string
  memory_limit: number
  timeout: number
}

export interface Language {
  id: number
  language_id: string
  version: string
  name: string
  is_active: boolean
  extension: string
  default_template: string
  executor_config: ExecutorConfig
  created_at: Date
  updated_at: Date
}

export interface Tag {
  id: number
  name: string
  slug: string
  created_at: string
  updated_at: string
}

export interface Category {
  id: number
  name: string
  slug: string
  created_at: string
  updated_at: string
}

export interface Parameter {
  name: string
  type: string
  is_custom: boolean
}

export interface Problem {
  id: number
  title: string
  slug: string
  description: string
  difficulty: "easy" | "medium" | "hard"
  time_limit: number
  memory_limit: number
  current_step: 1 | 2 | 3 | 4
  validator_type: "exact_match"
  input_format: string
  output_format: string
  constraints: string
  status: "draft" | "published"
  validation_status: "draft" | "validated"
  is_active: boolean
  acceptance_rate: number
  total_submissions: number
  total_accepted: number

  // V2 Fields
  function_name?: string
  return_type?: string
  parameters?: Parameter[]
  test_cases?: TestCase[]
  validation_type?: string

  tags?: Tag[]
  categories?: Category[]
  created_at: Date
  updated_at: Date
}

export interface TestCase {
  id: number;
  problem_id: number;
  input: string;
  expected_output: string;
  is_hidden: boolean;
  is_sample: boolean;
  order: number;
  created_at: string;
  updated_at: string;
}

export interface TrendingProblem {
  id: number
  title: string
  slug: string
  submission_count: number
}

export interface LanguageStat {
  language_name: string
  count: number
}

export interface AdminAnalytics {
  total_users: number
  active_users: number
  inactive_users: number
  verified_users: number
  total_submissions: number
  pending_submissions: number
  active_workers: number
  queue_size: number
  submission_history?: { date: string; count: number }[]
  trending_problems?: TrendingProblem[]
  language_stats?: LanguageStat[]
}

export interface LoginCredentials {
  email: string
  password: string
}


export * from './request'
export * from './response'
export * from './problemLanguage'
