export interface UserStats {
  total_submissions: number
  accepted_submissions: number
  problems_solved: number
  acceptance_rate: number
}

export interface User {
  id: number
  email: string
  username: string
  role: string
  email_verified: boolean
  created_at: string
  stats?: UserStats
}

export interface PublicUser {
  id: number
  username: string
  is_verified: boolean
  created_at: string
  stats?: UserStats
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterResponse {
  data: {
    message: string
    user: User
  }
}

export interface LoginResponse {
  data: {
    message: string
    user: User
  }
}

export interface ValidationErrors {
  email?: string
  username?: string
  password?: string
}
