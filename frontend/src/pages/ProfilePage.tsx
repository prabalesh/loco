import { motion } from 'framer-motion'
import { useProfile } from '@/features/auth/hooks/useProfile'
import { Card } from '@/shared/components/ui/Card'
import { Button } from '@/shared/components/ui/Button'
import { formatDistanceToNow } from 'date-fns'
import { Skeleton } from '@/shared/components/ui/Skeleton'
import { StatsHeatmap } from '@/features/users/components/StatsHeatmap'
import { SolvedDistribution } from '@/features/users/components/SolvedDistribution'
import { achievementsApi } from '@/features/achievements/api/achievementsApi'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import {
  Trophy as TrophyIcon,
  CheckCircle2,
  Activity as ActivityIcon,
  User,
  Mail,
  Calendar,
  Shield,
  XCircle,
  Zap,
  Award,
  ArrowRight
} from 'lucide-react'

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
  const { data: userAchievements = [] } = useQuery({
    queryKey: ['my-achievements'],
    queryFn: achievementsApi.getMyAchievements,
    enabled: !!user,
  })

  const level = user?.level || 1
  const xp = user?.xp || 0
  const xpForNextLevel = level * 100
  const xpProgress = (xp % 100) / 100 * 100
  const recentAchievements = userAchievements.slice(0, 3)

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
          <Card className="p-4 sm:p-6 lg:p-8 mb-4 sm:mb-6 border-0 shadow-xl shadow-gray-200/50 rounded-[2.5rem] overflow-hidden relative">
            <div className="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-blue-500 to-indigo-500" />
            <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-6">
              <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6 w-full sm:w-auto">
                {/* Avatar */}
                <div className="relative">
                  <div className="h-24 w-24 bg-gradient-to-br from-blue-100 to-indigo-100 rounded-3xl flex items-center justify-center border-4 border-white shadow-xl">
                    <User className="h-12 w-12 text-blue-600" />
                  </div>
                  {user.is_verified && (
                    <div className="absolute -bottom-1 -right-1 bg-white p-1 rounded-full shadow-md">
                      <CheckCircle2 className="h-5 w-5 text-emerald-500 fill-emerald-50" />
                    </div>
                  )}
                </div>

                {/* User Info */}
                <div className="text-center sm:text-left w-full sm:w-auto">
                  <h1 className="text-3xl font-black text-gray-900 break-words mb-1">
                    @{user.username}
                  </h1>
                  <div className="flex items-center justify-center sm:justify-start space-x-2 text-gray-500">
                    <Mail className="h-4 w-4 flex-shrink-0" />
                    <span className="text-sm font-semibold">{user.email}</span>
                  </div>
                  <div className="flex flex-wrap justify-center sm:justify-start gap-2 mt-4">
                    <span className="inline-flex items-center px-4 py-1.5 rounded-xl text-xs font-bold bg-blue-50 text-blue-600 uppercase tracking-widest border border-blue-100">
                      <Shield className="h-3 w-3 mr-2" />
                      {user.role}
                    </span>
                    <span className="inline-flex items-center px-4 py-1.5 rounded-xl text-xs font-bold bg-gray-50 text-gray-500 uppercase tracking-widest border border-gray-100">
                      <Calendar className="h-3 w-3 mr-2" />
                      Joined {new Date(user.created_at).toLocaleDateString(undefined, { month: 'long', year: 'numeric' })}
                    </span>
                  </div>
                </div>
              </div>

              {/* Edit Button */}
              <Button variant="outline" className="rounded-2xl border-2 font-bold hover:bg-gray-50 px-8">
                Edit Profile
              </Button>
            </div>
          </Card>

          {/* XP and Level Card */}
          <Card className="p-6 border-0 shadow-xl shadow-purple-200/30 rounded-[2.5rem] bg-gradient-to-br from-purple-600 to-pink-600 text-white mb-6 relative overflow-hidden">
            <div className="absolute top-0 right-0 p-8 opacity-10">
              <Zap className="h-32 w-32" />
            </div>
            <div className="relative z-10">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <div className="flex items-center gap-2 mb-1">
                    <Award className="h-5 w-5" />
                    <span className="text-sm font-bold text-purple-100">LEVEL</span>
                  </div>
                  <div className="text-5xl font-black">{level}</div>
                </div>
                <div className="text-right">
                  <div className="flex items-center gap-2 mb-1 justify-end">
                    <Zap className="h-5 w-5" />
                    <span className="text-sm font-bold text-purple-100">EXPERIENCE</span>
                  </div>
                  <div className="text-3xl font-black">{xp.toLocaleString()}</div>
                </div>
              </div>
              <div className="space-y-2">
                <div className="flex justify-between text-sm font-semibold">
                  <span className="text-purple-100">Progress to Level {level + 1}</span>
                  <span>{xp % 100} / {xpForNextLevel}</span>
                </div>
                <div className="h-3 w-full bg-white/20 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-white transition-all duration-500 rounded-full"
                    style={{ width: `${xpProgress}%` }}
                  />
                </div>
              </div>
            </div>
          </Card>

          {/* Stats Grid */}
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-8">
            <motion.div whileHover={{ y: -5 }}>
              <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-blue-50 hover:ring-blue-200 transition-all">
                <div className="p-2 bg-blue-50 rounded-xl group-hover:bg-blue-500 group-hover:text-white transition-colors text-blue-600 w-fit mb-3">
                  <TrophyIcon className="h-4 w-4" />
                </div>
                <div className="text-2xl font-black text-gray-900">#{user.stats?.rank || 'N/A'}</div>
                <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Global Rank</div>
              </Card>
            </motion.div>

            <motion.div whileHover={{ y: -5 }}>
              <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-emerald-50 hover:ring-emerald-200 transition-all">
                <div className="p-2 bg-emerald-50 rounded-xl group-hover:bg-emerald-500 group-hover:text-white transition-colors text-emerald-600 w-fit mb-3">
                  <CheckCircle2 className="h-4 w-4" />
                </div>
                <div className="text-2xl font-black text-gray-900">{user.stats?.problems_solved || 0}</div>
                <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Solved</div>
              </Card>
            </motion.div>

            <motion.div whileHover={{ y: -5 }}>
              <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-amber-50 hover:ring-amber-200 transition-all">
                <div className="p-2 bg-amber-50 rounded-xl group-hover:bg-amber-500 group-hover:text-white transition-colors text-amber-600 w-fit mb-3">
                  <ActivityIcon className="h-4 w-4" />
                </div>
                <div className="text-2xl font-black text-gray-900">{(user.stats?.acceptance_rate || 0).toFixed(1)}%</div>
                <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Accuracy</div>
              </Card>
            </motion.div>

            <motion.div whileHover={{ y: -5 }}>
              <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-indigo-50 hover:ring-indigo-200 transition-all">
                <div className="p-2 bg-indigo-50 rounded-xl group-hover:bg-indigo-500 group-hover:text-white transition-colors text-indigo-600 w-fit mb-3">
                  <ActivityIcon className="h-4 w-4" />
                </div>
                <div className="text-2xl font-black text-gray-900">{user.stats?.total_submissions || 0}</div>
                <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Submissions</div>
              </Card>
            </motion.div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
            <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
              <SolvedDistribution
                distribution={user.solved_distribution || []}
                totalSolved={user.stats?.problems_solved || 0}
              />
            </Card>

            {/* Recent Achievements */}
            <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-bold text-gray-900">Recent Achievements</h3>
                <Link to="/achievements">
                  <Button variant="ghost" size="sm" className="text-purple-600 hover:text-purple-700">
                    View All
                    <ArrowRight className="h-4 w-4 ml-1" />
                  </Button>
                </Link>
              </div>
              {recentAchievements.length > 0 ? (
                <div className="space-y-3">
                  {recentAchievements.map((ua) => (
                    <div
                      key={ua.id}
                      className="flex items-center gap-3 p-3 bg-gradient-to-r from-purple-50 to-pink-50 rounded-xl border border-purple-100"
                    >
                      <div className="p-2 bg-gradient-to-br from-purple-500 to-pink-500 rounded-lg">
                        <TrophyIcon className="h-4 w-4 text-white" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="font-bold text-sm text-gray-900 truncate">
                          {ua.achievement.name}
                        </p>
                        <p className="text-xs text-gray-600 truncate">
                          +{ua.achievement.xp_reward} XP
                        </p>
                      </div>
                      <span className="text-xs text-purple-600 font-semibold">
                        {new Date(ua.unlocked_at).toLocaleDateString()}
                      </span>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <TrophyIcon className="h-12 w-12 text-gray-300 mx-auto mb-3" />
                  <p className="text-sm text-gray-500 font-medium">No achievements yet</p>
                  <p className="text-xs text-gray-400 mt-1">Start solving problems to unlock badges!</p>
                </div>
              )}
            </Card>
          </div>

          <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white mb-8">
            <StatsHeatmap data={user.submission_heatmap || []} />
          </Card>

          {/* Account Details */}
          <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
            <h2 className="text-2xl font-bold text-gray-900 mb-8">Account Details</h2>
            <div className="space-y-1">
              {[
                { label: 'Username', value: user.username, icon: User },
                { label: 'Email', value: user.email, icon: Mail },
                { label: 'Role', value: user.role, icon: Shield, capitalize: true },
                { label: 'Member Since', value: formatDistanceToNow(new Date(user.created_at), { addSuffix: true }), icon: Calendar },
                { label: 'Email Verified', value: user.email_verified ? 'Yes' : 'No', icon: CheckCircle2, success: user.email_verified, error: !user.email_verified }
              ].map((item, idx) => (
                <div key={item.label} className={`flex items-center justify-between py-4 ${idx !== 4 ? 'border-b border-gray-100' : ''}`}>
                  <div className="flex items-center space-x-4">
                    <div className="p-2 bg-gray-50 rounded-lg text-gray-400">
                      <item.icon className="h-5 w-5" />
                    </div>
                    <span className="font-bold text-gray-500 uppercase tracking-tight text-sm">{item.label}</span>
                  </div>
                  <span className={`font-black text-gray-900 ${item.capitalize ? 'capitalize' : ''} ${item.success ? 'text-emerald-600' : item.error ? 'text-rose-600' : ''}`}>
                    {item.value}
                  </span>
                </div>
              ))}
            </div>
          </Card>
        </motion.div>
      </div>
    </div>
  )
}
