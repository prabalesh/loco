import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { toast } from 'react-hot-toast'
import { authApi } from '../api/authApi'
import { authStore } from '../store/authStore'
import { ROUTES } from '../../../shared/constants/routes'

export const useLogout = () => {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      authStore.getState().logout()
      toast.success('Logged out successfully')
      navigate(ROUTES.LOGIN)
    },
    onError: () => {
      // Even if API fails, clear local state
      authStore.getState().logout()
      navigate(ROUTES.LOGIN)
    },
  })
}
