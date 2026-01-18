import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useState, useEffect } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { motion, AnimatePresence } from 'framer-motion'
import { Code, CheckCircle, Mail, AlertCircle, Loader2 } from 'lucide-react'
import { Input } from '@/shared/components/ui/Input'
import { Button } from '@/shared/components/ui/Button'
import { Card } from '@/shared/components/ui/Card'
import { useLogin } from '../hooks/useLogin'
import { authApi } from '../api/authApi'
import { ROUTES } from '@/shared/constants/routes'
import { CONFIG } from '@/shared/constants/config'
import { useAuth } from '@/shared/hooks/useAuth'
import { Navigate } from 'react-router-dom'
import { toast } from 'react-hot-toast'
import type { LoginResponse } from '@/shared/types/auth.types'
import type { AxiosError } from 'axios'

const loginSchema = z.object({
  email: z.string().min(1, 'Email is required').email('Invalid email format'),
  password: z.string().min(1, 'Password is required'),
})

type LoginFormData = z.infer<typeof loginSchema>

export const LoginPage = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { isAuthenticated } = useAuth()
  const { mutate: login, isPending } = useLogin()

  const [showResendOption, setShowResendOption] = useState(false)
  const [resendEmail, setResendEmail] = useState('')
  const [resendLoading, setResendLoading] = useState(false)
  const [resendCooldown, setResendCooldown] = useState(0)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  })

  useEffect(() => {
    if (searchParams.get('verified') === 'true') {
      toast.success('Email verified! You can now log in.', {
        icon: <CheckCircle className="h-5 w-5 text-green-600" />,
      })
    }
  }, [searchParams])

  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => setResendCooldown(resendCooldown - 1), 1000)
      return () => clearTimeout(timer)
    }
  }, [resendCooldown])

  const onSubmit = (data: LoginFormData) => {
    setShowResendOption(false)
    login(data, {
      onSuccess: () => {
        toast.success('Login successful!')
        navigate(ROUTES.HOME)
      },
      onError: (err: any) => {
        const error = err as AxiosError<LoginResponse>
        console.log(error.response?.data)
        const errorMessage = error.response?.data.data.message || 'Login failed'

        if (errorMessage.toLowerCase().includes('verify your email')) {
          setShowResendOption(true)
          setResendEmail(data.email)
          toast.error('Please verify your email first.', { duration: 5000 })
        } else {
          console.log("triggered")
          toast.error(errorMessage)
        }
      },
    })
  }

  const handleResendVerification = async () => {
    setResendLoading(true)
    try {
      await authApi.resendVerificationEmail(resendEmail)
      setResendCooldown(120)
      toast.success('Verification email sent! Check your inbox.', { duration: 5000 })
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to send verification email'
      const match = errorMsg.match(/(\d+) seconds/)
      if (match) {
        setResendCooldown(parseInt(match[1]))
        toast.error(`Please wait ${match[1]} seconds before requesting again`)
      } else {
        toast.error(errorMsg)
      }
    } finally {
      setResendLoading(false)
    }
  }

  if (isAuthenticated) {
    return <Navigate to={ROUTES.HOME} replace />
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center px-4 py-12">
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }} className="w-full max-w-md">
        <div className="text-center mb-8">
          <Link to={ROUTES.HOME} className="inline-flex items-center space-x-2 mb-6">
            <Code className="h-10 w-10 text-blue-600" />
            <span className="text-2xl font-bold text-gray-900">{CONFIG.APP_NAME}</span>
          </Link>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Welcome Back</h1>
          <p className="text-gray-600">Sign in to continue your coding journey</p>
        </div>

        <Card className="p-8">
          <AnimatePresence>
            {showResendOption && (
              <motion.div initial={{ opacity: 0, height: 0 }} animate={{ opacity: 1, height: 'auto' }} exit={{ opacity: 0, height: 0 }} className="mb-6 overflow-hidden">
                <div className="bg-amber-50 border border-amber-200 rounded-lg p-4">
                  <div className="flex items-start mb-3">
                    <AlertCircle className="h-5 w-5 text-amber-600 mt-0.5 mr-2 flex-shrink-0" />
                    <div className="flex-1">
                      <h3 className="text-sm font-semibold text-amber-900 mb-1">Email Not Verified</h3>
                      <p className="text-sm text-amber-800">Please verify your email address before logging in. Check your inbox for the verification link.</p>
                    </div>
                  </div>

                  <Button type="button" variant="outline" size="sm" onClick={handleResendVerification} disabled={resendCooldown > 0 || resendLoading} className="w-full text-amber-700 border-amber-300 hover:bg-amber-100">
                    {resendLoading ? (
                      <>
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                        Sending...
                      </>
                    ) : resendCooldown > 0 ? (
                      `Resend in ${Math.floor(resendCooldown / 60)}:${(resendCooldown % 60).toString().padStart(2, '0')}`
                    ) : (
                      <>
                        <Mail className="h-4 w-4 mr-2" />
                        Resend Verification Email
                      </>
                    )}
                  </Button>
                </div>
              </motion.div>
            )}
          </AnimatePresence>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <Input label="Email" type="email" placeholder="john@example.com" error={errors.email?.message} {...register('email')} />

            <div>
              <Input label="Password" type="password" placeholder="••••••••" error={errors.password?.message} {...register('password')} />

              <div className="mt-2 text-right">
                <p className="text-sm text-gray-600">
                  <Link to={ROUTES.FORGOT_PASSWORD} className="text-blue-600 hover:text-blue-700 font-medium">
                    Forgot Password?
                  </Link>
                </p>
              </div>
            </div>

            <Button type="submit" variant="primary" className="w-full" isLoading={isPending}>
              Sign In
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Don't have an account?{' '}
              <Link to={ROUTES.REGISTER} className="text-blue-600 hover:text-blue-700 font-medium">
                Sign up
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
