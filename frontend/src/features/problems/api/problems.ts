import { apiClient } from '@/shared/lib/axios'
import type { Problem, ProblemLanguage, ListProblemsRequest, TestCase, Tag, Category, ProblemResponse } from '../types'
import type { ApiResponse } from '@/shared/types/common.types'
import type { PaginatedResponse } from '@/shared/types/common.types'

export const problemsApi = {
    list: (params: ListProblemsRequest) =>
        apiClient.get<PaginatedResponse<Problem>>('/problems', { params }),

    get: (identifier: string) =>
        apiClient.get<ApiResponse<ProblemResponse>>(`/problems/${identifier}`),

    getLanguages: (problemId: number) =>
        apiClient.get<ApiResponse<ProblemLanguage[]>>(`/problems/${problemId}/languages`),

    getSampleTestCases: (problemId: number) =>
        apiClient.get<ApiResponse<TestCase[]>>(`/problems/${problemId}/test-cases/samples`),

    getTags: () =>
        apiClient.get<ApiResponse<Tag[]>>('/tags'),

    getCategories: () =>
        apiClient.get<ApiResponse<Category[]>>('/categories'),

    getStub: (problemId: number, language: string) =>
        apiClient.get<ApiResponse<{ stub_code: string }>>(`/api/v2/problems/${problemId}/stub`, {
            params: { language }
        }),

    submit: (problemId: number, data: { code: string; language_slug: string }) =>
        apiClient.post<ApiResponse<any>>(`/api/v2/problems/${problemId}/submit`, data),
}
