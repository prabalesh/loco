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

  getUserProfile: async (userId: 'me'): Promise<User> => {
    const response = await apiClient.get<User>(API_ENDPOINTS.USERS.PROFILE(userId))
    return response.data
  },

  getUserByUsername: async (username: string): Promise<PublicUser> => {
    const response = await apiClient.get<PublicUser>(API_ENDPOINTS.USERS.BY_USERNAME(username))
    return response.data
  },

  verifyEmail: async (email: string, otp: string): Promise<{ message: string }> => {
    const response = await apiClient.post<{ message: string }>(
      API_ENDPOINTS.AUTH.VERIFY_EMAIL,
      { email, otp }
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
}
