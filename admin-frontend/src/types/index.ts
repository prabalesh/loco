// src/types/index.ts
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
