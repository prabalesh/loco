import { useParams, Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { User, Calendar, ArrowLeft, Award, Trophy, Flame, Target, Clock } from 'lucide-react'
import { useUserProfile } from '@/features/auth/hooks/useUserProfile'
import { useAuth } from '@/shared/hooks/useAuth'
import { Card } from '@/shared/components/ui/Card'
import { Button } from '@/shared/components/ui/Button'
import { Loading } from '@/shared/components/common/Loading'
import { ROUTES } from '@/shared/constants/routes'
import { formatDistanceToNow, parseISO, format } from 'date-fns'
import { achievementsApi } from '@/features/achievements/api/achievementsApi'
import { useQuery } from '@tanstack/react-query'
import { cn } from '@/shared/lib/utils'

export const UserProfilePage = () => {
  const { username } = useParams<{ username: string }>()
  const { user: currentUser } = useAuth()
  const { data: user, isLoading, error } = useUserProfile(username!)

  const { data: userAchievements = [] } = useQuery({
    queryKey: ['user-achievements', username],
    queryFn: () => achievementsApi.getUserAchievements(username!),
    enabled: !!username,
  })

  const isOwnProfile = currentUser?.username === username

  const level = user?.level || 1
  const xp = user?.xp || 0
  const xpProgress = (xp % 100)

  const stats = user?.stats
  const streak = stats?.streak || 0

  // Prepare distribution data
  const distribution = stats?.solved_distribution || []
  const easyCount = distribution.find(d => d.difficulty === 'Easy')?.count || 0
  const mediumCount = distribution.find(d => d.difficulty === 'Medium')?.count || 0
  const hardCount = distribution.find(d => d.difficulty === 'Hard')?.count || 0
  const totalSolved = stats?.problems_solved || 0

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50/50">
        <Loading />
      </div>
    )
  }

  if (error || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50/50 p-4">
        <Card className="p-8 text-center max-w-md border-0 shadow-2xl shadow-gray-200/50 rounded-3xl">
          <div className="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <User className="h-10 w-10 text-gray-400" />
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">User Not Found</h2>
          <p className="text-gray-600 mb-8">
            The explorer @{username} hasn't joined our universe yet.
          </p>
          <Link to={ROUTES.HOME}>
            <Button variant="primary" className="w-full py-6 rounded-2xl">Return Home</Button>
          </Link>
        </Card>
      </div>
    )
  }


  return (
    <div className="min-h-screen bg-[#FDFDFF] py-12">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, x: -10 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.3 }}
        >
          <Link to={ROUTES.LEADERBOARD} className="inline-flex items-center text-sm font-semibold text-gray-500 hover:text-indigo-600 mb-8 transition-colors group">
            <ArrowLeft className="h-4 w-4 mr-2 group-hover:-translate-x-1 transition-transform" />
            BACK TO LEADERBOARD
          </Link>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Column: Profile Card & Quick Stats */}
          <div className="lg:col-span-1 space-y-6">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.4 }}
            >
              <Card className="p-8 border-0 shadow-2xl shadow-indigo-100/50 rounded-[2.5rem] bg-white text-center relative overflow-hidden group">
                <div className="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500" />

                <div className="relative mb-6 mx-auto w-32 h-32 p-1.5 rounded-full bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500">
                  <div className="w-full h-full rounded-full bg-white flex items-center justify-center overflow-hidden border-4 border-white">
                    <User className="h-16 w-16 text-indigo-400" />
                  </div>
                  <div className="absolute -bottom-1 -right-1 w-10 h-10 bg-indigo-600 rounded-full border-4 border-white flex items-center justify-center shadow-lg">
                    <Target className="h-5 w-5 text-white" />
                  </div>
                </div>

                <h1 className="text-3xl font-black text-gray-900 leading-none">{user.username}</h1>
                <p className="text-indigo-500 font-bold tracking-wider mt-2 uppercase text-xs">CODE MASTER</p>

                <div className="mt-8 pt-8 border-t border-gray-100 grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-2xl font-black text-gray-900">{totalSolved}</p>
                    <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest">SOLVED</p>
                  </div>
                  <div className="border-l border-gray-100">
                    <p className="text-2xl font-black text-indigo-600">{stats?.rank || 'N/A'}</p>
                    <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest">RANK</p>
                  </div>
                </div>

                {isOwnProfile && (
                  <Link to={ROUTES.PROFILE} className="block mt-8">
                    <Button variant="outline" className="w-full rounded-2xl border-2 border-gray-100 hover:border-indigo-100 hover:bg-indigo-50 transition-all font-bold">
                      Edit Profile
                    </Button>
                  </Link>
                )}
              </Card>
            </motion.div>

            {/* Streak Card */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: 0.1 }}
            >
              <Card className={cn(
                "p-6 border-0 shadow-2xl rounded-[2rem] text-white flex items-center justify-between overflow-hidden relative",
                streak > 0 ? "bg-gradient-to-br from-orange-500 to-red-600" : "bg-gradient-to-br from-gray-400 to-gray-500 opacity-80"
              )}>
                <div className="absolute -right-4 -bottom-4 opacity-10">
                  <Flame className="h-24 w-24" />
                </div>
                <div className="relative z-10 flex items-center gap-4">
                  <div className="p-3 bg-white/20 rounded-2xl backdrop-blur-md">
                    <Flame className={cn("h-8 w-8", streak > 0 ? "text-white animate-pulse" : "text-white/50")} />
                  </div>
                  <div>
                    <h3 className="text-sm font-bold text-white/80 tracking-widest uppercase">STREAK</h3>
                    <p className="text-4xl font-black leading-none">{streak} <span className="text-sm font-bold text-white/60">DAYS</span></p>
                  </div>
                </div>
                {streak > 5 && (
                  <div className="relative z-10 p-2 bg-yellow-400 rounded-full shadow-lg">
                    <Target className="h-4 w-4 text-orange-700" />
                  </div>
                )}
              </Card>
            </motion.div>

            {/* About / Info */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.4, delay: 0.2 }}
            >
              <Card className="p-8 border-0 shadow-xl shadow-gray-100/50 rounded-[2.5rem] bg-white">
                <h2 className="text-xs font-black text-gray-400 uppercase tracking-[0.2em] mb-6">EXPLORER IDENTITY</h2>
                <div className="space-y-6">
                  <div className="flex items-center gap-4">
                    <div className="p-2.5 bg-blue-50 text-blue-500 rounded-xl">
                      <Calendar className="h-5 w-5" />
                    </div>
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest">MEMBER SINCE</p>
                      <p className="text-sm font-black text-gray-900">{format(parseISO(user.created_at), 'MMM d, yyyy')}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-4">
                    <div className="p-2.5 bg-green-50 text-green-500 rounded-xl">
                      <Clock className="h-5 w-5" />
                    </div>
                    <div>
                      <p className="text-[10px] font-bold text-gray-400 uppercase tracking-widest">LAST ACTIVE</p>
                      <p className="text-sm font-black text-gray-900">{formatDistanceToNow(parseISO(user.created_at), { addSuffix: true })}</p>
                    </div>
                  </div>
                </div>
              </Card>
            </motion.div>
          </div>

          {/* Right Column: Progress & Achievements */}
          <div className="lg:col-span-2 space-y-8">
            {/* Level & XP Hero Card */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5 }}
            >
              <Card className="p-10 border-0 shadow-2xl shadow-indigo-100/50 rounded-[3rem] bg-gradient-to-br from-[#1E293B] to-[#0F172A] text-white relative overflow-hidden">
                <div className="absolute top-0 right-0 -mr-20 -mt-20 w-80 h-80 bg-indigo-500/10 rounded-full blur-3xl" />
                <div className="absolute bottom-0 left-0 -ml-20 -mb-20 w-80 h-80 bg-purple-500/10 rounded-full blur-3xl" />

                <div className="relative z-10">
                  <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-10">
                    <div className="flex items-center gap-6">
                      <div className="relative">
                        <div className="w-24 h-24 rounded-3xl bg-indigo-500/20 border border-indigo-400/30 flex items-center justify-center rotate-12 group-hover:rotate-0 transition-transform">
                          <Trophy className="h-12 w-12 text-indigo-400 -rotate-12 group-hover:rotate-0 transition-transform" />
                        </div>
                        <div className="absolute -top-3 -right-3 w-10 h-10 bg-gradient-to-tr from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center font-black text-xl shadow-xl shadow-indigo-500/30 border-2 border-[#1E293B]">
                          {level}
                        </div>
                      </div>
                      <div>
                        <h2 className="text-sm font-black text-indigo-400 uppercase tracking-[0.2em] mb-1">CURRENT RANKING</h2>
                        <p className="text-4xl font-black italic tracking-tight">Code Vanguard</p>
                      </div>
                    </div>

                    <div className="text-left md:text-right">
                      <h2 className="text-sm font-black text-purple-400 uppercase tracking-[0.2em] mb-1">TOTAL EXPERIENCE</h2>
                      <p className="text-5xl font-black text-white">{xp.toLocaleString()} <span className="text-lg font-bold text-gray-500 italic">XP</span></p>
                    </div>
                  </div>

                  <div className="space-y-4">
                    <div className="flex justify-between items-end">
                      <p className="text-sm font-bold text-gray-400 uppercase tracking-widest">GOAL: LEVEL {level + 1}</p>
                      <p className="text-lg font-black">{xpProgress}% <span className="text-xs font-bold text-gray-500">COMPLETE</span></p>
                    </div>
                    <div className="h-4 w-full bg-white/5 rounded-full overflow-hidden p-1 shadow-inner">
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: `${xpProgress}%` }}
                        transition={{ duration: 1, ease: "easeOut" }}
                        className="h-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-full shadow-[0_0_20px_rgba(99,102,241,0.5)]"
                      />
                    </div>
                    <div className="flex justify-between text-[10px] font-black tracking-widest text-gray-500">
                      <span>{xp % 100} XP COLLECTED</span>
                      <span>{(level + 1) * 100} XP TARGET</span>
                    </div>
                  </div>
                </div>
              </Card>
            </motion.div>

            {/* Difficulty Mastery */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.2 }}
              >
                <Card className="p-8 border-0 shadow-xl shadow-gray-100/50 rounded-[2.5rem] bg-white h-full">
                  <div className="flex items-center justify-between mb-8">
                    <h2 className="text-sm font-black text-gray-400 uppercase tracking-[0.2em]">MASTERY MAP</h2>
                    <Target className="h-5 w-5 text-gray-300" />
                  </div>

                  <div className="space-y-8">
                    {/* Easy */}
                    <div className="space-y-2">
                      <div className="flex justify-between items-center font-bold">
                        <span className="text-xs text-green-500 tracking-widest uppercase">EASY EXPLORER</span>
                        <span className="text-sm text-gray-900 font-black">{easyCount}</span>
                      </div>
                      <div className="h-2 w-full bg-gray-50 rounded-full overflow-hidden">
                        <div className="h-full bg-green-500 rounded-full" style={{ width: `${totalSolved ? (easyCount / totalSolved) * 100 : 0}%` }} />
                      </div>
                    </div>

                    {/* Medium */}
                    <div className="space-y-2">
                      <div className="flex justify-between items-center font-bold">
                        <span className="text-xs text-yellow-500 tracking-widest uppercase">MEDIUM WARRIOR</span>
                        <span className="text-sm text-gray-900 font-black">{mediumCount}</span>
                      </div>
                      <div className="h-2 w-full bg-gray-50 rounded-full overflow-hidden">
                        <div className="h-full bg-yellow-500 rounded-full" style={{ width: `${totalSolved ? (mediumCount / totalSolved) * 100 : 0}%` }} />
                      </div>
                    </div>

                    {/* Hard */}
                    <div className="space-y-2">
                      <div className="flex justify-between items-center font-bold">
                        <span className="text-xs text-red-500 tracking-widest uppercase">HARDCORE MASTER</span>
                        <span className="text-sm text-gray-900 font-black">{hardCount}</span>
                      </div>
                      <div className="h-2 w-full bg-gray-50 rounded-full overflow-hidden">
                        <div className="h-full bg-red-500 rounded-full" style={{ width: `${totalSolved ? (hardCount / totalSolved) * 100 : 0}%` }} />
                      </div>
                    </div>
                  </div>
                </Card>
              </motion.div>

              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.3 }}
              >
                <Card className="p-8 border-0 shadow-xl shadow-gray-100/50 rounded-[2.5rem] bg-white h-full">
                  <div className="flex items-center justify-between mb-8">
                    <h2 className="text-sm font-black text-gray-400 uppercase tracking-[0.2em]">STATISTICS</h2>
                    <Award className="h-5 w-5 text-gray-300" />
                  </div>

                  <div className="grid grid-cols-2 gap-6">
                    <div className="p-4 rounded-3xl bg-indigo-50/50 border border-indigo-50">
                      <p className="text-[10px] font-black text-indigo-400 uppercase tracking-widest mb-1">RATE</p>
                      <p className="text-2xl font-black text-indigo-900">{stats?.acceptance_rate?.toFixed(1) || 0}%</p>
                    </div>
                    <div className="p-4 rounded-3xl bg-pink-50/50 border border-pink-50">
                      <p className="text-[10px] font-black text-pink-400 uppercase tracking-widest mb-1">TOTAL</p>
                      <p className="text-2xl font-black text-pink-900">{stats?.total_submissions || 0}</p>
                    </div>
                    <div className="p-4 rounded-3xl bg-blue-50/50 border border-blue-50">
                      <p className="text-[10px] font-black text-blue-400 uppercase tracking-widest mb-1">ACCEPTED</p>
                      <p className="text-2xl font-black text-blue-900">{stats?.accepted_submissions || 0}</p>
                    </div>
                    <div className="p-4 rounded-3xl bg-orange-50/50 border border-orange-50">
                      <p className="text-[10px] font-black text-orange-400 uppercase tracking-widest mb-1">RANK</p>
                      <p className="text-2xl font-black text-orange-900">#{stats?.rank || '?'}</p>
                    </div>
                  </div>
                </Card>
              </motion.div>
            </div>

            {/* Achievements Section */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.4 }}
            >
              <Card className="p-8 border-0 shadow-2xl shadow-indigo-50/50 rounded-[3rem] bg-white">
                <div className="flex items-center justify-between mb-10">
                  <div>
                    <h2 className="text-sm font-black text-gray-400 uppercase tracking-[0.2em] mb-1">TROPHY CABINET</h2>
                    <p className="text-3xl font-black text-gray-900">Achievements</p>
                  </div>
                  <Link to={ROUTES.ACHIEVEMENTS} className="p-3 bg-indigo-50 text-indigo-600 rounded-2xl hover:bg-indigo-600 hover:text-white transition-all">
                    <Target className="h-6 w-6" />
                  </Link>
                </div>

                {userAchievements.length > 0 ? (
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {userAchievements.slice(0, 6).map((ua) => (
                      <motion.div
                        whileHover={{ scale: 1.02 }}
                        key={ua.id}
                        className="flex items-center gap-4 p-5 rounded-[2rem] bg-gray-50/50 border-2 border-transparent hover:border-indigo-100 hover:bg-white transition-all group"
                      >
                        <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center shrink-0 shadow-lg shadow-indigo-200 group-hover:rotate-6 transition-transform">
                          <Trophy className="h-8 w-8 text-white" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-black text-gray-900 leading-tight">
                            {ua.achievement.name}
                          </p>
                          <p className="text-[10px] font-black text-indigo-500 mt-1 uppercase tracking-widest bg-indigo-50 inline-block px-2 py-0.5 rounded-lg">
                            +{ua.achievement.xp_reward} XP
                          </p>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-20 bg-gray-50/50 rounded-[3rem] border-2 border-dashed border-gray-200">
                    <Trophy className="h-20 w-20 text-gray-300 mx-auto mb-6" />
                    <p className="text-xl font-black text-gray-900 tracking-tight">NO TROPHIES EARNED YET</p>
                    <p className="text-gray-500 font-medium mt-2">Solve challenges to fill this cabinet!</p>
                  </div>
                )}

                {userAchievements.length > 6 && (
                  <div className="mt-8 text-center">
                    <Link to={ROUTES.ACHIEVEMENTS}>
                      <Button variant="ghost" className="font-black text-indigo-600 hover:bg-indigo-50 rounded-2xl">
                        VIEW ALL ACHIEVEMENTS ({userAchievements.length})
                      </Button>
                    </Link>
                  </div>
                )}
              </Card>
            </motion.div>
          </div>
        </div>
      </div>
    </div>
  )
}
