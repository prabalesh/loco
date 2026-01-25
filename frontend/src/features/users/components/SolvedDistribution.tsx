import { motion } from 'framer-motion'
import type { DifficultyStat } from '../types'

interface SolvedDistributionProps {
    distribution: DifficultyStat[]
    totalSolved: number
}

export const SolvedDistribution = ({ distribution = [], totalSolved }: SolvedDistributionProps) => {
    const getCount = (difficulty: string) =>
        distribution?.find(d => d.difficulty.toLowerCase() === difficulty.toLowerCase())?.count || 0

    const easy = getCount('easy')
    const medium = getCount('medium')
    const hard = getCount('hard')

    const stats = [
        { label: 'Easy', count: easy, color: 'bg-emerald-500', textColor: 'text-emerald-600', bg: 'bg-emerald-50' },
        { label: 'Medium', count: medium, color: 'bg-amber-500', textColor: 'text-amber-600', bg: 'bg-amber-50' },
        { label: 'Hard', count: hard, color: 'bg-rose-500', textColor: 'text-rose-600', bg: 'bg-rose-50' }
    ]

    return (
        <div className="space-y-6">
            <div className="flex items-end justify-between mb-2">
                <div>
                    <h3 className="text-sm font-bold text-gray-400 uppercase tracking-widest mb-1">Solved Problems</h3>
                    <div className="text-4xl font-black text-gray-900">{totalSolved}</div>
                </div>
                <div className="text-right">
                    <span className="text-xs font-bold text-gray-400 uppercase tracking-tighter">Distribution</span>
                </div>
            </div>

            <div className="space-y-4">
                {stats.map((stat, idx) => {
                    const percentage = totalSolved > 0 ? (stat.count / totalSolved) * 100 : 0
                    return (
                        <div key={stat.label} className="space-y-2">
                            <div className="flex justify-between items-center text-sm">
                                <span className={`font-bold ${stat.textColor} uppercase tracking-tight`}>{stat.label}</span>
                                <span className="font-mono font-bold text-gray-900">{stat.count}</span>
                            </div>
                            <div className={`h-3 w-full ${stat.bg} rounded-full overflow-hidden`}>
                                <motion.div
                                    initial={{ width: 0 }}
                                    animate={{ width: `${percentage}%` }}
                                    transition={{ duration: 1, delay: idx * 0.1 }}
                                    className={`h-full ${stat.color} rounded-full`}
                                />
                            </div>
                        </div>
                    )
                })}
            </div>
        </div>
    )
}
