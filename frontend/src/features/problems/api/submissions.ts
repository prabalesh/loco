import { apiClient } from '@/shared/lib/axios'
import type { Submission, TestCaseResult } from '../types'
import type { ApiResponse } from '@/shared/types/common.types'

export interface RunCodeResult {
    status: string
    error_message?: string
    passed_test_cases: number
    total_test_cases: number
    results: TestCaseResult[]
}

export const submissionsApi = {
    runCode: (problemId: number, languageId: number, code: string) =>
        apiClient.post<ApiResponse<RunCodeResult>>(`/problems/${problemId}/run`, {
            language_id: languageId,
            code,
        }),

    submit: (problemId: number, languageId: number, code: string) =>
        apiClient.post<ApiResponse<Submission>>(`/problems/${problemId}/submissions`, {
            language_id: languageId,
            code,
        }),

    get: (problemId: number, id: number) =>
        apiClient.get<ApiResponse<Submission>>(`/problems/${problemId}/submissions/${id}`),

    list: (problemId: number, page = 1, limit = 10) =>
        apiClient.get<{ data: Submission[]; total: number }>(`/problems/${problemId}/submissions`, {
            params: { page, limit },
        }),

    listUserSubmissions: (page = 1, limit = 10) =>
        apiClient.get<{ data: Submission[]; total: number }>(`/submissions`, {
            params: { page, limit },
        }),
}
