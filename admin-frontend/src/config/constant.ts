export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const ROUTES = {
  LOGIN: '/login',
  DASHBOARD: '/',
  USERS: '/users',
  ANALYTICS: '/analytics',
  PROBLEMS: {
    HOME: '/problems',
    CREATE: '/problems/create',
    TESTCASES: (id:number) => `/problems/${id}/testcases`,
    LANGUAGES: (id:number) => `/problems/${id}/languages`,
    VALIDATE: (id:number) => `/problems/${id}/validate`
  }
} as const

