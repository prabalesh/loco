import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ConfigProvider, Spin } from 'antd'
import { Toaster } from 'react-hot-toast'
import { AdminLogin } from './features/auth/components/AdminLogin'
import { Dashboard } from './features/dashboard/components/Dashboard'
import { UsersList } from './features/users/components/UsersList'
import { AdminLayout } from './components/layout/AdminLayout'
import { ProtectedRoute } from './components/common/ProtectedRoute'
import { useAuthStore } from './features/auth/store/authStore'
import { useEffect, useState } from 'react'
import { adminAuthApi } from './api/adminApi'
import ProblemList from './features/problems/components/ProblemList'
import LanguageList from './features/languages/components/LangageList'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

function App() {
  const [isLoading, setIsLoading] = useState(true)
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const logout = useAuthStore((state) => state.logout)

  const setUser = useAuthStore((s) => s.setUser)

  useEffect(() => {
    const timer = setTimeout(() => setIsLoading(false), 100)
    return () => clearTimeout(timer)
  }, [])

  useEffect(() => {
    const initAuth = async () => {
      try {
        const response = await adminAuthApi.getProfile()
        setUser(response.data)
      } catch (err: any) {
        if (err?.response?.status === 401) {
          logout()
        }
        console.error('Auth check failed:', err)
      } finally {
        setIsLoading(false)
      }
    }

    initAuth()
  }, [setUser, logout])

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <Spin size="large" />
      </div>
    )
  }

  return (
    <ConfigProvider theme={{ token: { colorPrimary: '#1890ff' } }}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route
              path="/login"
              element={isAuthenticated? <Navigate to="/" replace /> : <AdminLogin />}
            />

            <Route
              element={
                <ProtectedRoute>
                  <AdminLayout />
                </ProtectedRoute>
              }
            >
              <Route path="/" element={<Dashboard />} />
              <Route path="/users" element={<UsersList />} />
              <Route path="/problems" element={<ProblemList />} />
              <Route path="/languages" element={<LanguageList />} />
            </Route>

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
        <Toaster position="top-right" />
      </QueryClientProvider>
    </ConfigProvider>
  )
}

export default App;