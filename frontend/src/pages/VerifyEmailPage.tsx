import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Mail, ArrowLeft, AlertCircle, Loader2, CheckCircle2 } from 'lucide-react'
import { Button } from '@/shared/components/ui/Button'
import { Input } from '@/shared/components/ui/Input'
import { Card } from '@/shared/components/ui/Card'
import { authApi } from '../features/auth/api/authApi'
import { ROUTES } from '@/shared/constants/routes'
import { toast } from 'react-hot-toast'

export const VerifyEmailPage = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const email = searchParams.get('email') || ''

  const [otp, setOtp] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [cooldown, setCooldown] = useState(0)
  const [resending, setResending] = useState(false)
  const [attemptsLeft, setAttemptsLeft] = useState(5)

  useEffect(() => {
    if (!email) {
      navigate(ROUTES.REGISTER)
    }
  }, [email, navigate])

  useEffect(() => {
    if (cooldown > 0) {
      const timer = setTimeout(() => setCooldown(cooldown - 1), 1000)
      return () => clearTimeout(timer)
    }
  }, [cooldown])

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await authApi.verifyEmail(email, otp)
      toast.success('Email verified successfully! 🎉')
      setTimeout(() => {
        navigate(ROUTES.LOGIN + '?verified=true')
      }, 1500)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Verification failed'
      setError(errorMsg)
      
      if (errorMsg.includes('maximum')) {
        setAttemptsLeft(0)
        toast.error('Maximum attempts exceeded. Please request a new code.')
      } else if (errorMsg.includes('invalid') || errorMsg.includes('expired')) {
        setAttemptsLeft(prev => Math.max(0, prev - 1))
        toast.error('Invalid or expired code')
      } else {
        toast.error(errorMsg)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleResend = async () => {
    setError('')
    setResending(true)

    try {
      await authApi.resendVerificationEmail(email)
      setCooldown(120) // 2 minutes cooldown
      setAttemptsLeft(5) // Reset attempts
      setOtp('') // Clear OTP input
      toast.success('New verification code sent! Check your email. 📧')
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to resend'
      
      // Extract cooldown from error message if present
      const match = errorMsg.match(/(\d+) seconds/)
      if (match) {
        const remainingSeconds = parseInt(match[1])
        setCooldown(remainingSeconds)
        toast.error(`Please wait ${remainingSeconds} seconds before requesting again`)
      } else {
        toast.error(errorMsg)
      }
    } finally {
      setResending(false)
    }
  }

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
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              Verify Your Email
            </h1>
            <p className="text-gray-600 text-sm">
              We've sent a 6-digit verification code to
            </p>
            <p className="font-semibold text-gray-900 mt-1 break-all">{email}</p>
          </div>

          {/* Attempts Warning */}
          {attemptsLeft <= 2 && attemptsLeft > 0 && (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mb-4"
            >
              <div className="flex items-start">
                <AlertCircle className="h-5 w-5 text-yellow-600 mt-0.5 mr-2 flex-shrink-0" />
                <div className="text-sm text-yellow-800">
                  <strong>Warning:</strong> You have {attemptsLeft} attempt{attemptsLeft !== 1 ? 's' : ''} remaining.
                </div>
              </div>
            </motion.div>
          )}

          {/* Max Attempts Exceeded */}
          {attemptsLeft === 0 && (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="bg-red-50 border border-red-200 rounded-lg p-3 mb-4"
            >
              <div className="flex items-start">
                <AlertCircle className="h-5 w-5 text-red-600 mt-0.5 mr-2 flex-shrink-0" />
                <div className="text-sm text-red-800">
                  <strong>Maximum attempts exceeded.</strong> Please request a new verification code below.
                </div>
              </div>
            </motion.div>
          )}

          <form onSubmit={handleVerify} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Verification Code
              </label>
              <Input
                type="text"
                value={otp}
                onChange={(e) => setOtp(e.target.value.replace(/\D/g, '').slice(0, 6))}
                placeholder="000000"
                maxLength={6}
                error={error}
                disabled={attemptsLeft === 0}
                className="text-center text-2xl tracking-[0.5em] font-mono"
                autoComplete="off"
              />
            </div>

            <Button
              type="submit"
              variant="primary"
              className="w-full"
              isLoading={loading}
              disabled={otp.length !== 6 || attemptsLeft === 0}
            >
              {loading ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Verifying...
                </>
              ) : (
                <>
                  <CheckCircle2 className="h-4 w-4 mr-2" />
                  Verify Email
                </>
              )}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600 mb-3">
              Didn't receive the code?
            </p>
            <Button
              type="button"
              variant="ghost"
              onClick={handleResend}
              disabled={cooldown > 0 || resending}
              className="w-full"
            >
              {resending ? (
                <>
                  <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                  Sending...
                </>
              ) : cooldown > 0 ? (
                `Resend in ${Math.floor(cooldown / 60)}:${(cooldown % 60).toString().padStart(2, '0')}`
              ) : (
                <>
                  <Mail className="h-4 w-4 mr-2" />
                  Resend Code
                </>
              )}
            </Button>
          </div>

          <div className="mt-6 p-4 bg-blue-50 rounded-lg border border-blue-100">
            <h3 className="text-sm font-semibold text-blue-900 mb-2 flex items-center">
              <Mail className="h-4 w-4 mr-2" />
              Check your email
            </h3>
            <ul className="text-sm text-blue-800 space-y-1.5">
              <li className="flex items-start">
                <span className="mr-2">•</span>
                <span>Code expires in 10 minutes</span>
              </li>
              <li className="flex items-start">
                <span className="mr-2">•</span>
                <span>Check your spam/junk folder</span>
              </li>
              <li className="flex items-start">
                <span className="mr-2">•</span>
                <span>You have 5 attempts per code</span>
              </li>
              <li className="flex items-start">
                <span className="mr-2">•</span>
                <span>New code can be requested after 2 minutes</span>
              </li>
            </ul>
          </div>
        </Card>

        <p className="text-center text-sm text-gray-600 mt-6">
          Already verified?{' '}
          <Link to={ROUTES.LOGIN} className="text-blue-600 hover:text-blue-700 font-medium">
            Sign in
          </Link>
        </p>
      </motion.div>
    </div>
  )
}
