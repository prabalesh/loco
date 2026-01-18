import { apiClient } from '@/shared/lib/axios'

export interface LeaderboardEntry {
    rank: number
    user_id: number
    username: string
    problems_solved: number
    total_submissions: number
    acceptance_rate: number
}

export const leaderboardApi = {
    getLeaderboard: (limit: number = 100) =>
        apiClient.get<LeaderboardEntry[]>(`/leaderboard?limit=${limit}`),
}
