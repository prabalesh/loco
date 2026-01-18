import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { leaderboardApi } from '../api/leaderboard'
import { Card } from '@/shared/components/ui/Card'
import { Skeleton } from '@/shared/components/ui/Skeleton'
import { Trophy, Medal, User, ArrowRight, TrendingUp, Zap } from 'lucide-react'
import { Button } from '@/shared/components/ui/Button'
import { motion } from 'framer-motion'

export const LeaderboardPage = () => {
    const navigate = useNavigate()
    const { data: leaderboard, isLoading } = useQuery({
        queryKey: ['leaderboard'],
        queryFn: () => leaderboardApi.getLeaderboard().then(res => (res.data as any).data),
    })

    if (isLoading) {
        return (
            <div className="max-w-5xl mx-auto px-4 py-12">
                <Skeleton className="h-12 w-64 mb-8" />
                <div className="space-y-4">
                    {[...Array(5)].map((_, i) => (
                        <Skeleton key={i} className="h-20 w-full rounded-2xl" />
                    ))}
                </div>
            </div>
        )
    }

    const getRankIcon = (rank: number) => {
        switch (rank) {
            case 1: return <Trophy className="h-6 w-6 text-yellow-500" />
            case 2: return <Medal className="h-6 w-6 text-gray-400" />
            case 3: return <Medal className="h-6 w-6 text-amber-600" />
            default: return <span className="text-lg font-bold text-gray-400">#{rank}</span>
        }
    }

    return (
        <div className="min-h-screen bg-gray-50/50 py-12">
            <div className="max-w-5xl mx-auto px-4">
                <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12">
                    <div>
                        <motion.div
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            className="flex items-center gap-3 mb-4"
                        >
                            <div className="p-2 bg-blue-600 rounded-xl shadow-lg shadow-blue-200">
                                <Trophy className="h-5 w-5 text-white" />
                            </div>
                            <span className="text-sm font-bold text-blue-600 uppercase tracking-wider">Hall of Fame</span>
                        </motion.div>
                        <h1 className="text-4xl md:text-5xl font-black text-gray-900 tracking-tight">
                            Global <span className="text-blue-600">Leaderboard</span>
                        </h1>
                    </div>
                </div>

                <div className="grid grid-cols-1 gap-6">
                    {leaderboard?.map((entry: any, index: number) => (
                        <motion.div
                            key={entry.user_id}
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: index * 0.05 }}
                        >
                            <Card
                                className={`p-6 border-0 shadow-xl shadow-gray-200/50 rounded-[2rem] bg-white group hover:ring-2 hover:ring-blue-500/20 transition-all cursor-pointer ${entry.rank <= 3 ? 'ring-1 ring-blue-100' : ''
                                    }`}
                                onClick={() => navigate(`/users/${entry.username}`)}
                            >
                                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-6">
                                    <div className="flex items-center gap-6">
                                        <div className="w-12 h-12 flex items-center justify-center bg-gray-50 rounded-2xl group-hover:bg-blue-50 transition-colors">
                                            {getRankIcon(entry.rank)}
                                        </div>

                                        <div className="flex items-center gap-4">
                                            <div className="h-14 w-14 bg-gradient-to-br from-blue-100 to-indigo-100 rounded-2xl flex items-center justify-center border-2 border-white shadow-sm">
                                                <User className="h-7 w-7 text-blue-600" />
                                            </div>
                                            <div>
                                                <h3 className="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors">
                                                    @{entry.username}
                                                </h3>
                                                <div className="flex items-center gap-3 text-sm text-gray-500 font-medium">
                                                    <span className="flex items-center gap-1">
                                                        <TrendingUp className="h-3.5 w-3.5 text-emerald-500" />
                                                        Rank {entry.rank}
                                                    </span>
                                                    <span>•</span>
                                                    <span>Player</span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-4 sm:gap-12 ml-auto sm:ml-0">
                                        <div className="text-center">
                                            <div className="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1 pointer-events-none">Solved</div>
                                            <div className="text-2xl font-black text-gray-900 group-hover:text-blue-600 transition-colors">
                                                {entry.problems_solved}
                                            </div>
                                        </div>

                                        <div className="text-center">
                                            <div className="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1 pointer-events-none">Accuracy</div>
                                            <div className="text-2xl font-black text-gray-900">
                                                {entry.acceptance_rate.toFixed(1)}<span className="text-sm font-bold text-gray-400 ml-0.5">%</span>
                                            </div>
                                        </div>

                                        <div className="text-center">
                                            <div className="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1 pointer-events-none flex items-center justify-center gap-1">
                                                <Zap className="h-3 w-3" />
                                                Level
                                            </div>
                                            <div className="text-2xl font-black text-purple-600">
                                                {entry.level || 1}
                                            </div>
                                        </div>

                                        <div className="hidden md:block text-center border-l border-gray-100 pl-12">
                                            <div className="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1 flex items-center justify-center gap-1">
                                                <Zap className="h-3 w-3 text-amber-500" />
                                                XP
                                            </div>
                                            <div className="text-2xl font-black text-amber-600">{(entry.xp || 0).toLocaleString()}</div>
                                        </div>

                                        <div className="sm:ml-4">
                                            <Button variant="ghost" size="sm" className="rounded-xl group-hover:bg-blue-600 group-hover:text-white">
                                                <ArrowRight className="h-5 w-5" />
                                            </Button>
                                        </div>
                                    </div>
                                </div>
                            </Card>
                        </motion.div>
                    ))}
                </div>
            </div>
        </div>
    )
}
