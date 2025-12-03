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
  name: string
  is_active: boolean
  extension: string
  default_template: string
  executor_config: ExecutorConfig
  created_at: Date
  updated_at: Date
}

export interface AdminAnalytics {
  total_users: number
  active_users: number
  inactive_users: number
  verified_users: number
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface ApiError {
  error: string
}
