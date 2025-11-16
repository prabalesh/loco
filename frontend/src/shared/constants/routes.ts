export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  PROBLEMS: '/problems',
  PROBLEM_DETAIL: (slug: string) => `/problems/${slug}`,
  SUBMISSIONS: '/submissions',
  LEADERBOARD: '/leaderboard',
  PROFILE: '/profile',
  USER_PROFILE: (username: string) => `/users/${username}`,
  NOT_FOUND: '/404',
} as const
