import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from './../../../shared/components/ui/Input'
import { Button } from './../../../shared/components/ui/Button'
import { useRegister } from '../hooks/useRegister'

const registerSchema = z.object({
  email: z
    .string()
    .min(1, 'Email is required')
    .email('Invalid email format'),
  username: z
    .string()
    .min(3, 'Username must be at least 3 characters')
    .max(50, 'Username must be less than 50 characters')
    .regex(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers, and underscores'),
  password: z
    .string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
    .regex(/[0-9]/, 'Password must contain at least one number'),
})

type RegisterFormData = z.infer<typeof registerSchema>

export const RegisterForm = () => {
  const { mutate: register, isPending } = useRegister()

  const {
    register: registerField,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  })

  const onSubmit = (data: RegisterFormData) => {
    register(data)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <Input
        label="Email"
        type="email"
        placeholder="john@example.com"
        error={errors.email?.message}
        {...registerField('email')}
      />

      <Input
        label="Username"
        type="text"
        placeholder="john_doe"
        error={errors.username?.message}
        {...registerField('username')}
      />

      <Input
        label="Password"
        type="password"
        placeholder="••••••••"
        error={errors.password?.message}
        {...registerField('password')}
      />

      <Button
        type="submit"
        variant="primary"
        className="w-full"
        isLoading={isPending}
      >
        Create Account
      </Button>
    </form>
  )
}
