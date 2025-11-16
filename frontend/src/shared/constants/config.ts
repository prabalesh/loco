export const CONFIG = {
  API_BASE_URL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  APP_NAME: 'Loco',
  APP_DESCRIPTION: 'Master coding challenges and climb the leaderboard',
} as const
