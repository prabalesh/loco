import { motion } from 'framer-motion'
import { User, Mail, Calendar, Shield, CheckCircle, XCircle } from 'lucide-react'
import { useProfile } from '@/features/auth/hooks/useProfile'
import { Card } from '@/shared/components/ui/Card'
import { Button } from '@/shared/components/ui/Button'
import { Loading } from '@/shared/components/common/Loading'
import { formatDistanceToNow } from 'date-fns'

export const ProfilePage = () => {
  const { data: user, isLoading, error } = useProfile()

  if (isLoading) {
    return <Loading />
  }

  if (error || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Card className="p-8 text-center max-w-md">
          <XCircle className="h-16 w-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            Failed to Load Profile
          </h2>
          <p className="text-gray-600">Please try again later</p>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          {/* Header Card */}
          <Card className="p-8 mb-6">
            <div className="flex items-start justify-between">
              <div className="flex items-center space-x-4">
                {/* Avatar */}
                <div className="w-20 h-20 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center">
                  <User className="h-10 w-10 text-white" />
                </div>

                {/* User Info */}
                <div>
                  <h1 className="text-3xl font-bold text-gray-900">
                    {user.username}
                  </h1>
                  <div className="flex items-center space-x-2 mt-1">
                    <Mail className="h-4 w-4 text-gray-500" />
                    <span className="text-gray-600">{user.email}</span>
                  </div>
                </div>
              </div>

              {/* Edit Button (future feature) */}
              <Button variant="outline" size="sm">
                Edit Profile
              </Button>
            </div>

            {/* Badges */}
            <div className="flex flex-wrap gap-2 mt-6">
              {/* Role Badge */}
              <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                <Shield className="h-4 w-4 mr-1" />
                {user.role}
              </span>

              {/* Email Verified Badge */}
              {user.email_verified ? (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                  <CheckCircle className="h-4 w-4 mr-1" />
                  Email Verified
                </span>
              ) : (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800">
                  <XCircle className="h-4 w-4 mr-1" />
                  Email Not Verified
                </span>
              )}
            </div>
          </Card>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
            <Card className="p-6 text-center">
              <div className="text-3xl font-bold text-blue-600 mb-2">0</div>
              <div className="text-gray-600 font-medium">Problems Solved</div>
            </Card>

            <Card className="p-6 text-center">
              <div className="text-3xl font-bold text-green-600 mb-2">0</div>
              <div className="text-gray-600 font-medium">Submissions</div>
            </Card>

            <Card className="p-6 text-center">
              <div className="text-3xl font-bold text-purple-600 mb-2">-</div>
              <div className="text-gray-600 font-medium">Global Rank</div>
            </Card>
          </div>

          {/* Account Details */}
          <Card className="p-8">
            <h2 className="text-xl font-bold text-gray-900 mb-6">
              Account Details
            </h2>

            <div className="space-y-4">
              <div className="flex items-center justify-between py-3 border-b border-gray-200">
                <div className="flex items-center space-x-3">
                  <User className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Username</span>
                </div>
                <span className="text-gray-900 font-semibold">{user.username}</span>
              </div>

              <div className="flex items-center justify-between py-3 border-b border-gray-200">
                <div className="flex items-center space-x-3">
                  <Mail className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Email</span>
                </div>
                <span className="text-gray-900">{user.email}</span>
              </div>

              <div className="flex items-center justify-between py-3 border-b border-gray-200">
                <div className="flex items-center space-x-3">
                  <Shield className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Role</span>
                </div>
                <span className="text-gray-900 capitalize">{user.role}</span>
              </div>

              <div className="flex items-center justify-between py-3 border-b border-gray-200">
                <div className="flex items-center space-x-3">
                  <Calendar className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Member Since</span>
                </div>
                <span className="text-gray-900">
                  {formatDistanceToNow(new Date(user.created_at), { addSuffix: true })}
                </span>
              </div>

              <div className="flex items-center justify-between py-3">
                <div className="flex items-center space-x-3">
                  <CheckCircle className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Email Verified</span>
                </div>
                <span className={user.email_verified ? 'text-green-600' : 'text-red-600'}>
                  {user.email_verified ? 'Yes' : 'No'}
                </span>
              </div>
            </div>
          </Card>

          {/* Recent Activity (Placeholder) */}
          <Card className="p-8 mt-6">
            <h2 className="text-xl font-bold text-gray-900 mb-6">
              Recent Activity
            </h2>
            <div className="text-center py-12">
              <Calendar className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-600">No recent activity</p>
              <p className="text-gray-500 text-sm mt-2">
                Start solving problems to see your activity here
              </p>
            </div>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}
