import { apiClient } from '@/shared/lib/axios'
import type { Submission, TestCaseResult } from '../types'

export interface RunCodeResult {
    status: string
    error_message?: string
    passed_test_cases: number
    total_test_cases: number
    results: TestCaseResult[]
}

export const submissionsApi = {
    runCode: (problemId: number, languageId: number, code: string) =>
        apiClient.post<RunCodeResult>(`/problems/${problemId}/run`, {
            language_id: languageId,
            code,
        }),

    submit: (problemId: number, languageId: number, code: string) =>
        apiClient.post<Submission>(`/problems/${problemId}/submissions`, {
            language_id: languageId,
            code,
        }),

    get: (id: number) =>
        apiClient.get<Submission>(`/problems/${id}/submissions/${id}`),

    list: (problemId: number, page = 1, limit = 10) =>
        apiClient.get<{ data: Submission[]; total: number }>(`/problems/${problemId}/submissions`, {
            params: { page, limit },
        }),

    listUserSubmissions: (page = 1, limit = 10) =>
        apiClient.get<{ data: { data: Submission[]; }, total: number, limit: number, page: number }>(`/submissions`, {
            params: { page, limit },
        }),
}
