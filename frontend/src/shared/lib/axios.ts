import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { CONFIG } from '../constants/config'
import { API_ENDPOINTS } from '../constants/apiEndpoints'
import { authStore } from '@/features/auth/store/authStore'

export const apiClient = axios.create({
  baseURL: CONFIG.API_BASE_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
})

let isRefreshing = false

// Request interceptor
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Response interceptor
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    // Not a 401 or no config
    if (!originalRequest || error.response?.status !== 401) {
      return Promise.reject(error)
    }

    // ⭐ IMPORTANT: Don't retry auth endpoints (login, register, refresh)
    const authEndpoints = [
      API_ENDPOINTS.AUTH.LOGIN,
      API_ENDPOINTS.AUTH.REGISTER,
      API_ENDPOINTS.AUTH.REFRESH,
      API_ENDPOINTS.AUTH.LOGOUT,
    ]

    if (authEndpoints.some(endpoint => originalRequest.url?.includes(endpoint))) {
      console.log('Auth endpoint failed, not retrying:', originalRequest.url)
      
      // If refresh endpoint fails, logout
      if (originalRequest.url?.includes(API_ENDPOINTS.AUTH.REFRESH)) {
        authStore.getState().logout()
      }
      
      return Promise.reject(error)
    }

    // Already retried or currently refreshing
    if (originalRequest._retry || isRefreshing) {
      return Promise.reject(error)
    }

    originalRequest._retry = true
    isRefreshing = true

    try {
      console.log('Token expired, attempting refresh...')
      await apiClient.post(API_ENDPOINTS.AUTH.REFRESH)
      isRefreshing = false
      console.log('Token refreshed successfully, retrying original request')
      return apiClient(originalRequest)
    } catch (refreshError) {
      console.log('Token refresh failed, logging out')
      isRefreshing = false
      authStore.getState().logout()
      return Promise.reject(refreshError)
    }
  }
)
