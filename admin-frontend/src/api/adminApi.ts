import axios from './axios'
import type { User, AdminAnalytics, LoginCredentials, Language, Problem, TestCase } from '../types'
import type { PaginatedResponse, Response, SimpleResponse } from '../types/repsonse'
import type { CreateOrUpdateLanguageRequest, CreateOrUpdateProblemRequest, CreateTestCaseRequest } from '../types/request'

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
  getAllActive: () => axios.get<Response<Language[]>>("/admin/languages/active"),
  create: (values: CreateOrUpdateLanguageRequest) => axios.post<Response<Language>>("/admin/languages", values),
  update: (id: number, values: CreateOrUpdateLanguageRequest) => axios.put<Response<Language>>(`/admin/languages/${id}`, values),
  getById: (id: number) => axios.get<Response<Language>>(`/admin/languages/${id}`),
  activate: (id: number) => axios.post<SimpleResponse>(`/admin/languages/${id}/activate`),
  deactivate: (id: number) => axios.post<SimpleResponse>(`/admin/languages/${id}/deactivate`),
  delete: (id: number) => axios.delete<SimpleResponse>(`/admin/languages/${id}`)
}

export const adminProblemApi = {
  getAll: () => axios.get<PaginatedResponse<Problem[]>>("/admin/problems"),
  getById: (id: number) => axios.get<Response<Problem>>(`/admin/problems/${id}`),
  create: (values: CreateOrUpdateProblemRequest) =>
    axios.post<Problem>("/admin/problems", values),
  update: (id: number, values: CreateOrUpdateProblemRequest) =>
    axios.put<Problem>(`/admin/problems/${id}`, values),
  delete: (id: number) => axios.delete<void>(`/admin/problems/${id}`),
  validateTestCases: (id: number) =>
    axios.post(`/admin/problems/${id}/test-cases/validate`),
}

export const adminTestcaseApi = {
  getAll: (problemId: number) => 
    axios.get<PaginatedResponse<TestCase[]>>(`/admin/problems/${problemId}/test-cases`),
  create: (problemId: number, data: CreateTestCaseRequest) =>
    axios.post<Response<TestCase>>(`/admin/problems/${problemId}/test-cases`, data),
  update: (testcaseId: number, data: CreateTestCaseRequest) =>
    axios.put<Response<TestCase>>(`/admin/test-cases/${testcaseId}`, data),
  delete: (problemId: number, testcaseId: number) =>
    axios.delete(`/admin/problems/${problemId}/test-cases/${testcaseId}`),
};

export const adminAnalyticsApi = {
  getAnalytics: () => axios.get<AdminAnalytics>('/admin/analytics'),
}
