import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Form, Input, Button, Card, Typography } from 'antd'
import { LockOutlined, UserOutlined, SafetyCertificateOutlined } from '@ant-design/icons'
import { adminAuthApi } from '../../../api/adminApi'
import { useAuthStore } from '../store/authStore'
import toast from 'react-hot-toast'
import type { LoginCredentials } from '../../../types'

const { Title, Text } = Typography

export const AdminLogin = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const setUser = useAuthStore((state) => state.setUser)

  const onFinish = async (values: LoginCredentials) => {
  setLoading(true)
  try {
    const response = await adminAuthApi.login(values)
    const user = response.data
    
    setUser(user)
    
    toast.success('Login successful!')
    navigate('/', { replace: true })
  } catch (error: any) {
    // ... error handling
  } finally {
    setLoading(false)
  }
}


  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <Card className="w-full max-w-md shadow-xl">
        <div className="text-center mb-8">
          <SafetyCertificateOutlined className="text-6xl text-blue-600 mb-4" />
          <Title level={2} className="mb-2">Admin Portal</Title>
          <Text type="secondary">Sign in to access the admin dashboard</Text>
        </div>

        <Form
          name="admin-login"
          onFinish={onFinish}
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="email"
            rules={[
              { required: true, message: 'Please enter your email' },
              { type: 'email', message: 'Please enter a valid email' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="Admin Email"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: 'Please enter your password' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="Password"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              size="large"
            >
              Sign In as Admin
            </Button>
          </Form.Item>
        </Form>

        <div className="text-center mt-4">
          <Text type="secondary" className="text-xs">
            Authorized personnel only. All activities are logged.
          </Text>
        </div>
      </Card>
    </div>
  )
}
