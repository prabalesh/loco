// src/api/adminApi.ts
import axios from './axios'
import type { User, AdminAnalytics, LoginCredentials } from '../types'

export const adminAuthApi = {
  login: (credentials: LoginCredentials) =>
    axios.post('/admin/auth/login', credentials),

  logout: () => axios.post('/admin/auth/logout'),

  refreshToken: () => axios.post('/admin/auth/refresh'),

  getProfile: () => axios.get<User>('/admin/auth/me'),
}

export const adminUsersApi = {
  getAll: () => axios.get<User[]>('/admin/users'),

  getById: (id: number) => axios.get<User>(`/admin/users/${id}`),

  deleteUser: (id: number) => axios.delete(`/admin/users/${id}`),

  updateRole: (id: number, role: string) =>
    axios.patch(`/admin/users/${id}/role`, { role }),

  updateStatus: (id: number, isActive: boolean) =>
    axios.patch(`/admin/users/${id}/status`, { is_active: isActive }),
}

export const adminAnalyticsApi = {
  getAnalytics: () => axios.get<AdminAnalytics>('/admin/analytics'),
}
