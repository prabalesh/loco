import { Navigate } from 'react-router-dom'
import { useAuth } from '../../../shared/hooks/useAuth'
import { ROUTES } from '../../../shared/constants/routes'

interface ProtectedRouteProps {
  children: React.ReactNode
}

export const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  return <>{children}</>
}
