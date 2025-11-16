import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Home, SearchX } from 'lucide-react'
import { Button } from '../shared/components/ui/Button'
import { ROUTES } from '../shared/constants/routes'

export const NotFoundPage = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center px-4">
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.5 }}
        className="text-center max-w-2xl"
      >
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="mb-8"
        >
          <SearchX className="h-24 w-24 text-gray-400 mx-auto mb-6" />
          <h1 className="text-9xl font-bold text-gray-800 mb-4">404</h1>
          <h2 className="text-3xl font-bold text-gray-900 mb-2">
            Page Not Found
          </h2>
          <p className="text-lg text-gray-600 mb-8">
            Oops! The page you're looking for doesn't exist. It might have been moved or deleted.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="flex flex-col sm:flex-row gap-4 justify-center"
        >
          <Link to={ROUTES.HOME}>
            <Button variant="primary" size="lg">
              <Home className="mr-2 h-5 w-5" />
              Go Home
            </Button>
          </Link>
          <Link to={ROUTES.PROBLEMS}>
            <Button variant="outline" size="lg">
              Browse Problems
            </Button>
          </Link>
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="mt-12"
        >
          <p className="text-sm text-gray-500">
            If you believe this is an error, please contact support
          </p>
        </motion.div>
      </motion.div>
    </div>
  )
}
