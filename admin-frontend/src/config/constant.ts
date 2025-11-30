export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

export const ROUTES = {
  LOGIN: '/login',
  DASHBOARD: '/',
  USERS: '/users',
  ANALYTICS: '/analytics',
} as const
