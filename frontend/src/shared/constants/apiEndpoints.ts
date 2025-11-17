export const API_ENDPOINTS = {
  AUTH: {
    REGISTER: '/auth/register',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me',
    VERIFY_EMAIL: '/auth/verify-email',
    RESEND_VERIFICATION: '/auth/resend-verification',
    FORGOT_PASSWORD: '/auth/forgot-password',
    RESET_PASSWORD: '/auth/reset-password',
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
