import { lazy, Suspense } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { MainLayout } from '../shared/components/layout/MainLayout'
import { ProtectedRoute } from '../shared/components/common/ProtectedRoute'
import { ROUTES } from '../shared/constants/routes'
import { Skeleton } from '@/shared/components/ui/Skeleton'

// Lazy loaded components
const HomePage = lazy(() => import('../pages/HomePage').then(m => ({ default: m.HomePage })))
const NotFoundPage = lazy(() => import('../pages/PageNotFound').then(m => ({ default: m.NotFoundPage })))
const LoginPage = lazy(() => import('../features/auth/pages/LoginPage').then(m => ({ default: m.LoginPage })))
const RegisterPage = lazy(() => import('../features/auth/pages/RegisterPage').then(m => ({ default: m.RegisterPage })))
const ProfilePage = lazy(() => import('@/pages/ProfilePage').then(m => ({ default: m.ProfilePage })))
const VerifyEmailPage = lazy(() => import('@/features/auth/pages/VerifyEmailPage').then(m => ({ default: m.VerifyEmailPage })))
const ForgotPasswordPage = lazy(() => import('@/features/auth/pages/ForgetPasswordPage').then(m => ({ default: m.ForgotPasswordPage })))
const ResetPasswordPage = lazy(() => import('@/features/auth/pages/ResetPasswordPage').then(m => ({ default: m.ResetPasswordPage })))
const ProblemsPage = lazy(() => import('@/features/problems/pages/ProblemsPage').then(m => ({ default: m.ProblemsPage })))
const ProblemDetailPage = lazy(() => import('@/features/problems/pages/ProblemDetailPage').then(m => ({ default: m.ProblemDetailPage })))
const UserProfileView = lazy(() => import('@/features/users/pages/ProfilePage').then(m => ({ default: m.ProfilePage })))
const SubmissionsPage = lazy(() => import('@/pages/SubmissionsPage').then(m => ({ default: m.SubmissionsPage })))
const LeaderboardPage = lazy(() => import('@/features/users/pages/LeaderboardPage').then(m => ({ default: m.LeaderboardPage })))
const AchievementsPage = lazy(() => import('@/features/achievements/pages/AchievementsPage').then(m => ({ default: m.AchievementsPage })))

const PageLoader = () => (
  <div className="p-8 space-y-4">
    <Skeleton className="h-12 w-3/4" />
    <Skeleton className="h-64 w-full" />
    <Skeleton className="h-32 w-full" />
  </div>
)

export const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <Suspense fallback={<PageLoader />}>
        <MainLayout />
      </Suspense>
    ),
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
        element: <LeaderboardPage />,
      },
      {
        path: ROUTES.ACHIEVEMENTS,
        element: <AchievementsPage />,
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
