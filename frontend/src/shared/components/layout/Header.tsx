import { Link } from 'react-router-dom'
import { Code, LogOut, User } from 'lucide-react'
import { useAuth } from '../../../shared/hooks/useAuth'
import { useLogout } from '../../../features/auth/hooks/useLogout'
import { Button } from '../ui/Button'
import { ROUTES } from '../../../shared/constants/routes'
import { CONFIG } from '../../../shared/constants/config'

export const Header = () => {
  const { isAuthenticated, user } = useAuth()
  const { mutate: logout, isPending } = useLogout()

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50 shadow-sm">
      <nav className="max-w-7xl m-auto">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to={ROUTES.HOME} className="flex items-center space-x-2 hover:opacity-80 transition-opacity">
            <Code className="h-8 w-8 text-blue-600" />
            <span className="text-xl font-bold text-gray-900">
              {CONFIG.APP_NAME}
            </span>
          </Link>

          {/* Center Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            <Link
              to={ROUTES.PROBLEMS}
              className="text-gray-700 hover:text-blue-600 font-medium transition-colors px-3 py-2 rounded-md hover:bg-gray-50"
            >
              Problems
            </Link>
            <Link
              to={ROUTES.LEADERBOARD}
              className="text-gray-700 hover:text-blue-600 font-medium transition-colors px-3 py-2 rounded-md hover:bg-gray-50"
            >
              Leaderboard
            </Link>
            {isAuthenticated && (
              <Link
                to={ROUTES.SUBMISSIONS}
                className="text-gray-700 hover:text-blue-600 font-medium transition-colors px-3 py-2 rounded-md hover:bg-gray-50"
              >
                Submissions
              </Link>
            )}
          </div>

          {/* Right Side - Auth Buttons */}
          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <Link
                  to={ROUTES.PROFILE}
                  className="hidden md:flex items-center space-x-2 text-gray-700 hover:text-blue-600 transition-colors px-3 py-2 rounded-md hover:bg-gray-50"
                >
                  <User className="h-5 w-5" />
                  <span className="font-medium">{user?.username}</span>
                </Link>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => logout()}
                  isLoading={isPending}
                  className="flex items-center"
                >
                  <LogOut className="h-4 w-4 mr-2" />
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Link to={ROUTES.LOGIN}>
                  <Button variant="ghost" size="sm">
                    Login
                  </Button>
                </Link>
                <Link to={ROUTES.REGISTER}>
                  <Button variant="primary" size="sm">
                    Sign Up
                  </Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </nav>
    </header>
  )
}
