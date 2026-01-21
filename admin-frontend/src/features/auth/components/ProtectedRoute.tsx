import { Navigate } from 'react-router-dom'
import { CircularProgress, Box } from '@mui/material'
import { useAuthStore } from '../store/authStore'

interface ProtectedRouteProps {
  children: React.ReactNode
}

export const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated, isInitialized, isLoading } = useAuthStore()

  if (!isInitialized || isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <CircularProgress />
      </Box>
    )
  }

  if (isAuthenticated === false) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}
