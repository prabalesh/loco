import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { adminAnalyticsApi } from '../../../lib/api/admin'
import {
  Users,
  UserCheck,
  UserX,
  Users2,
  AlertTriangle,
} from 'lucide-react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import dayjs from 'dayjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

import { Skeleton, Grid, Box } from '@mui/material'

const DashboardSkeleton = () => (
  <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 p-8">
    {/* Header Skeleton */}
    <Box className="mb-10">
      <div className="flex items-start justify-between mb-6">
        <div className="flex items-center">
          <Skeleton variant="rounded" width={48} height={48} className="rounded-2xl mr-4" />
          <div className="space-y-1">
            <Skeleton variant="text" width={300} height={40} />
            <Skeleton variant="text" width={200} height={20} />
          </div>
        </div>
        <div className="text-right">
          <Skeleton variant="text" width={80} className="ml-auto" />
          <Skeleton variant="rounded" width={100} height={24} className="rounded-full ml-auto" />
        </div>
      </div>
    </Box>

    {/* Stats Grid Skeleton */}
    <Grid container spacing={3}>
      {[1, 2, 3, 4].map((i) => (
        <Grid size={{ xs: 12, md: 6, lg: 3 }} key={i}>
          <Box className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-44">
            <Skeleton variant="rounded" width={64} height={64} className="rounded-2xl mb-6" />
            <Skeleton variant="text" width="60%" height={36} />
            <Skeleton variant="text" width="40%" height={24} />
          </Box>
        </Grid>
      ))}
    </Grid>

    {/* Chart and System Health Skeleton */}
    <Grid container spacing={3} className="mt-10">
      <Grid size={{ xs: 12, lg: 8 }}>
        <Box className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-[450px]">
          <Skeleton variant="text" width={200} height={28} className="mb-2" />
          <Skeleton variant="text" width={250} height={20} className="mb-6" />
          <Skeleton variant="rectangular" width="100%" height={300} className="rounded-xl" />
        </Box>
      </Grid>
      <Grid size={{ xs: 12, lg: 4 }}>
        <Box className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-full">
          <Skeleton variant="text" width={150} height={28} className="mb-6" />
          <div className="space-y-6">
            <Skeleton variant="rounded" width="100%" height={80} className="rounded-2xl" />
            <Skeleton variant="rounded" width="100%" height={80} className="rounded-2xl" />
          </div>
        </Box>
      </Grid>
    </Grid>
  </div>
)

export const Dashboard = () => {
  const { data: analytics, isLoading } = useQuery({
    queryKey: ['admin-analytics'],
    queryFn: async () => {
      const response = await adminAnalyticsApi.getStats()
      return response.data.data
    },
  })

  if (isLoading) {
    return <DashboardSkeleton />
  }

  const totalUsers = analytics?.total_users || 0
  const activeUsers = analytics?.active_users || 0
  const inactiveUsers = analytics?.inactive_users || 0
  const verifiedUsers = analytics?.verified_users || 0

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 p-8">
      {/* Header */}
      <motion.header
        initial={{ opacity: 0, y: -30 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-10"
      >
        <div className="flex items-start justify-between mb-6">
          <div>
            <div className="flex items-center mb-2">
              <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg mr-4">
                <Users className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                  Admin Dashboard
                </h1>
                <p className="text-gray-600 font-medium">Real-time platform analytics</p>
              </div>
            </div>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-500 mb-1">Last updated</div>
            <div className="font-mono text-xs bg-white/50 px-3 py-1 rounded-full backdrop-blur-sm border border-white/60">
              {new Date().toLocaleTimeString()}
            </div>
          </div>
        </div>
      </motion.header>

      {/* Warning Alert for inactive workers */}
      {
        analytics?.active_workers === 0 && (
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            className="mb-8 p-4 bg-red-50 border border-red-200 rounded-2xl flex items-center justify-between shadow-sm"
          >
            <div className="flex items-center">
              <div className="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center mr-4">
                <AlertTriangle className="w-5 h-5 text-red-600" />
              </div>
              <div>
                <h3 className="text-red-900 font-bold">No Active Workers Detected</h3>
                <p className="text-red-700 text-sm">Submission processing is currently paused. Please check worker status.</p>
              </div>
            </div>
            <div className="hidden md:block">
              <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                Critical
              </span>
            </div>
          </motion.div>
        )
      }

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Total Users */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1 }}
          whileHover={{ y: -8 }}
          className="group"
        >
          <div className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl hover:shadow-3xl transition-all duration-500 hover:-translate-y-2 hover:bg-white/90">
            <div className="flex items-center justify-between mb-6">
              <div className="w-16 h-16 bg-gradient-to-br from-blue-100 to-blue-200 rounded-2xl flex items-center justify-center group-hover:from-blue-200 group-hover:to-blue-300 transition-all duration-300">
                <Users className="w-8 h-8 text-blue-600" />
              </div>
            </div>
            <div>
              <p className="text-3xl font-bold text-gray-900 mb-1">{totalUsers.toLocaleString()}</p>
              <p className="text-sm font-medium text-gray-600 uppercase tracking-wide">Total Users</p>
            </div>
          </div>
        </motion.div>

        {/* Active Users */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          whileHover={{ y: -8 }}
          className="group"
        >
          <div className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl hover:shadow-3xl transition-all duration-500 hover:-translate-y-2 hover:bg-white/90">
            <div className="flex items-center justify-between mb-6">
              <div className="w-16 h-16 bg-gradient-to-br from-emerald-100 to-emerald-200 rounded-2xl flex items-center justify-center group-hover:from-emerald-200 group-hover:to-emerald-300 transition-all duration-300">
                <UserCheck className="w-8 h-8 text-emerald-600" />
              </div>
            </div>
            <div>
              <p className="text-3xl font-bold text-gray-900 mb-1">{activeUsers.toLocaleString()}</p>
              <p className="text-sm font-medium text-gray-600 uppercase tracking-wide">Active Users</p>
            </div>
          </div>
        </motion.div>

        {/* Inactive Users */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
          whileHover={{ y: -8 }}
          className="group"
        >
          <div className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl hover:shadow-3xl transition-all duration-500 hover:-translate-y-2 hover:bg-white/90">
            <div className="flex items-center justify-between mb-6">
              <div className="w-16 h-16 bg-gradient-to-br from-orange-100 to-orange-200 rounded-2xl flex items-center justify-center group-hover:from-orange-200 group-hover:to-orange-300 transition-all duration-300">
                <UserX className="w-8 h-8 text-orange-600" />
              </div>
            </div>
            <div>
              <p className="text-3xl font-bold text-gray-900 mb-1">{inactiveUsers.toLocaleString()}</p>
              <p className="text-sm font-medium text-gray-600 uppercase tracking-wide">Inactive Users</p>
            </div>
          </div>
        </motion.div>

        {/* Verified Users */}
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          whileHover={{ y: -8 }}
          className="group"
        >
          <div className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl hover:shadow-3xl transition-all duration-500 hover:-translate-y-2 hover:bg-white/90">
            <div className="flex items-center justify-between mb-6">
              <div className="w-16 h-16 bg-gradient-to-br from-purple-100 to-purple-200 rounded-2xl flex items-center justify-center group-hover:from-purple-200 group-hover:to-purple-300 transition-all duration-300">
                <Users2 className="w-8 h-8 text-purple-600" />
              </div>
            </div>
            <div>
              <p className="text-3xl font-bold text-gray-900 mb-1">{verifiedUsers.toLocaleString()}</p>
              <p className="text-sm font-medium text-gray-600 uppercase tracking-wide">Verified Users</p>
            </div>
          </div>
        </motion.div>
      </div>

      {/* Analytics Chart */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-10">
        <div className="lg:col-span-2">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.6 }}
            className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl"
          >
            <div className="flex items-center justify-between mb-6">
              <div>
                <h3 className="text-xl font-bold text-gray-900">Submission Traffic</h3>
                <p className="text-sm text-gray-500">Daily submissions over last 7 days</p>
              </div>
            </div>
            <div className="h-[300px] w-full">
              <Line
                data={{
                  labels: analytics?.submission_history?.map(h => dayjs(h.date).format('MMM D')) || [],
                  datasets: [
                    {
                      label: 'Submissions',
                      data: analytics?.submission_history?.map(h => h.count) || [],
                      borderColor: '#4F46E5',
                      backgroundColor: 'rgba(79, 70, 229, 0.1)',
                      tension: 0.4,
                      fill: true,
                      pointBackgroundColor: '#fff',
                      pointBorderColor: '#4F46E5',
                      pointBorderWidth: 2,
                      pointRadius: 4,
                      pointHoverRadius: 6,
                    }
                  ]
                }}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: {
                    legend: {
                      display: false
                    },
                    tooltip: {
                      backgroundColor: 'rgba(17, 24, 39, 0.9)',
                      padding: 12,
                      cornerRadius: 8,
                      displayColors: false,
                    }
                  },
                  scales: {
                    y: {
                      beginAtZero: true,
                      grid: {
                        color: 'rgba(0, 0, 0, 0.05)',
                      },
                      ticks: {
                        stepSize: 1
                      }
                    },
                    x: {
                      grid: {
                        display: false
                      }
                    }
                  }
                }}
              />
            </div>
          </motion.div>
        </div>

        {/* System Health Card (Active Workers moved here) */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.7 }}
          className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-full"
        >
          <h3 className="text-xl font-bold text-gray-900 mb-6">System Status</h3>
          <div className="space-y-6">
            <div className={`flex items-center justify-between p-4 rounded-2xl border ${analytics?.active_workers === 0 ? 'bg-red-50 border-red-100' : 'bg-green-50 border-green-100'}`}>
              <div className="flex items-center">
                <div className={`w-10 h-10 rounded-xl flex items-center justify-center mr-3 ${analytics?.active_workers === 0 ? 'bg-red-100' : 'bg-green-100'}`}>
                  {analytics?.active_workers === 0 ? (
                    <AlertTriangle className="w-5 h-5 text-red-600" />
                  ) : (
                    <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse" />
                  )}
                </div>
                <div>
                  <p className="font-semibold text-gray-900">Worker Nodes</p>
                  <p className={`text-xs font-medium ${analytics?.active_workers === 0 ? 'text-red-600' : 'text-green-600'}`}>
                    {analytics?.active_workers === 0 ? 'Systems Critical' : 'Active & Healthy'}
                  </p>
                </div>
              </div>
              <span className={`text-2xl font-bold ${analytics?.active_workers === 0 ? 'text-red-700' : 'text-green-700'}`}>
                {analytics?.active_workers}
              </span>
            </div>

            <div className="flex items-center justify-between p-4 bg-blue-50 rounded-2xl border border-blue-100">
              <div className="flex items-center">
                <div className="w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center mr-3">
                  <Users2 className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <p className="font-semibold text-gray-900">Queue Status</p>
                  <p className="text-xs text-blue-600 font-medium">Jobs Pending</p>
                </div>
              </div>
              <span className="text-2xl font-bold text-blue-700">{analytics?.queue_size || 0}</span>
            </div>
          </div>
        </motion.div>
      </div>
    </div >
  )
}
