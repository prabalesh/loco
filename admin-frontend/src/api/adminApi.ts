import axios from './axios'
import type { User, AdminAnalytics, LoginCredentials, Language } from '../types'
import type { Response, SimpleResponse } from '../types/repsonse'
import type { CreateOrUpdateLanguageRequest } from '../types/request'

export const adminAuthApi = {
  login: (credentials: LoginCredentials) =>
    axios.post('/admin/auth/login', credentials),

  logout: () => axios.post('/admin/auth/logout'),

  refreshToken: () => axios.post('/admin/auth/refresh'),

  getProfile: () => axios.get<User>('/admin/auth/me'),
}

export const adminUsersApi = {
  getAll: () => axios.get<Response<User[]>>('/admin/users'),

  getById: (id: number) => axios.get<User>(`/admin/users/${id}`),

  deleteUser: (id: number) => axios.delete(`/admin/users/${id}`),

  updateRole: (id: number, role: string) =>
    axios.patch(`/admin/users/${id}/role`, { role }),

  updateStatus: (id: number, isActive: boolean) =>
    axios.patch(`/admin/users/${id}/status`, { is_active: isActive }),
}

export const adminLanguagesApi = {
  getAll: () => axios.get<Response<Language[]>>("/admin/languages"),
  create: (values: CreateOrUpdateLanguageRequest) => axios.post<Response<Language>>("/admin/languages", values),
  update: (id: number, values: CreateOrUpdateLanguageRequest) => axios.put<Response<Language>>(`/admin/languages/${id}`, values),
  getById: (id: number) => axios.get<Response<Language>>(`/admin/languages/${id}`),
  activate: (id: number) => axios.post<SimpleResponse>(`/admin/languages/${id}/activate`),
  deactivate: (id: number) => axios.post<SimpleResponse>(`/admin/languages/${id}/deactivate`),
  delete: (id: number) => axios.delete<SimpleResponse>(`/admin/languages/${id}`)
}

export const adminAnalyticsApi = {
  getAnalytics: () => axios.get<AdminAnalytics>('/admin/analytics'),
}
