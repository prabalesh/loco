import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Mail, ArrowLeft, AlertCircle, CheckCircle2 } from 'lucide-react'
import { Card } from '@/shared/components/ui/Card'
import { authApi } from '../features/auth/api/authApi'
import { ROUTES } from '@/shared/constants/routes'
import { toast } from 'react-hot-toast'
import { Button } from '@/shared/components/ui/Button'

export const VerifyEmailPage = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const token = searchParams.get('token') || ''

  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    if (!token) {
      // No token, redirect to register page or show message
      setError('Verification token missing.')
      return
    }

    const verify = async () => {
      setLoading(true)
      setError('')
      try {
        await authApi.verifyEmail(token)
        setSuccess(true)
        toast.success('Email verified successfully! 🎉')

        // Optional: Redirect to login after short delay
        setTimeout(() => {
          navigate(ROUTES.LOGIN + '?verified=true')
        }, 1500)
      } catch (err: any) {
        const errorMsg = err.response?.data?.error || 'Verification failed'
        setError(errorMsg)
        toast.error(errorMsg)
      } finally {
        setLoading(false)
      }
    }

    verify()
  }, [token, navigate])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4 py-12">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md"
      >
        <Link
          to={ROUTES.REGISTER}
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6 transition-colors"
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Register
        </Link>

        <Card className="p-8">
          <div className="text-center mb-6">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Mail className="h-8 w-8 text-blue-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Verify Your Email</h1>
            {loading && <p className="text-gray-600">Verifying your email...</p>}

            {error && (
              <>
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700 mb-4">
                  <AlertCircle className="inline-block h-5 w-5 mr-2 align-text-bottom" />
                  {error}
                </div>
                <Button onClick={() => navigate(ROUTES.LOGIN)} className="w-full">
                  Go to Login
                </Button>
              </>
            )}

            {success && (
              <div className="bg-green-50 border border-green-200 rounded-lg p-4 text-green-700">
                <CheckCircle2 className="inline-block h-5 w-5 mr-2 align-text-bottom" />
                Email verified successfully! Redirecting to login...
              </div>
            )}
          </div>
        </Card>
      </motion.div>
    </div>
  )
}
