import { Link } from 'react-router-dom'
import { Code, LogOut, User, Menu, X } from 'lucide-react'
import { useState } from 'react'
import { useAuth } from '../../../shared/hooks/useAuth'
import { useLogout } from '../../../features/auth/hooks/useLogout'
import { Button } from '../ui/Button'
import { ROUTES } from '../../../shared/constants/routes'
import { CONFIG } from '../../../shared/constants/config'


export const Header = () => {
  const { isAuthenticated, user } = useAuth()
  const { mutate: logout, isPending } = useLogout()
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)

  const closeMobileMenu = () => setMobileMenuOpen(false)

  return (
    <header className="bg-white border-b border-gray-200 sticky top-0 z-50 shadow-sm">
      <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link
            to={ROUTES.HOME}
            className="flex items-center space-x-2 hover:opacity-80 transition-opacity z-50"
            onClick={closeMobileMenu}
          >
            <Code className="h-6 w-6 sm:h-8 sm:w-8 text-blue-600" />
            <span className="text-lg sm:text-xl font-bold text-gray-900">
              {CONFIG.APP_NAME}
            </span>
          </Link>

          {/* Desktop Navigation */}
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

          {/* Desktop Auth Buttons */}
          <div className="hidden md:flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <Link
                  to={ROUTES.PROFILE}
                  className="flex items-center space-x-2 text-gray-700 hover:text-blue-600 transition-colors px-3 py-2 rounded-md hover:bg-gray-50"
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

          {/* Mobile Menu Button */}
          <button
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            className="md:hidden p-2 rounded-md text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors z-50"
            aria-label="Toggle mobile menu"
            aria-expanded={mobileMenuOpen}
          >
            {mobileMenuOpen ? (
              <X className="h-6 w-6" />
            ) : (
              <Menu className="h-6 w-6" />
            )}
          </button>
        </div>

        {/* Mobile Menu */}
        <div
          className={`md:hidden overflow-hidden transition-all duration-300 ease-in-out ${mobileMenuOpen ? 'max-h-screen opacity-100' : 'max-h-0 opacity-0'
            }`}
        >
          <div className="py-4 space-y-2 border-t border-gray-200">
            {/* Mobile Navigation Links */}
            <Link
              to={ROUTES.PROBLEMS}
              className="block text-gray-700 hover:text-blue-600 hover:bg-gray-50 font-medium px-4 py-3 rounded-md transition-colors"
              onClick={closeMobileMenu}
            >
              Problems
            </Link>
            <Link
              to={ROUTES.LEADERBOARD}
              className="block text-gray-700 hover:text-blue-600 hover:bg-gray-50 font-medium px-4 py-3 rounded-md transition-colors"
              onClick={closeMobileMenu}
            >
              Leaderboard
            </Link>
            {isAuthenticated && (
              <Link
                to={ROUTES.SUBMISSIONS}
                className="block text-gray-700 hover:text-blue-600 hover:bg-gray-50 font-medium px-4 py-3 rounded-md transition-colors"
                onClick={closeMobileMenu}
              >
                Submissions
              </Link>
            )}

            {/* Mobile Auth Section */}
            <div className="pt-4 mt-4 border-t border-gray-200 space-y-2">
              {isAuthenticated ? (
                <>
                  <Link
                    to={ROUTES.PROFILE}
                    className="flex items-center space-x-2 text-gray-700 hover:text-blue-600 hover:bg-gray-50 px-4 py-3 rounded-md transition-colors"
                    onClick={closeMobileMenu}
                  >
                    <User className="h-5 w-5" />
                    <span className="font-medium">{user?.username}</span>
                  </Link>
                  <button
                    onClick={() => {
                      logout()
                      closeMobileMenu()
                    }}
                    disabled={isPending}
                    className="w-full flex items-center justify-center space-x-2 text-gray-700 hover:text-red-600 hover:bg-red-50 px-4 py-3 rounded-md transition-colors disabled:opacity-50"
                  >
                    <LogOut className="h-5 w-5" />
                    <span className="font-medium">
                      {isPending ? 'Logging out...' : 'Logout'}
                    </span>
                  </button>
                </>
              ) : (
                <div className="space-y-2 px-4">
                  <Link
                    to={ROUTES.LOGIN}
                    className="block"
                    onClick={closeMobileMenu}
                  >
                    <Button variant="ghost" size="md" className="w-full">
                      Login
                    </Button>
                  </Link>
                  <Link
                    to={ROUTES.REGISTER}
                    className="block"
                    onClick={closeMobileMenu}
                  >
                    <Button variant="primary" size="md" className="w-full">
                      Sign Up
                    </Button>
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>
      </nav>
    </header>
  )
}
