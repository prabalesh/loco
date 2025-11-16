import { useQuery } from '@tanstack/react-query'
import { authApi } from '../api/authApi'
import { useAuth } from '@/shared/hooks/useAuth'

export const useProfile = (userId?: number | 'me') => {
  const { user } = useAuth()
  const targetUserId = userId || 'me'

  return useQuery({
    queryKey: ['user-profile', targetUserId],
    queryFn: () => authApi.getUserProfile(targetUserId),
    enabled: !!user, // Only fetch if logged in
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}
