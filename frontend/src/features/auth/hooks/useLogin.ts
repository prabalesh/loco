import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { toast } from 'react-hot-toast'
import { authApi } from '../api/authApi'
import { authStore } from '../store/authStore'
import type { LoginRequest } from '../../../shared/types/auth.types'
import { ROUTES } from '@/shared/constants/routes'

export const useLogin = () => {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (response) => {
      authStore.getState().setUser(response.user)
      toast.success(`Welcome back, ${response.user.username}!`)
      navigate(ROUTES.HOME)
    }
  })
}
