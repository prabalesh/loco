import { authStore } from '../../features/auth/store/authStore'

export const useAuth = () => {
  const { user, isAuthenticated } = authStore()

  return {
    user,
    isAuthenticated,
    isAdmin: user?.role === 'admin',
    isUser: user?.role === 'user',
  }
}
