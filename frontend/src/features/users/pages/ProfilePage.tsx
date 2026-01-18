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
                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-6 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-emerald-50 hover:ring-emerald-200 transition-all">
                                    <div className="flex items-center gap-4 mb-4">
                                        <div className="p-3 bg-emerald-50 rounded-2xl group-hover:bg-emerald-500 group-hover:text-white transition-colors text-emerald-600">
                                            <Trophy className="h-5 w-5" />
                                        </div>
                                        <div className="text-sm font-bold text-gray-500 uppercase tracking-tight">Solved</div>
                                    </div>
                                    <div className="text-3xl font-black text-gray-900">{profile.stats.problems_solved}</div>
                                </Card>
                            </motion.div>

                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-6 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-blue-50 hover:ring-blue-200 transition-all">
                                    <div className="flex items-center gap-4 mb-4">
                                        <div className="p-3 bg-blue-50 rounded-2xl group-hover:bg-blue-500 group-hover:text-white transition-colors text-blue-600">
                                            <Activity className="h-5 w-5" />
                                        </div>
                                        <div className="text-sm font-bold text-gray-500 uppercase tracking-tight">Accuracy</div>
                                    </div>
                                    <div className="text-3xl font-black text-gray-900">{profile.stats.acceptance_rate.toFixed(1)}%</div>
                                </Card>
                            </motion.div>

                            <motion.div whileHover={{ y: -5 }}>
                                <Card className="p-6 border-0 shadow-lg shadow-gray-200/40 rounded-3xl bg-white group ring-1 ring-indigo-50 hover:ring-indigo-200 transition-all">
                                    <div className="flex items-center gap-4 mb-4">
                                        <div className="p-3 bg-indigo-50 rounded-2xl group-hover:bg-indigo-500 group-hover:text-white transition-colors text-indigo-600">
                                            <Activity className="h-5 w-5" />
                                        </div>
                                        <div className="text-sm font-bold text-gray-500 uppercase tracking-tight">Submissions</div>
                                    </div>
                                    <div className="text-3xl font-black text-gray-900">{profile.stats.total_submissions}</div>
                                </Card>
                            </motion.div>
                        </div>

                        {/* Recent Activity placeholder or Solved Problems */}
                        <Card className="p-8 border-0 shadow-xl shadow-gray-200/30 rounded-[2.5rem] bg-white">
                            <div className="flex items-center justify-between mb-8">
                                <h2 className="text-2xl font-bold text-gray-900">Solved Problems</h2>
                                <Button variant="ghost" size="sm" className="text-blue-600 font-bold hover:bg-blue-50 px-4 rounded-xl">
                                    View All
                                </Button>
                            </div>

                            <div className="space-y-4">
                                <div className="py-12 text-center text-gray-400 bg-gray-50/50 rounded-3xl border-2 border-dashed border-gray-100">
                                    <Activity className="h-10 w-10 mx-auto mb-3 opacity-20" />
                                    <p className="font-medium">Recent solved problems will appear here</p>
                                </div>
                            </div>
                        </Card>
                    </div>
                </motion.div>
            </div>
        </div>
    )
}
