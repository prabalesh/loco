import { useQuery } from '@tanstack/react-query'
import { authApi } from '../api/authApi'

export const useUserProfile = (username: string) => {
  return useQuery({
    queryKey: ['user-profile-username', username],
    queryFn: () => authApi.getUserByUsername(username),
    enabled: !!username,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}
