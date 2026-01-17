import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export interface QueueStatus {
    queue_size: number
    active_workers: number
    oldest_pending_age_seconds: number
    health_status: 'healthy' | 'warning' | 'critical'
    warning_message?: string
    pending_count: number
    estimated_wait_time_seconds?: number
}

export const queueApi = {
    getStatus: () => axios.get<{ data: QueueStatus }>(`${API_BASE_URL}/queue/status`)
}
