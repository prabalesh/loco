import { createBrowserRouter, Navigate } from 'react-router-dom'
import { MainLayout } from '../shared/components/layout/MainLayout'
import { ProtectedRoute } from '../shared/components/common/ProtectedRoute'
import { HomePage } from '../pages/HomePage'
import { NotFoundPage } from '../pages/PageNotFound'
import { LoginPage } from '../features/auth/pages/LoginPage'
import { RegisterPage } from '../features/auth/pages/RegisterPage'
import { ROUTES } from '../shared/constants/routes'
import { ProfilePage } from '@/pages/ProfilePage'
import { VerifyEmailPage } from '@/features/auth/pages/VerifyEmailPage'
import { ForgotPasswordPage } from '@/features/auth/pages/ForgetPasswordPage'
import { ResetPasswordPage } from '@/features/auth/pages/ResetPasswordPage'
import { ProblemsPage } from '@/features/problems/pages/ProblemsPage'
import { ProblemDetailPage } from '@/features/problems/pages/ProblemDetailPage'
import { ProfilePage as UserProfileView } from '@/features/users/pages/ProfilePage'
import { SubmissionsPage } from '@/pages/SubmissionsPage'

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
        path: ROUTES.VERIFY_EMAIL,
        element: <VerifyEmailPage />,
      },
      {
        path: ROUTES.PROBLEMS,
        element: <ProblemsPage />,
      },
      {
        path: '/problems/:slug',
        element: <ProblemDetailPage />,
      },
      {
        path: ROUTES.FORGOT_PASSWORD,
        element: <ForgotPasswordPage />,
      },
      {
        path: ROUTES.RESET_PASSWORD,
        element: <ResetPasswordPage />,
      },
      {
        path: '/users/:username',
        element: <UserProfileView />,
      },
      {
        path: ROUTES.SUBMISSIONS,
        element: (
          <ProtectedRoute>
            <SubmissionsPage />
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
