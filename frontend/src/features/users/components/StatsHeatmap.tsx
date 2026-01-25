import { useMemo } from 'react'
import { motion } from 'framer-motion'
import type { HeatmapEntry } from '../types'
import { format, subDays, startOfToday, eachDayOfInterval } from 'date-fns'

interface StatsHeatmapProps {
    data: HeatmapEntry[]
}


export const StatsHeatmap = ({ data = [] }: StatsHeatmapProps) => {
    const days = useMemo(() => {
        const today = startOfToday()
        const startDate = subDays(today, 364) // Last 365 days
        const interval = eachDayOfInterval({ start: startDate, end: today })

        return interval.map(date => {
            const dateStr = format(date, 'yyyy-MM-dd')
            const entry = data?.find(d => d.date === dateStr)
            return {
                date: dateStr,
                count: entry?.count || 0,
                weekday: date.getDay()
            }
        })
    }, [data])

    const getLevel = (count: number) => {
        if (count === 0) return 'bg-gray-100'
        if (count < 3) return 'bg-blue-200'
        if (count < 6) return 'bg-blue-400'
        if (count < 10) return 'bg-blue-600'
        return 'bg-blue-800'
    }

    return (
        <div className="w-full overflow-x-auto pb-4">
            <div className="flex items-center gap-3 mb-6">
                <h3 className="text-lg font-bold text-gray-900">Activity Heatmap</h3>
                <div className="flex gap-1 items-center ml-auto">
                    <span className="text-[10px] font-bold text-gray-400 uppercase mr-1">Less</span>
                    {[0, 2, 5, 8, 12].map(c => (
                        <div key={c} className={`w-3 h-3 rounded-sm ${getLevel(c)}`} />
                    ))}
                    <span className="text-[10px] font-bold text-gray-400 uppercase ml-1">More</span>
                </div>
            </div>

            <div className="grid grid-flow-col grid-rows-7 gap-1.5 min-w-fit">
                {days.map((day, i) => (
                    <motion.div
                        key={day.date}
                        initial={{ opacity: 0, scale: 0 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: (i % 50) * 0.005 }}
                        className={`w-3.5 h-3.5 rounded-sm ${getLevel(day.count)} transition-all hover:ring-2 hover:ring-blue-300 relative group cursor-pointer`}
                    >
                        <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-1.5 bg-gray-900 text-white text-[10px] font-bold rounded-lg opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10 pointer-events-none shadow-xl">
                            {day.count} submissions on {format(new Date(day.date), 'MMM d, yyyy')}
                        </div>
                    </motion.div>
                ))}
            </div>
            <div className="flex justify-between mt-4 text-[10px] font-bold text-gray-400 uppercase tracking-widest">
                <span>{format(subDays(new Date(), 364), 'MMM yyyy')}</span>
                <span>{format(new Date(), 'MMM yyyy')}</span>
            </div>
        </div>
    )
}
