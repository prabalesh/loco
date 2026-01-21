import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ThemeProvider, createTheme, CssBaseline, CircularProgress, Box } from '@mui/material'
import { Toaster } from 'react-hot-toast'
import { lazy, Suspense, useEffect, useState } from 'react'

const AdminLogin = lazy(() => import('./features/auth/components/AdminLogin').then(m => ({ default: m.AdminLogin })))
const Dashboard = lazy(() => import('./features/dashboard/components/Dashboard').then(m => ({ default: m.Dashboard })))
const UsersList = lazy(() => import('./features/users/components/UsersList').then(m => ({ default: m.UsersList })))
const ProblemList = lazy(() => import('./features/problems/components/ProblemList'))
const LanguageList = lazy(() => import('./features/languages/components/LanguageList'))
const TagList = lazy(() => import('./features/tags/components/TagList'))
const CategoryList = lazy(() => import('./features/categories/components/CategoryList'))
const CreateProblem = lazy(() => import('./features/problems/pages/CreateProblem'))
const ProblemTestCases = lazy(() => import('./features/problems/pages/ProblemTestCases'))
const ProblemLanguage = lazy(() => import('./features/problems/pages/ProblemLanguage'))
const ProblemValidate = lazy(() => import('./features/problems/pages/ProblemValidate'))
const ProblemCreationFormV2 = lazy(() => import('./components/v2/ProblemCreationForm').then(m => ({ default: m.ProblemCreationForm })))

import { AdminLayout } from './components/layout/AdminLayout'
import { ProtectedRoute } from './features/auth/components/ProtectedRoute'
import { useAuthStore } from './features/auth/store/authStore'
import { adminAuthApi } from './lib/api/admin'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

const theme = createTheme({
  palette: {
    primary: {
      main: '#6366f1',
    },
  },
  typography: {
    fontFamily: "'Inter', sans-serif",
  },
  shape: {
    borderRadius: 8,
  },
})

const PageLoader = () => (
  <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '80vh' }}>
    <CircularProgress size={40} thickness={4} />
  </Box>
)

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
        setUser(response.data.data)
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
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <CircularProgress size={48} />
      </Box>
    )
  }

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Suspense fallback={<PageLoader />}>
            <Routes>
              <Route
                path="/login"
                element={isAuthenticated ? <Navigate to="/" replace /> : <AdminLogin />}
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
                <Route path="/problems/create" element={<CreateProblem />} />
                <Route path="/problems/create/v2" element={<ProblemCreationFormV2 />} />
                <Route path="/problems/edit/:id" element={<CreateProblem />} />
                <Route path="/languages" element={<LanguageList />} />
                <Route path="/tags" element={<TagList />} />
                <Route path="/categories" element={<CategoryList />} />
                <Route path="/problems/:problemId/testcases" element={<ProblemTestCases />} />
                <Route path="/problems/:problemId/languages" element={<ProblemLanguage />} />
                <Route path="/problems/:problemId/validate" element={<ProblemValidate />} />
              </Route>
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </Suspense>
        </BrowserRouter>
        <Toaster position="top-right" />
      </QueryClientProvider>
    </ThemeProvider>
  )
}

export default App;