import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { useAuthStore } from './features/auth/store/authStore'

import { AdminLayout } from './components/layout/AdminLayout'
import { ProtectedRoute } from './features/auth/components/ProtectedRoute'
import { AdminLogin } from './features/auth/components/AdminLogin'

// Core Components
import { Dashboard } from './features/dashboard/components/Dashboard'
import { UsersList as Users } from './features/users/components/UsersList'
import ProblemCreateV2 from './pages/ProblemCreateV2'
import Languages from './features/languages/components/LanguageList'
import Tags from './features/tags/components/TagList'
import Categories from './features/categories/components/CategoryList'
import ProblemValidate from './features/problems/pages/ProblemValidate'
import BulkImport from './pages/BulkImport'
import ProblemsList from './pages/ProblemsList'
import ProblemManagement from './pages/ProblemManagement'
import PistonExecutions from './pages/PistonExecutions'
import SubmissionsList from './pages/SubmissionsList'

export const App = () => {
  const checkAuth = useAuthStore((state) => state.checkAuth)

  useEffect(() => {
    checkAuth()
  }, [checkAuth])

  return (
    <>
      <Toaster position="top-right" />
      <Routes>
        <Route path="/login" element={<AdminLogin />} />

        <Route element={<ProtectedRoute><AdminLayout /></ProtectedRoute>}>
          <Route path="/" element={<Dashboard />} />
          <Route path="/users" element={<Users />} />
          <Route path="/problems/create" element={<ProblemCreateV2 />} />
          <Route path="/languages" element={<Languages />} />
          <Route path="/tags" element={<Tags />} />
          <Route path="/categories" element={<Categories />} />
          <Route path="/problems/:problemId/validate" element={<ProblemValidate />} />
          <Route path="/problems/bulk-import" element={<BulkImport />} />
          <Route path="/problems" element={<ProblemsList />} />
          <Route path="/problems/:id/manage" element={<ProblemManagement />} />
          <Route path="/submissions" element={<SubmissionsList />} />
          <Route path="/piston/executions" element={<PistonExecutions />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  )
}

export default App