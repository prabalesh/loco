export interface ApiResponse<T> {
  data?: T
  message?: string
  error?: string
}

export interface ApiError {
  error: string
  fields?: Record<string, string>
}

export interface PaginationParams {
  page: number
  limit: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  totalPages: number
}
