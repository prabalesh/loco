import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { usersApi } from '../api/users'
import { Card } from '@/shared/components/ui/Card'
import { Skeleton } from '@/shared/components/ui/Skeleton'
import {
    User,
    Calendar,
    CheckCircle2,
    Trophy,
    Activity,
    ChevronLeft,
    Mail
} from 'lucide-react'
import { Button } from '@/shared/components/ui/Button'
import { motion } from 'framer-motion'
import { format } from 'date-fns'

import type { UserProfile } from '../types'
import { StatsHeatmap } from '../components/StatsHeatmap'
import { SolvedDistribution } from '../components/SolvedDistribution'

export const ProfilePage = () => {
    const { username } = useParams<{ username: string }>()
    const navigate = useNavigate()

    const { data: profile, isLoading, isError } = useQuery({
        queryKey: ['user-profile', username],
        queryFn: () => usersApi.getProfile(username!).then(res => (res.data as any).data as UserProfile),
        enabled: !!username
    })

    if (isLoading) {
        return (
            <div className="max-w-4xl mx-auto px-4 py-12">
                <Skeleton className="h-32 w-full mb-8 rounded-2xl" />
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    <Skeleton className="h-40 w-full rounded-2xl" />
                    <Skeleton className="h-40 w-full rounded-2xl" />
                    <Skeleton className="h-40 w-full rounded-2xl" />
                </div>
            </div>
        )
    }

    if (isError || !profile) {
        return (
            <div className="max-w-4xl mx-auto px-4 py-20 text-center">
                <User className="h-16 w-16 text-gray-300 mx-auto mb-4" />
                <h2 className="text-2xl font-bold text-gray-900 mb-2">User not found</h2>
                <Button variant="ghost" onClick={() => navigate(-1)}>
                    <ChevronLeft className="h-4 w-4 mr-2" />
                    Go Back
                </Button>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-gray-50/50 py-12">
            <div className="max-w-5xl mx-auto px-4">
                <Button
                    variant="ghost"
                    onClick={() => navigate(-1)}
                    className="mb-8 text-gray-500 hover:text-gray-900"
                >
                    <ChevronLeft className="h-4 w-4 mr-2" />
                    Back
                </Button>

                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="grid grid-cols-1 lg:grid-cols-3 gap-8"
                >
                    {/* Sidebar / Main Info */}
                    <div className="lg:col-span-1 space-y-6">
                        <Card className="p-8 text-center bg-white border-0 shadow-xl shadow-gray-200/50 rounded-3xl overflow-hidden relative">
                            <div className="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-blue-500 to-indigo-500" />
                            <div className="relative mb-6">
                                <div className="h-24 w-24 bg-gradient-to-br from-blue-100 to-indigo-100 rounded-3xl mx-auto flex items-center justify-center border-4 border-white shadow-lg">
                                    <User className="h-12 w-12 text-blue-600" />
                                </div>
                                {profile.is_verified && (
                                    <div className="absolute -bottom-1 -right-1 bg-white p-1 rounded-full shadow-md">
                                        <CheckCircle2 className="h-5 w-5 text-emerald-500 fill-emerald-50" />
                                    </div>
                                )}
                            </div>

                            <h1 className="text-2xl font-bold text-gray-900 mb-1">@{profile.username}</h1>
                            <p className="text-sm text-gray-500 mb-6 flex items-center justify-center gap-1.5">
                                <Calendar className="h-3.5 w-3.5" />
                                Joined {format(new Date(profile.created_at), 'MMMM yyyy')}
                            </p>

                            <div className="flex flex-col gap-2">
                                <Button variant="primary" className="w-full rounded-xl">
                                    Follow
                                </Button>
                                <Button variant="ghost" className="w-full rounded-xl gap-2">
                                    <Mail className="h-4 w-4" />
                                    Message
                                </Button>
                            </div>
                        </Card>

                        <Card className="p-6 bg-white border-0 shadow-lg shadow-gray-200/40 rounded-3xl">
                            <h3 className="text-sm font-bold text-gray-900 uppercase tracking-wider mb-4 px-1">Socials</h3>
                            <div className="space-y-3">
                                <div className="p-3 bg-gray-50 rounded-xl text-sm text-gray-500 text-center italic">
                                    No social links added yet.
                                </div>
                            </div>
                        </Card>
                    </div>

                    {/* Main Content / Stats */}
                    <div className="lg:col-span-2 space-y-8">
                        {/* Stats Overview */}
                        <div className="grid grid-cols-1 sm:grid-cols-4 gap-4">
                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-blue-50 hover:ring-blue-200 transition-all">
                                    <div className="p-2 bg-blue-50 rounded-xl group-hover:bg-blue-500 group-hover:text-white transition-colors text-blue-600 w-fit mb-3">
                                        <Trophy className="h-4 w-4" />
                                    </div>
                                    <div className="text-2xl font-black text-gray-900">#{profile.stats.rank || 'N/A'}</div>
                                    <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Global Rank</div>
                                </Card>
                            </motion.div>

                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-emerald-50 hover:ring-emerald-200 transition-all">
                                    <div className="p-2 bg-emerald-50 rounded-xl group-hover:bg-emerald-500 group-hover:text-white transition-colors text-emerald-600 w-fit mb-3">
                                        <CheckCircle2 className="h-4 w-4" />
                                    </div>
                                    <div className="text-2xl font-black text-gray-900">{profile.stats.problems_solved}</div>
                                    <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Problems Solved</div>
                                </Card>
                            </motion.div>

                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-amber-50 hover:ring-amber-200 transition-all">
                                    <div className="p-2 bg-amber-50 rounded-xl group-hover:bg-amber-500 group-hover:text-white transition-colors text-amber-600 w-fit mb-3">
                                        <Activity className="h-4 w-4" />
                                    </div>
                                    <div className="text-2xl font-black text-gray-900">{profile.stats.acceptance_rate.toFixed(1)}%</div>
                                    <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Accuracy</div>
                                </Card>
                            </motion.div>

                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-5 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-indigo-50 hover:ring-indigo-200 transition-all">
                                    <div className="p-2 bg-indigo-50 rounded-xl group-hover:bg-indigo-500 group-hover:text-white transition-colors text-indigo-600 w-fit mb-3">
                                        <Activity className="h-4 w-4" />
                                    </div>
                                    <div className="text-2xl font-black text-gray-900">{profile.stats.total_submissions}</div>
                                    <div className="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">Submissions</div>
                                </Card>
                            </motion.div>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                            <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
                                <SolvedDistribution
                                    distribution={profile.solved_distribution}
                                    totalSolved={profile.stats.problems_solved}
                                />
                            </Card>

                            <div className="space-y-8">
                                <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-indigo-600 text-white relative overflow-hidden">
                                    <div className="absolute top-0 right-0 p-8 opacity-10">
                                        <Trophy className="h-32 w-32" />
                                    </div>
                                    <h3 className="text-lg font-bold mb-2">Badge Progress</h3>
                                    <p className="text-indigo-100 text-sm mb-6 font-medium">Solve 50 problems to unlock the next achievement!</p>
                                    <div className="h-2 w-full bg-white/20 rounded-full overflow-hidden">
                                        <div className="h-full bg-white transition-all" style={{ width: `${Math.min((profile.stats.problems_solved / 50) * 100, 100)}%` }} />
                                    </div>
                                </Card>
                            </div>
                        </div>

                        <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
                            <StatsHeatmap data={profile.submission_heatmap} />
                        </Card>

                        {/* Achievements Section */}
                        <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
                            <div className="flex items-center justify-between mb-8">
                                <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-3">
                                    <Trophy className="h-6 w-6 text-amber-500" />
                                    Achievements
                                </h2>
                                <span className="px-4 py-1.5 bg-gray-50 text-gray-500 rounded-xl text-sm font-bold">
                                    {profile.achievements?.length || 0} Unlocked
                                </span>
                            </div>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                {profile.achievements && profile.achievements.length > 0 ? (
                                    profile.achievements.map((ua, index) => (
                                        <motion.div
                                            key={ua.id}
                                            initial={{ opacity: 0, scale: 0.9 }}
                                            animate={{ opacity: 1, scale: 1 }}
                                            transition={{ delay: index * 0.05 }}
                                            className="group relative p-4 bg-gradient-to-br from-gray-50 to-white rounded-2xl border border-gray-100 hover:border-amber-200 hover:shadow-lg transition-all"
                                        >
                                            <div className="flex items-center gap-4">
                                                <div className="h-12 w-12 rounded-xl bg-amber-50 flex items-center justify-center text-amber-600 group-hover:scale-110 transition-transform">
                                                    {ua.achievement.icon_url ? (
                                                        <img src={ua.achievement.icon_url} alt={ua.achievement.name} className="h-8 w-8 object-contain" />
                                                    ) : (
                                                        <Trophy className="h-6 w-6" />
                                                    )}
                                                </div>
                                                <div className="flex-1 min-w-0">
                                                    <h4 className="font-bold text-gray-900 group-hover:text-amber-700 transition-colors truncate">
                                                        {ua.achievement.name}
                                                    </h4>
                                                    <p className="text-xs text-gray-500 line-clamp-1">{ua.achievement.description}</p>
                                                </div>
                                            </div>
                                            <div className="mt-3 flex items-center justify-between border-t border-gray-50 pt-3">
                                                <span className="text-[10px] font-bold text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full uppercase tracking-wider">
                                                    +{ua.achievement.xp_reward} XP
                                                </span>
                                                <span className="text-[10px] text-gray-400 font-medium">
                                                    {format(new Date(ua.unlocked_at), 'MMM d, yyyy')}
                                                </span>
                                            </div>
                                        </motion.div>
                                    ))
                                ) : (
                                    <div className="col-span-full py-10 text-center text-gray-400 bg-gray-50/50 rounded-2xl border-2 border-dashed border-gray-100">
                                        <Trophy className="h-10 w-10 mx-auto mb-3 opacity-20" />
                                        <p className="font-medium">No achievements unlocked yet</p>
                                    </div>
                                )}
                            </div>
                        </Card>

                        {/* Recent Activity placeholder or Solved Problems */}
                        <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
                            <div className="flex items-center justify-between mb-8">
                                <h2 className="text-2xl font-bold text-gray-900">Solved Problems</h2>
                                <Button variant="ghost" size="sm" className="text-blue-600 font-bold hover:bg-blue-50 px-4 rounded-xl">
                                    View All
                                </Button>
                            </div>

                            <div className="space-y-4">
                                {profile.recent_problems && profile.recent_problems.length > 0 ? (
                                    profile.recent_problems.map((problem, index) => (
                                        <motion.div
                                            key={problem.id}
                                            initial={{ opacity: 0, x: -20 }}
                                            animate={{ opacity: 1, x: 0 }}
                                            transition={{ delay: index * 0.1 }}
                                            className="group flex items-center justify-between p-4 bg-gray-50/50 hover:bg-white hover:shadow-md rounded-2xl border border-transparent hover:border-gray-100 transition-all cursor-pointer"
                                            onClick={() => navigate(`/problems/${problem.slug}`)}
                                        >
                                            <div className="flex items-center gap-4">
                                                <div className={`p-2 rounded-xl text-white shadow-sm ${problem.difficulty === 'easy' ? 'bg-emerald-500' :
                                                    problem.difficulty === 'medium' ? 'bg-amber-500' : 'bg-rose-500'
                                                    }`}>
                                                    <Trophy className="h-4 w-4" />
                                                </div>
                                                <div>
                                                    <h4 className="font-bold text-gray-900 group-hover:text-blue-600 transition-colors uppercase tracking-tight">
                                                        {problem.title}
                                                    </h4>
                                                    <div className="flex items-center gap-3 text-xs font-semibold uppercase tracking-wider">
                                                        <span className={
                                                            problem.difficulty === 'easy' ? 'text-emerald-600' :
                                                                problem.difficulty === 'medium' ? 'text-amber-600' : 'text-rose-600'
                                                        }>
                                                            {problem.difficulty}
                                                        </span>
                                                        <span className="text-gray-400">•</span>
                                                        <span className="text-gray-500">{problem.acceptance_rate.toFixed(1)}% Accuracy</span>
                                                    </div>
                                                </div>
                                            </div>
                                            <ChevronLeft className="h-5 w-5 text-gray-300 group-hover:text-blue-500 rotate-180 transition-all transform group-hover:translate-x-1" />
                                        </motion.div>
                                    ))
                                ) : (
                                    <div className="py-12 text-center text-gray-400 bg-gray-50/50 rounded-3xl border-2 border-dashed border-gray-100">
                                        <Activity className="h-10 w-10 mx-auto mb-3 opacity-20" />
                                        <p className="font-medium">Recent solved problems will appear here</p>
                                    </div>
                                )}
                            </div>
                        </Card>
                    </div>
                </motion.div>
            </div>
        </div>
    )
}
