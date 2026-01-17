import { useQuery } from '@tanstack/react-query'
import { AlertTriangle, Clock, XCircle } from 'lucide-react'
import { queueApi } from '../api/queue'

interface QueueStatusBannerProps {
    show: boolean
}

export const QueueStatusBanner = ({ show }: QueueStatusBannerProps) => {
    const { data: queueStatus } = useQuery({
        queryKey: ['queue-status'],
        queryFn: () => queueApi.getStatus().then(res => res.data.data),
        enabled: show,
        refetchInterval: show ? 5000 : false, // Refresh every 5 seconds when visible
    })

    if (!show || !queueStatus) return null

    // Don't show banner if queue is healthy
    if (queueStatus.health_status === 'healthy') return null

    const isCritical = queueStatus.health_status === 'critical'
    const isWarning = queueStatus.health_status === 'warning'

    const bgColor = isCritical ? 'bg-red-50 border-red-200' : 'bg-yellow-50 border-yellow-200'
    const textColor = isCritical ? 'text-red-800' : 'text-yellow-800'
    const iconColor = isCritical ? 'text-red-600' : 'text-yellow-600'

    const Icon = isCritical ? XCircle : isWarning ? AlertTriangle : Clock

    const formatTime = (seconds: number) => {
        if (seconds < 60) return `${seconds}s`
        const minutes = Math.floor(seconds / 60)
        const remainingSeconds = seconds % 60
        return remainingSeconds > 0 ? `${minutes}m ${remainingSeconds}s` : `${minutes}m`
    }

    return (
        <div className={`border-l-4 p-4 mb-4 rounded-r-lg ${bgColor}`}>
            <div className="flex items-start">
                <Icon className={`h-5 w-5 ${iconColor} mt-0.5 mr-3 flex-shrink-0`} />
                <div className="flex-1">
                    <h3 className={`text-sm font-semibold ${textColor} mb-1`}>
                        {isCritical ? 'Queue System Issue' : 'Queue Delay'}
                    </h3>
                    <div className={`text-sm ${textColor} space-y-1`}>
                        {queueStatus.warning_message && (
                            <p>{queueStatus.warning_message}</p>
                        )}
                        {queueStatus.active_workers === 0 && (
                            <p className="font-medium">⚠️ No workers are currently processing submissions.</p>
                        )}
                        {queueStatus.oldest_pending_age_seconds > 0 && (
                            <p>
                                Oldest submission has been waiting for{' '}
                                <span className="font-semibold">
                                    {formatTime(queueStatus.oldest_pending_age_seconds)}
                                </span>
                            </p>
                        )}
                        {queueStatus.queue_size > 0 && (
                            <p>
                                <span className="font-semibold">{queueStatus.queue_size}</span> submission(s) in queue
                            </p>
                        )}
                        {queueStatus.estimated_wait_time_seconds && queueStatus.estimated_wait_time_seconds > 0 && (
                            <p className="text-xs mt-2 opacity-90">
                                Estimated wait time: ~{formatTime(queueStatus.estimated_wait_time_seconds)}
                            </p>
                        )}
                    </div>
                </div>
            </div>
        </div>
    )
}
