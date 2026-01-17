import { motion } from 'framer-motion'
import { User, Mail, Calendar, Shield, CheckCircle, XCircle } from 'lucide-react'
import { useProfile } from '@/features/auth/hooks/useProfile'
import { Card } from '@/shared/components/ui/Card'
import { Button } from '@/shared/components/ui/Button'
import { formatDistanceToNow } from 'date-fns'
import { Skeleton } from '@/shared/components/ui/Skeleton'

const ProfileSkeleton = () => (
  <div className="min-h-screen bg-gray-50 py-6 sm:py-12">
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      {/* Header Card Skeleton */}
      <Card className="p-4 sm:p-6 lg:p-8 mb-4 sm:mb-6">
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <div className="flex flex-col sm:flex-row items-center sm:items-start gap-4 w-full sm:w-auto">
            <Skeleton className="w-16 h-16 sm:w-20 sm:h-20 rounded-full" />
            <div className="text-center sm:text-left w-full sm:w-auto space-y-2">
              <Skeleton className="h-8 w-48 mx-auto sm:mx-0" />
              <div className="flex items-center justify-center sm:justify-start space-x-2">
                <Skeleton className="h-4 w-4" />
                <Skeleton className="h-4 w-40" />
              </div>
            </div>
          </div>
          <Skeleton className="h-9 w-32 mx-auto sm:mx-0" />
        </div>
        <div className="flex flex-wrap justify-center sm:justify-start gap-2 mt-4 sm:mt-6">
          <Skeleton className="h-6 w-24 rounded-full" />
          <Skeleton className="h-6 w-32 rounded-full" />
        </div>
      </Card>

      {/* Stats Grid Skeleton */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 sm:gap-6 mb-4 sm:mb-6">
        {[1, 2, 3].map((i) => (
          <Card key={i} className="p-4 sm:p-6 text-center space-y-2">
            <Skeleton className="h-10 w-12 mx-auto" />
            <Skeleton className="h-5 w-32 mx-auto" />
          </Card>
        ))}
      </div>

      {/* Account Details Skeleton */}
      <Card className="p-4 sm:p-6 lg:p-8">
        <Skeleton className="h-7 w-48 mb-6" />
        <div className="space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-3 border-b border-gray-100 gap-2">
              <div className="flex items-center space-x-3">
                <Skeleton className="h-5 w-5" />
                <Skeleton className="h-5 w-32" />
              </div>
              <Skeleton className="h-5 w-48" />
            </div>
          ))}
        </div>
      </Card>
    </div>
  </div>
)

export const ProfilePage = () => {
  const { data: user, isLoading, error } = useProfile()

  if (isLoading) {
    return <ProfileSkeleton />
  }

  if (error || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center px-4">
        <Card className="p-6 sm:p-8 text-center max-w-md w-full">
          <XCircle className="h-12 w-12 sm:h-16 sm:w-16 text-red-500 mx-auto mb-4" />
          <h2 className="text-xl sm:text-2xl font-bold text-gray-900 mb-2">
            Failed to Load Profile
          </h2>
          <p className="text-sm sm:text-base text-gray-600">Please try again later</p>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-6 sm:py-12">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          {/* Header Card */}
          <Card className="p-4 sm:p-6 lg:p-8 mb-4 sm:mb-6">
            <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
              <div className="flex flex-col sm:flex-row items-center sm:items-start gap-4 w-full sm:w-auto">
                {/* Avatar */}
                <div className="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center flex-shrink-0">
                  <User className="h-8 w-8 sm:h-10 sm:w-10 text-white" />
                </div>

                {/* User Info */}
                <div className="text-center sm:text-left w-full sm:w-auto">
                  <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 break-words">
                    {user.username}
                  </h1>
                  <div className="flex items-center justify-center sm:justify-start space-x-2 mt-1">
                    <Mail className="h-4 w-4 text-gray-500 flex-shrink-0" />
                    <span className="text-sm sm:text-base text-gray-600 break-all">{user.email}</span>
                  </div>
                </div>
              </div>

              {/* Edit Button */}
              <Button variant="outline" size="sm" className="w-full sm:w-auto">
                Edit Profile
              </Button>
            </div>

            {/* Badges */}
            <div className="flex flex-wrap justify-center sm:justify-start gap-2 mt-4 sm:mt-6">
              {/* Role Badge */}
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs sm:text-sm font-medium bg-blue-100 text-blue-800">
                <Shield className="h-3 w-3 sm:h-4 sm:w-4 mr-1" />
                {user.role}
              </span>

              {/* Email Verified Badge */}
              {user.email_verified ? (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-xs sm:text-sm font-medium bg-green-100 text-green-800">
                  <CheckCircle className="h-3 w-3 sm:h-4 sm:w-4 mr-1" />
                  Email Verified
                </span>
              ) : (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-xs sm:text-sm font-medium bg-yellow-100 text-yellow-800">
                  <XCircle className="h-3 w-3 sm:h-4 sm:w-4 mr-1" />
                  Email Not Verified
                </span>
              )}
            </div>
          </Card>

          {/* Stats Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 sm:gap-6 mb-4 sm:mb-6">
            <Card className="p-4 sm:p-6 text-center">
              <div className="text-2xl sm:text-3xl font-bold text-blue-600 mb-1 sm:mb-2">0</div>
              <div className="text-sm sm:text-base text-gray-600 font-medium">Problems Solved</div>
            </Card>

            <Card className="p-4 sm:p-6 text-center">
              <div className="text-2xl sm:text-3xl font-bold text-green-600 mb-1 sm:mb-2">0</div>
              <div className="text-sm sm:text-base text-gray-600 font-medium">Submissions</div>
            </Card>

            <Card className="p-4 sm:p-6 text-center">
              <div className="text-2xl sm:text-3xl font-bold text-purple-600 mb-1 sm:mb-2">-</div>
              <div className="text-sm sm:text-base text-gray-600 font-medium">Global Rank</div>
            </Card>
          </div>

          {/* Account Details */}
          <Card className="p-4 sm:p-6 lg:p-8">
            <h2 className="text-lg sm:text-xl font-bold text-gray-900 mb-4 sm:mb-6">
              Account Details
            </h2>

            <div className="space-y-3 sm:space-y-4">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-3 border-b border-gray-200 gap-2 sm:gap-4">
                <div className="flex items-center space-x-3">
                  <User className="h-4 w-4 sm:h-5 sm:w-5 text-gray-500 flex-shrink-0" />
                  <span className="text-sm sm:text-base text-gray-600 font-medium">Username</span>
                </div>
                <span className="text-sm sm:text-base text-gray-900 font-semibold break-words sm:text-right">
                  {user.username}
                </span>
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-3 border-b border-gray-200 gap-2 sm:gap-4">
                <div className="flex items-center space-x-3">
                  <Mail className="h-4 w-4 sm:h-5 sm:w-5 text-gray-500 flex-shrink-0" />
                  <span className="text-sm sm:text-base text-gray-600 font-medium">Email</span>
                </div>
                <span className="text-sm sm:text-base text-gray-900 break-all sm:text-right">
                  {user.email}
                </span>
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-3 border-b border-gray-200 gap-2 sm:gap-4">
                <div className="flex items-center space-x-3">
                  <Shield className="h-4 w-4 sm:h-5 sm:w-5 text-gray-500 flex-shrink-0" />
                  <span className="text-sm sm:text-base text-gray-600 font-medium">Role</span>
                </div>
                <span className="text-sm sm:text-base text-gray-900 capitalize sm:text-right">
                  {user.role}
                </span>
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-3 border-b border-gray-200 gap-2 sm:gap-4">
                <div className="flex items-center space-x-3">
                  <Calendar className="h-4 w-4 sm:h-5 sm:w-5 text-gray-500 flex-shrink-0" />
                  <span className="text-sm sm:text-base text-gray-600 font-medium">Member Since</span>
                </div>
                <span className="text-sm sm:text-base text-gray-900 sm:text-right">
                  {formatDistanceToNow(new Date(user.created_at), { addSuffix: true })}
                </span>
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between py-3">
                <div className="flex items-center space-x-3">
                  <CheckCircle className="h-4 w-4 sm:h-5 sm:w-5 text-gray-500 flex-shrink-0" />
                  <span className="text-sm sm:text-base text-gray-600 font-medium">Email Verified</span>
                </div>
                <span className={`text-sm sm:text-base sm:text-right ${user.email_verified ? 'text-green-600' : 'text-red-600'}`}>
                  {user.email_verified ? 'Yes' : 'No'}
                </span>
              </div>
            </div>
          </Card>

          {/* Recent Activity (Placeholder) */}
          <Card className="p-4 sm:p-6 lg:p-8 mt-4 sm:mt-6">
            <h2 className="text-lg sm:text-xl font-bold text-gray-900 mb-4 sm:mb-6">
              Recent Activity
            </h2>
            <div className="text-center py-8 sm:py-12">
              <Calendar className="h-12 w-12 sm:h-16 sm:w-16 text-gray-400 mx-auto mb-3 sm:mb-4" />
              <p className="text-sm sm:text-base text-gray-600">No recent activity</p>
              <p className="text-xs sm:text-sm text-gray-500 mt-2 px-4">
                Start solving problems to see your activity here
              </p>
            </div>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}
