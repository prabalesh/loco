import axiosInstance from '../axios'
import type { User, AdminAnalytics, LoginCredentials, Language, Problem, TestCase, Tag, Category } from '../../types'
import type { PaginatedResponse, Response, SimpleResponse } from '../../types/repsonse'
import type { CreateOrUpdateLanguageRequest, CreateOrUpdateProblemRequest, CreateTestCaseRequest } from '../../types/request'
import type { ProblemLanguage, CreateProblemLanguageRequest, UpdateProblemLanguageRequest } from '../../types/problemLanguage'

export const adminAuthApi = {
  login: (credentials: LoginCredentials) =>
    axiosInstance.post<Response<{ token: string; user: User }>>('/admin/auth/login', credentials),

  getProfile: () =>
    axiosInstance.get<Response<User>>('/admin/auth/me'),

  logout: () =>
    axiosInstance.post<SimpleResponse>('/admin/auth/logout'),
}

export const adminAnalyticsApi = {
  getStats: () =>
    axiosInstance.get<Response<AdminAnalytics>>('/admin/analytics'),
}

export const adminUsersApi = {
  getAll: (page = 1, limit = 10, search = '') =>
    axiosInstance.get<PaginatedResponse<User[]>>('/admin/users', {
      params: { page, limit, search }
    }),

  getById: (id: number) => axiosInstance.get<Response<User>>(`/admin/users/${id}`),

  deleteUser: (id: number) => axiosInstance.delete(`/admin/users/${id}`),

  updateRole: (id: number, role: string) =>
    axiosInstance.patch(`/admin/users/${id}/role`, { role }),

  updateStatus: (id: number, isActive: boolean) =>
    axiosInstance.patch(`/admin/users/${id}/status`, { is_active: isActive }),
}

export const adminLanguagesApi = {
  getAll: () => axiosInstance.get<Response<Language[]>>('/admin/languages'),
  getAllActive: () => axiosInstance.get<Response<Language[]>>('/admin/languages/active'),
  create: (data: CreateOrUpdateLanguageRequest) =>
    axiosInstance.post<Response<Language>>('/admin/languages', data),
  update: (id: number, data: CreateOrUpdateLanguageRequest) =>
    axiosInstance.put<Response<Language>>(`/admin/languages/${id}`, data),
  delete: (id: number) =>
    axiosInstance.delete<SimpleResponse>(`/admin/languages/${id}`),
  activate: (id: number) =>
    axiosInstance.post<SimpleResponse>(`/admin/languages/${id}/activate`),
  deactivate: (id: number) =>
    axiosInstance.post<SimpleResponse>(`/admin/languages/${id}/deactivate`),
  getById: (id: number) => axiosInstance.get<Response<Language>>(`/admin/languages/${id}`),
}

export const adminProblemApi = {
  getAll: () => axiosInstance.get<Response<Problem[]>>('/admin/problems'),
  getById: (id: string) => axiosInstance.get<Response<Problem>>(`/admin/problems/${id}`),
  create: (data: CreateOrUpdateProblemRequest) =>
    axiosInstance.post<Response<Problem>>('/admin/problems', data),
  update: (id: string, data: CreateOrUpdateProblemRequest) =>
    axiosInstance.put<Response<Problem>>(`/admin/problems/${id}`, data),
  delete: (id: string) =>
    axiosInstance.delete<SimpleResponse>(`/admin/problems/${id}`),

  // Test case management
  getTestCases: (problemId: string) =>
    axiosInstance.get<Response<TestCase[]>>(`/admin/problems/${problemId}/test-cases`),

  validateTestCases: (problemId: number) =>
    axiosInstance.post<SimpleResponse>(`/admin/problems/${problemId}/test-cases/validate`),

  publish: (problemId: string) =>
    axiosInstance.post<SimpleResponse>(`/admin/problems/${problemId}/publish`),

  v2Create: (data: any) =>
    axiosInstance.post<Response<Problem>>('/api/v2/admin/problems', data),

  v2GetById: (id: string | number) =>
    axiosInstance.get<Response<Problem>>(`/api/v2/admin/problems/${id}`),

  v2Publish: (id: string | number) =>
    axiosInstance.post<Response<any>>(`/api/v2/admin/problems/${id}/publish`),

  v2Validate: (id: string | number, data: { language_slug: string; code: string }) =>
    axiosInstance.post<Response<any>>(`/api/v2/admin/problems/${id}/validate`, data),

  v2GetValidationStatus: (id: string | number) =>
    axiosInstance.get<Response<any>>(`/api/v2/admin/problems/${id}/validation-status`),

  getTags: () => axiosInstance.get<Response<Tag[]>>('/tags'),
  getCategories: () => axiosInstance.get<Response<Category[]>>('/categories'),
}

export const adminTestcaseApi = {
  getAll: (problemId: string) =>
    axiosInstance.get<Response<TestCase[]>>(`/admin/problems/${problemId}/test-cases`),
  create: (problemId: string, data: CreateTestCaseRequest) =>
    axiosInstance.post<Response<TestCase>>(`/admin/problems/${problemId}/test-cases`, data),
  delete: (id: string) =>
    axiosInstance.delete<SimpleResponse>(`/admin/test-cases/${id}`),
  update: (id: string, data: Partial<CreateTestCaseRequest>) =>
    axiosInstance.put<Response<TestCase>>(`/admin/test-cases/${id}`, data),
}

export const adminProblemLanguagesApi = {
  getAll: (problemId: string) =>
    axiosInstance.get<Response<ProblemLanguage[]>>(`/admin/problems/${problemId}/languages`),

  create: (problemId: string, data: CreateProblemLanguageRequest) =>
    axiosInstance.post<Response<ProblemLanguage>>(`/admin/problems/${problemId}/languages`, data),

  update: (problemId: string, languageId: number, data: UpdateProblemLanguageRequest) =>
    axiosInstance.put<Response<ProblemLanguage>>(
      `/admin/problems/${problemId}/languages/${languageId}`,
      data
    ),

  delete: (problemId: string, languageId: number) =>
    axiosInstance.delete<SimpleResponse>(`/admin/problems/${problemId}/languages/${languageId}`),

  validate: (problemId: string, languageId: number) =>
    axiosInstance.post<Response<any>>(`/admin/problems/${problemId}/languages/${languageId}/validate`),

  preview: (problemId: string, languageId: number) =>
    axiosInstance.get<Response<{ combined_code: string }>>(`/admin/problems/${problemId}/languages/${languageId}/preview`),
};

export const adminTagApi = {
  create: (data: { name: string; slug: string }) =>
    axiosInstance.post<Response<Tag>>('/admin/tags', data),
  update: (id: number, data: { name?: string; slug?: string }) =>
    axiosInstance.put<Response<Tag>>(`/admin/tags/${id}`, data),
  delete: (id: number) =>
    axiosInstance.delete<SimpleResponse>(`/admin/tags/${id}`),
}

export const adminCategoryApi = {
  create: (data: { name: string; slug: string }) =>
    axiosInstance.post<Response<Category>>('/admin/categories', data),
  update: (id: number, data: { name?: string; slug?: string }) =>
    axiosInstance.put<Response<Category>>(`/admin/categories/${id}`, data),
  delete: (id: number) =>
    axiosInstance.delete<SimpleResponse>(`/admin/categories/${id}`),
}

export const adminSubmissionsApi = {
  getById: (problemId: number, submissionId: number) => axiosInstance.get<Response<any>>(`problems/${problemId}/submissions/${submissionId}`),
}

export const adminCodeGenApi = {
  generateStub: (data: {
    function_name: string;
    return_type: string;
    parameters: Array<{ name: string; type: string; is_custom: boolean }>;
    language_slug: string;
  }) => axiosInstance.post<Response<{ stub_code: string }>>('/api/v2/codegen/stub', data),

  getBoilerplateStats: (problemId: number) =>
    axiosInstance.get<Response<{ total_languages: number; languages: string[] }>>(`/api/v2/problems/${problemId}/boilerplates`),
}
