import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { authApi } from '../api/authApi'
import { Button } from '@/shared/components/ui/Button'
import { Input } from '@/shared/components/ui/Input'
import { Card } from '@/shared/components/ui/Card'
import { toast } from 'react-hot-toast'
import { ROUTES } from '@/shared/constants/routes'

const resetPasswordSchema = z.object({
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ['confirmPassword'],
})

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>

export const ResetPasswordPage = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const token = searchParams.get('token') || ''
  const [submitting, setSubmitting] = useState(false)

  const { register, handleSubmit, formState: { errors } } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
  })

  useEffect(() => {
    if (!token) {
      toast.error('Invalid or missing password reset token.')
      navigate(ROUTES.FORGOT_PASSWORD)
    }
  }, [token, navigate])

  const onSubmit = async (data: ResetPasswordFormData) => {
    setSubmitting(true)
    try {
      await authApi.resetPassword(token, data.password)
      toast.success('Password reset successful! Please login.')
      navigate(ROUTES.LOGIN)
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to reset password')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
      <Card className="max-w-md w-full p-8">
        <h1 className="text-2xl font-bold mb-4">Reset Password</h1>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input
            label="New Password"
            type="password"
            placeholder="********"
            error={errors.password?.message}
            {...register('password')}
          />
          <Input
            label="Confirm New Password"
            type="password"
            placeholder="********"
            error={errors.confirmPassword?.message}
            {...register('confirmPassword')}
          />
          <Button type="submit" variant="primary" className="w-full" isLoading={submitting}>
            Reset Password
          </Button>
        </form>
      </Card>
    </div>
  )
}
