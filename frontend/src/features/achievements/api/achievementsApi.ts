import { apiClient as api } from '@/shared/lib/axios'

export interface Achievement {
    id: number
    slug: string
    name: string
    description: string
    icon_url: string
    xp_reward: number
    category: string
    condition_type: string
    condition_value: string
}

export interface UserAchievement {
    id: number
    user_id: number
    achievement_id: number
    unlocked_at: string
    achievement: Achievement
}

export const achievementsApi = {
    list: async (): Promise<Achievement[]> => {
        const response = await api.get<{ data: Achievement[] }>('/achievements')
        return response.data.data
    },

    getMyAchievements: async (): Promise<UserAchievement[]> => {
        const response = await api.get<{ data: UserAchievement[] }>('/users/me/achievements')
        return response.data.data
    },

    getUserAchievements: async (username: string): Promise<UserAchievement[]> => {
        const response = await api.get<{ data: UserAchievement[] }>(`/users/${username}/achievements`)
        return response.data.data
    },
}
