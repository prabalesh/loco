// src/api/axios.ts
import axios from 'axios'
import { API_BASE_URL } from '../config/constant'

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Important for cookies
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor for handling token refresh
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // If 401 and not already retried, try to refresh token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        await axios.post(
          `${API_BASE_URL}/admin/auth/refresh`,
          {},
          { withCredentials: true }
        )
        return axiosInstance(originalRequest)
      } catch (refreshError) {
        // Refresh failed, redirect to login
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default axiosInstance
