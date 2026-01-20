import type { Problem } from '@/features/problems/types'

export interface UserStats {
    total_submissions: number
    accepted_submissions: number
    problems_solved: number
    acceptance_rate: number
    rank: number
}

export interface DifficultyStat {
    difficulty: string
    count: number
}

export interface HeatmapEntry {
    date: string
    count: number
}

export interface Achievement {
    id: number
    slug: string
    name: string
    description: string
    icon_url: string
    xp_reward: number
    category: string
}

export interface UserAchievement {
    id: number
    unlocked_at: string
    achievement: Achievement
}

export interface UserProfile {
    id: number
    username: string
    is_verified: boolean
    created_at: string
    stats: UserStats
    recent_problems: Problem[]
    submission_heatmap: HeatmapEntry[]
    solved_distribution: DifficultyStat[]
    achievements: UserAchievement[]
}
