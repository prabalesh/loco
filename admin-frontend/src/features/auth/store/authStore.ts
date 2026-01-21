import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../../../types'
import { adminAuthApi } from '../../../lib/api/admin'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  isInitialized: boolean
  setUser: (user: User | null) => void
  logout: () => void
  checkAuth: () => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      isInitialized: false,
      setUser: (user) => set({ user, isAuthenticated: !!user, isInitialized: true }),
      logout: () => set({ user: null, isAuthenticated: false, isInitialized: true }),
      checkAuth: async () => {
        set({ isLoading: true })
        try {
          const response = await adminAuthApi.getProfile()
          set({
            user: response.data.data,
            isAuthenticated: true,
            isInitialized: true,
            isLoading: false
          })
        } catch (error) {
          set({
            user: null,
            isAuthenticated: false,
            isInitialized: true,
            isLoading: false
          })
        }
      },
    }),
    {
      name: 'admin-auth-storage',
      partialize: (state) => ({
        isAuthenticated: state.isAuthenticated,
        user: state.user
      }),
    }
  )
)
