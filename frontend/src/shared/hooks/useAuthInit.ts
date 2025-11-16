import { useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { authApi } from '@/features/auth/api/authApi'
import { authStore } from '@/features/auth/store/authStore'

export const useAuthInit = () => {
  const { isAuthenticated } = authStore()

  const { data, error, isLoading } = useQuery({
    queryKey: ['auth-me'],
    queryFn: () => authApi.getCurrentUser(),
    enabled: isAuthenticated, // Only fetch if Zustand says we're logged in
    retry: false, // Don't retry on failure
    staleTime: Infinity, // Don't refetch automatically
  })

  useEffect(() => {
    if (data) {
      // Update Zustand with fresh user data from backend
      authStore.getState().setUser(data)
    } else if (error) {
      // Token expired or invalid, logout
      authStore.getState().logout()
    }
  }, [data, error])

  return { isLoading, isAuthenticated: !!data }
}
