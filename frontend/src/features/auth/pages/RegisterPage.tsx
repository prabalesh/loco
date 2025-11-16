import { Link, Navigate } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Code } from 'lucide-react'
import { RegisterForm } from '../components/RegisterForm'
import { Card } from '../../../shared/components/ui/Card'
import { useAuth } from '../../../shared/hooks/useAuth'
import { ROUTES } from '../../../shared/constants/routes'
import { CONFIG } from '../../../shared/constants/config'

export const RegisterPage = () => {
  const { isAuthenticated } = useAuth()

  // Redirect if already logged in
  if (isAuthenticated) {
    return <Navigate to={ROUTES.HOME} replace />
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4 py-12">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md"
      >
        <div className="text-center mb-8">
          <Link to={ROUTES.HOME} className="inline-flex items-center space-x-2 mb-6">
            <Code className="h-10 w-10 text-blue-600" />
            <span className="text-2xl font-bold text-gray-900">
              {CONFIG.APP_NAME}
            </span>
          </Link>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Create Your Account
          </h1>
          <p className="text-gray-600">
            Start your journey to coding mastery
          </p>
        </div>

        <Card className="p-8">
          <RegisterForm />

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Already have an account?{' '}
              <Link
                to={ROUTES.LOGIN}
                className="text-blue-600 hover:text-blue-700 font-medium"
              >
                Sign in
              </Link>
            </p>
          </div>
        </Card>

        <p className="text-center text-sm text-gray-600 mt-8">
          By continuing, you agree to our Terms of Service and Privacy Policy
        </p>
      </motion.div>
    </div>
  )
}
