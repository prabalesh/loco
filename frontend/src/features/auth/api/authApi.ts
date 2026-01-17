import { apiClient } from '../../../shared/lib/axios'
import { API_ENDPOINTS } from '../../../shared/constants/apiEndpoints'
import type {
  RegisterRequest,
  LoginRequest,
  RegisterResponse,
  LoginResponse,
  User,
  PublicUser,
} from '../../../shared/types/auth.types'
import type { ApiResponse } from '@/shared/types/common.types'

export const authApi = {
  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await apiClient.post<RegisterResponse>(
      API_ENDPOINTS.AUTH.REGISTER,
      data
    )
    return response.data
  },

  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>(
      API_ENDPOINTS.AUTH.LOGIN,
      data
    )
    return response.data
  },

  logout: async (): Promise<void> => {
    await apiClient.post(API_ENDPOINTS.AUTH.LOGOUT)
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await apiClient.get<User>(API_ENDPOINTS.AUTH.ME)
    return response.data
  },

  refreshToken: async (): Promise<void> => {
    await apiClient.post(API_ENDPOINTS.AUTH.REFRESH)
  },

  getUserProfile: async (userId: 'me'): Promise<User | undefined> => {
    const response = await apiClient.get<ApiResponse<User>>(API_ENDPOINTS.USERS.PROFILE(userId))
    return response.data.data
  },

  getUserByUsername: async (username: string): Promise<PublicUser> => {
    const response = await apiClient.get<PublicUser>(API_ENDPOINTS.USERS.BY_USERNAME(username))
    return response.data
  },

  verifyEmail: async (token: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      API_ENDPOINTS.AUTH.VERIFY_EMAIL,
      { token }
    )
    return response.data
  },

  resendVerificationEmail: async (email: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      API_ENDPOINTS.AUTH.RESEND_VERIFICATION,
      { email }
    )
    return response.data
  },

  forgotPassword: async (email: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      API_ENDPOINTS.AUTH.FORGOT_PASSWORD,
      { email }
    )
    return response.data
  },

  resetPassword: async (token: string, newPassword: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      API_ENDPOINTS.AUTH.RESET_PASSWORD,
      { token, new_password: newPassword }  // adjust key names to match backend
    )
    return response.data
  },
}
