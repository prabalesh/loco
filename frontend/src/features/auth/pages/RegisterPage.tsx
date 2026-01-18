import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { Code } from 'lucide-react'
import { Input } from '@/shared/components/ui/Input'
import { Button } from '@/shared/components/ui/Button'
import { Card } from '@/shared/components/ui/Card'
import { authApi } from '../api/authApi'
import { ROUTES } from '@/shared/constants/routes'
import { CONFIG } from '@/shared/constants/config'
import { toast } from 'react-hot-toast'
import type { AxiosError } from 'axios'
import type { RegisterResponse } from '@/shared/types/auth.types'

const registerSchema = z.object({
  email: z.string().min(1, 'Email is required').email('Invalid email format'),
  username: z
    .string()
    .min(3, 'Username must be at least 3 characters')
    .max(20, 'Username must be less than 20 characters')
    .regex(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers, and underscores'),
  password: z
    .string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
    .regex(/[0-9]/, 'Password must contain at least one number'),
})

type RegisterFormData = z.infer<typeof registerSchema>

export const RegisterPage = () => {
  const [isLoading, setIsLoading] = useState(false)
  const [isSuccess, setIsSuccess] = useState(false)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  })

  const onSubmit = async (data: RegisterFormData) => {
    setIsLoading(true)
    setIsSuccess(false)
    setSuccessMessage(null)
    try {
      const response = await authApi.register(data)
      const message = response.data.message || 'Registration successful! Please verify your email.'
      setSuccessMessage(message)
      toast.success('Registration successful! Please verify your email.')
      setIsSuccess(true)
    } catch (err) {
      const error = err as AxiosError<RegisterResponse>
      const errorMsg = error.response?.data.data.message || 'Registration failed'
      toast.error(errorMsg)
      setIsSuccess(false)
    } finally {
      setIsLoading(false)
    }
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
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Create Account</h1>
          <p className="text-gray-600">Start your coding journey today</p>
        </div>

        <Card className="p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <Input
              label="Email"
              type="email"
              placeholder="john@example.com"
              error={errors.email?.message}
              {...register('email')}
            />

            <Input
              label="Username"
              type="text"
              placeholder="johndoe"
              error={errors.username?.message}
              {...register('username')}
            />

            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              error={errors.password?.message}
              {...register('password')}
            />

            <Button
              type="submit"
              variant="primary"
              className="w-full"
              isLoading={isLoading}
              disabled={isLoading || isSuccess}
            >
              {isLoading ? (
                'Creating Account...'
              ) : (
                'Sign Up'
              )}
            </Button>
          </form>

          {successMessage && (
            <div className="mt-6 text-green-700 font-semibold text-center select-none drop-shadow-md bg-green-100 px-4 py-3 rounded-md border border-green-300">
              <p>{successMessage}</p>
            </div>
          )}

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Already have an account?{' '}
              <Link to={ROUTES.LOGIN} className="text-blue-600 hover:text-blue-700 font-medium">
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
