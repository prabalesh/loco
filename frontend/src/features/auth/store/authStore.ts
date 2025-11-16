import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '../../../shared/types/auth.types'

interface AuthStore {
  user: User | null
  isAuthenticated: boolean
  setUser: (user: User) => void
  logout: () => void
  clearUser: () => void
}

export const authStore = create<AuthStore>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,

      setUser: (user: User) => {
        set({ user, isAuthenticated: true })
      },

      logout: () => {
        set({ user: null, isAuthenticated: false })
      },

      clearUser: () => {
        set({ user: null, isAuthenticated: false })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
