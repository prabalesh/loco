import { apiClient } from '@/shared/lib/axios'
import type { UserProfile } from '../types'

export const usersApi = {
    getProfile: (username: string) =>
        apiClient.get<UserProfile>(`/users/${username}`),
}
