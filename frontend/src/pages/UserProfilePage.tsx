import { useParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { User, Calendar, ArrowLeft } from 'lucide-react'
import { useUserProfile } from '@/features/auth/hooks/useUserProfile'
import { useAuth } from '@/shared/hooks/useAuth'
import { Card } from '@/shared/components/ui/Card'
import { Button } from '@/shared/components/ui/Button'
import { Loading } from '@/shared/components/common/Loading'
import { ROUTES } from '@/shared/constants/routes'
import { formatDistanceToNow, parseISO, format } from 'date-fns'

export const UserProfilePage = () => {
  const { username } = useParams<{ username: string }>()
  const { user: currentUser } = useAuth()
  const { data: user, isLoading, error } = useUserProfile(username!)

  const isOwnProfile = currentUser?.username === username

  if (isLoading) {
    return <Loading />
  }

  if (error || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Card className="p-8 text-center max-w-md">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">User Not Found</h2>
          <p className="text-gray-600 mb-4">
            The user @{username} doesn't exist
          </p>
          <Link to={ROUTES.HOME}>
            <Button variant="primary">Go Home</Button>
          </Link>
        </Card>
      </div>
    )
  }

  // Check if this is public profile (limited data)
  const isPublicProfile = !('email' in user)

  return (
    <div className="min-h-screen bg-gray-50 py-12">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <Link to={ROUTES.LEADERBOARD} className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6">
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Leaderboard
        </Link>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
        >
          {/* Header Card */}
          <Card className="p-8 mb-6">
            <div className="flex items-start justify-between">
              <div className="flex items-center space-x-4">
                <div className="w-20 h-20 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center">
                  <User className="h-10 w-10 text-white" />
                </div>

                <div>
                  <h1 className="text-3xl font-bold text-gray-900">{user.username}</h1>
                  {/* ⭐ REMOVED email - only show username */}
                  <p className="text-gray-600 mt-1">@{user.username}</p>
                </div>
              </div>

              {isOwnProfile && (
                <Link to={ROUTES.PROFILE}>
                  <Button variant="outline" size="sm">
                    Edit Profile
                  </Button>
                </Link>
              )}
            </div>

            {/* Badges */}
            <div className="flex flex-wrap gap-2 mt-6">
              
              {/* Only show verification for own profile */}
              {user.is_verified && (
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                  Verified
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

          {/* About */}
          <Card className="p-8">
            <h2 className="text-xl font-bold text-gray-900 mb-6">About</h2>

            <div className="space-y-4">
              <div className="flex items-center justify-between py-3 border-b border-gray-200">
                <div className="flex items-center space-x-3">
                  <Calendar className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Member Since</span>
                </div>
                <span className="text-gray-900">
                  {format(parseISO(user.created_at), 'MMMM d, yyyy')}
                </span>
              </div>

              <div className="flex items-center justify-between py-3 border-b border-gray-200">
                <div className="flex items-center space-x-3">
                  <Calendar className="h-5 w-5 text-gray-500" />
                  <span className="text-gray-600 font-medium">Joined</span>
                </div>
                <span className="text-gray-900">
                  {/* ⭐ FIXED: Use parseISO */}
                  {formatDistanceToNow(parseISO(user.created_at), { addSuffix: true })}
                </span>
              </div>

              {/* Only show account status for own profile */}
              {!isPublicProfile && 'is_active' in user && (
                <div className="flex items-center justify-between py-3 border-t border-gray-200">
                  <span className="text-gray-600 font-medium">Account Status</span>
                  <span className={user.is_active ? 'text-green-600' : 'text-red-600'}>
                    {user.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              )}
            </div>
          </Card>

          {/* Recent Submissions */}
          <Card className="p-8 mt-6">
            <h2 className="text-xl font-bold text-gray-900 mb-6">Recent Submissions</h2>
            <div className="text-center py-12">
              <Calendar className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-600">No recent submissions</p>
            </div>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}
