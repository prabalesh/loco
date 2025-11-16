export const API_ENDPOINTS = {
  AUTH: {
    REGISTER: '/auth/register',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me',
  },
  USERS: {
    PROFILE: (id: 'me') => `/users/${id}`,
    ME: '/users/me',
    BY_USERNAME: (username: string) => `/users/${username}`
  },
  PROBLEMS: {
    LIST: '/problems',
    DETAIL: (slug: string) => `/problems/${slug}`,
    SUBMIT: '/problems/submit',
  },
  SUBMISSIONS: {
    LIST: '/submissions',
    DETAIL: (id: string) => `/submissions/${id}`,
  },
  LEADERBOARD: {
    GET: '/leaderboard',
  },
} as const
