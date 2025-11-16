import { createBrowserRouter, Navigate } from 'react-router-dom'
import { MainLayout } from '../shared/components/layout/MainLayout'
import { ProtectedRoute } from '../shared/components/common/ProtectedRoute'
import { HomePage } from '../pages/HomePage'
import { NotFoundPage } from '../pages/PageNotFound'
import { LoginPage } from '../features/auth/pages/LoginPage'
import { RegisterPage } from '../features/auth/pages/RegisterPage'
import { ROUTES } from '../shared/constants/routes'
import { ProfilePage } from '@/pages/ProfilePage'
import { UserProfilePage } from '@/pages/UserProfilePage'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: ROUTES.LOGIN,
        element: <LoginPage />,
      },
      {
        path: ROUTES.REGISTER,
        element: <RegisterPage />,
      },
      {
        path: ROUTES.PROBLEMS,
        element: (
          <div className="p-8 text-center">
            <h1 className="text-3xl font-bold">Problems Page (Coming Soon)</h1>
          </div>
        ),
      },
      // Protected routes (add more as you build them)
      {
        path: '/users/:username',
        element: (
          <ProtectedRoute>
            <UserProfilePage />
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.SUBMISSIONS,
        element: (
          <ProtectedRoute>
            <div className="p-8 text-center">
              <h1 className="text-3xl font-bold">Submissions Page (Coming Soon)</h1>
            </div>
          </ProtectedRoute>
        ),
      },
      {
        path: ROUTES.LEADERBOARD,
        element: (
          <div className="p-8 text-center">
            <h1 className="text-3xl font-bold">Leaderboard Page (Coming Soon)</h1>
          </div>
        ),
      },
      {
        path: ROUTES.PROFILE,
        element: (
          <ProtectedRoute>
            <ProfilePage />
          </ProtectedRoute>
        ),
      },
      {
        path: '404',
        element: <NotFoundPage />,
      },
      {
        path: '*',
        element: <Navigate to="/404" replace />,
      },
    ],
  },
])
