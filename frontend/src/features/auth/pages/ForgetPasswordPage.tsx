import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { authApi } from '../api/authApi'
import { Button } from '@/shared/components/ui/Button'
import { Input } from '@/shared/components/ui/Input'
import { Card } from '@/shared/components/ui/Card'
import { toast } from 'react-hot-toast'

const forgotPasswordSchema = z.object({
  email: z.string().min(1, 'Email is required').email('Invalid email format'),
})

type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>

export const ForgotPasswordPage = () => {
  const [success, setSuccess] = useState(false)
  const { register, handleSubmit, formState: { errors }, reset } = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordSchema),
  })

  const onSubmit = async (data: ForgotPasswordFormData) => {
    try {
      await authApi.forgotPassword(data.email)
      setSuccess(true)
      reset()
      toast.success('If the email exists, a reset link has been sent.')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to send reset email')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
      <Card className="max-w-md w-full p-8">
        <h1 className="text-2xl font-bold mb-4">Forgot Password</h1>
        {success ? (
          <p className="text-green-700">
            If an account with that email exists, a password reset link has been sent.
          </p>
        ) : (
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <Input
              label="Email"
              type="email"
              placeholder="your.email@example.com"
              error={errors.email?.message}
              {...register('email')}
            />
            <Button type="submit" variant="primary" className="w-full">
              Send Reset Link
            </Button>
          </form>
        )}
      </Card>
    </div>
  )
}
