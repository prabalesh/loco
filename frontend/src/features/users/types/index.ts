import type { Problem } from '@/features/problems/types'

export interface UserStats {
    total_submissions: number
    accepted_submissions: number
    problems_solved: number
    acceptance_rate: number
}

export interface UserProfile {
    id: number
    username: string
    is_verified: boolean
    created_at: string
    stats: UserStats
    recent_problems: Problem[]
}
