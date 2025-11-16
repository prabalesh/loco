import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { toast } from 'react-hot-toast'
import { authApi } from '../api/authApi'
import { authStore } from '../store/authStore'
import type { RegisterRequest } from '../../../shared/types/auth.types'
import { ROUTES } from '../../../shared/constants/routes'
import { AxiosError } from 'axios'
import type { ApiError } from '../../../shared/types/common.types'

export const useRegister = () => {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: (response) => {
      authStore.getState().setUser(response.user)
      toast.success('Registration successful! Welcome!')
      navigate(ROUTES.HOME)
    },
    onError: (error: AxiosError<ApiError>) => {
      const errorMessage = error.response?.data?.error || 'Registration failed'
      
      if (error.response?.data?.fields) {
        const fields = error.response.data.fields
        Object.values(fields).forEach((msg) => toast.error(msg))
      } else {
        toast.error(errorMessage)
      }
    },
  })
}
