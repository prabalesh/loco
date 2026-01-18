import { apiClient } from '@/shared/lib/axios'
import type { Problem, ProblemLanguage, ListProblemsRequest, TestCase } from '../types'
import type { PaginatedResponse } from '@/shared/types/common.types'

export const problemsApi = {
    list: (params: ListProblemsRequest) =>
        apiClient.get<PaginatedResponse<Problem>>('/problems', { params }),

    get: (identifier: string) =>
        apiClient.get<Problem>(`/problems/${identifier}`),

    getLanguages: (problemId: number) =>
        apiClient.get<ProblemLanguage[]>(`/problems/${problemId}/languages`),

    getSampleTestCases: (problemId: number) =>
        apiClient.get<TestCase[]>(`/problems/${problemId}/test-cases/samples`),
}
