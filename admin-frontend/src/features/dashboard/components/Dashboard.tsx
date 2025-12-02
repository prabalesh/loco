import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { adminAnalyticsApi } from '../../../api/adminApi'
import type { AdminAnalytics } from '../../../types'
import {
  Users,
  UserCheck,
  UserX,
  Users2,
} from 'lucide-react'

interface EnhancedAnalytics extends AdminAnalytics {
  total_admins?: number
  recent_registrations?: number
  avg_daily_users?: number
}

export const Dashboard = () => {
  const { data: analytics, isLoading, error } = useQuery({
    queryKey: ['admin-analytics'],
    queryFn: async () => {
      const response = await adminAnalyticsApi.getAnalytics()
      return response.data as EnhancedAnalytics
    },
  })

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-blue-50">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          className="w-16 h-16 border-4 border-blue-200 border-t-blue-500 rounded-full"
        />
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-blue-50">
        <div className="text-center p-12 max-w-md">
          <div className="w-24 h-24 bg-red-100 rounded-2xl flex items-center justify-center mx-auto mb-6">
            <Users className="w-12 h-12 text-red-500" />
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Loading Failed</h2>
          <p className="text-gray-600 mb-6">Unable to fetch analytics data</p>
          <button
            onClick={() => window.location.reload()}
            className="px-6 py-2 bg-blue-600 text-white rounded-xl font-medium hover:bg-blue-700 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
          >
            Retry
          </button>
        </div>
      </div>
    )
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
    </div>
  )
}
