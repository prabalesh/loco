import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { adminAnalyticsApi } from '../../../lib/api/admin'
import {
  Users,
  UserCheck,
  UserX,
  Users2,
  AlertTriangle,
  Flame,
  Globe,
  ArrowUpRight,
} from 'lucide-react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'
import { Line, Doughnut } from 'react-chartjs-2'
import dayjs from 'dayjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

import { Skeleton, Grid, Box, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Chip } from '@mui/material'

const DashboardSkeleton = () => (
  <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 p-8">
    <Box className="mb-10">
      <div className="flex items-start justify-between mb-6">
        <div className="flex items-center">
          <Skeleton variant="rounded" width={48} height={48} className="rounded-2xl mr-4" />
          <div className="space-y-1">
            <Skeleton variant="text" width={300} height={40} />
            <Skeleton variant="text" width={200} height={20} />
          </div>
        </div>
      </div>
    </Box>

    <Grid container spacing={3}>
      {[1, 2, 3, 4].map((i) => (
        <Grid item xs={12} md={6} lg={3} key={i}>
          <Box className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-44">
            <Skeleton variant="rounded" width={64} height={64} className="rounded-2xl mb-6" />
            <Skeleton variant="text" width="60%" />
          </Box>
        </Grid>
      ))}
    </Grid>

    <Grid container spacing={3} className="mt-10">
      <Grid item xs={12} lg={8}>
        <Box className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-[450px]">
          <Skeleton variant="rectangular" width="100%" height="100%" className="rounded-xl" />
        </Box>
      </Grid>
      <Grid item xs={12} lg={4}>
        <Box className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-[450px]">
          <Skeleton variant="rectangular" width="100%" height="100%" className="rounded-xl" />
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

  const languageChartData = {
    labels: analytics?.language_stats?.map(s => s.language_name) || [],
    datasets: [
      {
        data: analytics?.language_stats?.map(s => s.count) || [],
        backgroundColor: [
          '#4F46E5',
          '#10B981',
          '#F59E0B',
          '#EF4444',
          '#8B5CF6',
          '#EC4899',
        ],
        borderWidth: 0,
        hoverOffset: 15,
      }
    ]
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 p-8">
      <motion.header
        initial={{ opacity: 0, y: -30 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-10"
      >
        <div className="flex items-start justify-between mb-6">
          <div className="flex items-center">
            <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-2xl flex items-center justify-center shadow-lg mr-4">
              <Users className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                Platform Insights
              </h1>
              <p className="text-gray-600 font-medium">Real-time Admin Analytics</p>
            </div>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-500">Last updated</div>
            <div className="font-mono text-xs bg-white/50 px-3 py-1 rounded-full border border-white/60">
              {new Date().toLocaleTimeString()}
            </div>
          </div>
        </div>
      </motion.header>

      {analytics?.active_workers === 0 && (
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className="mb-8 p-4 bg-red-50 border border-red-200 rounded-2xl flex items-center gap-4 shadow-sm"
        >
          <div className="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
            <AlertTriangle className="w-5 h-5 text-red-600" />
          </div>
          <div>
            <h3 className="text-red-900 font-bold">Offline Workers</h3>
            <p className="text-red-700 text-sm">Submission execution is currently halted.</p>
          </div>
        </motion.div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {[
          { icon: Users, label: 'Total Users', value: totalUsers, color: 'text-blue-600', bg: 'bg-blue-100' },
          { icon: UserCheck, label: 'Active', value: activeUsers, color: 'text-emerald-600', bg: 'bg-emerald-100' },
          { icon: UserX, label: 'Inactive', value: inactiveUsers, color: 'text-orange-600', bg: 'bg-orange-100' },
          { icon: Users2, label: 'Verified', value: verifiedUsers, color: 'text-purple-600', bg: 'bg-purple-100' }
        ].map((stat, idx) => (
          <motion.div
            key={stat.label}
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: idx * 0.1 }}
            className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl hover:bg-white/90 transition-all duration-300"
          >
            <div className={`w-14 h-14 ${stat.bg} rounded-2xl flex items-center justify-center mb-6`}>
              <stat.icon className={`w-7 h-7 ${stat.color}`} />
            </div>
            <p className="text-3xl font-bold text-gray-900">{stat.value.toLocaleString()}</p>
            <p className="text-xs font-bold text-gray-500 uppercase tracking-widest">{stat.label}</p>
          </motion.div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-10">
        <div className="lg:col-span-2">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-full"
          >
            <h3 className="text-xl font-bold text-gray-900 mb-6">Submission Trends</h3>
            <div className="h-[300px]">
              <Line
                data={{
                  labels: analytics?.submission_history?.map(h => dayjs(h.date).format('MMM D')) || [],
                  datasets: [{
                    label: 'Submissions',
                    data: analytics?.submission_history?.map(h => h.count) || [],
                    borderColor: '#4F46E5',
                    backgroundColor: 'rgba(79, 70, 229, 0.1)',
                    tension: 0.4,
                    fill: true,
                  }]
                }}
                options={{ responsive: true, maintainAspectRatio: false, plugins: { legend: { display: false } } }}
              />
            </div>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-full flex flex-col"
        >
          <div className="flex items-center gap-3 mb-8">
            <div className="p-2 bg-emerald-100 rounded-xl">
              <Globe className="w-5 h-5 text-emerald-600" />
            </div>
            <h3 className="text-xl font-bold text-gray-900">Languages</h3>
          </div>
          <div className="flex-grow flex items-center justify-center">
            <div className="h-[250px] w-full">
              <Doughnut
                data={languageChartData}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  cutout: '70%',
                  plugins: { legend: { position: 'bottom' } }
                }}
              />
            </div>
          </div>
        </motion.div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-10">
        <div className="lg:col-span-2">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-full"
          >
            <div className="flex items-center justify-between mb-8">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-orange-100 rounded-xl">
                  <Flame className="w-5 h-5 text-orange-600" />
                </div>
                <h3 className="text-xl font-bold text-gray-900">Trending Problems</h3>
              </div>
              <Chip label="Last 7 Days" size="small" className="bg-orange-100 text-orange-700 font-bold" />
            </div>

            <TableContainer>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell sx={{ fontWeight: 800, color: '#64748b' }}>Problem</TableCell>
                    <TableCell align="right" sx={{ fontWeight: 800, color: '#64748b' }}>Submissions</TableCell>
                    <TableCell align="right"></TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {analytics?.trending_problems?.map((p) => (
                    <TableRow key={p.id} hover>
                      <TableCell sx={{ py: 2 }}>
                        <div className="font-bold text-gray-900">{p.title}</div>
                        <div className="text-xs text-gray-400 font-mono">{p.slug}</div>
                      </TableCell>
                      <TableCell align="right">
                        <span className="px-3 py-1 rounded-full text-xs font-black bg-indigo-50 text-indigo-700">
                          {p.submission_count}
                        </span>
                      </TableCell>
                      <TableCell align="right">
                        <ArrowUpRight className="w-4 h-4 text-gray-300" />
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          className="bg-white/70 backdrop-blur-xl rounded-3xl p-8 border border-white/50 shadow-2xl h-full"
        >
          <h3 className="text-xl font-bold text-gray-900 mb-8">Operational Status</h3>
          <div className="space-y-6">
            <div className="p-5 bg-gray-50 rounded-2xl border border-gray-100">
              <div className="flex justify-between items-center mb-1">
                <span className="text-xs font-bold text-gray-400 uppercase tracking-widest">Active Workers</span>
                <span className="text-2xl font-black text-emerald-600">{analytics?.active_workers}</span>
              </div>
              <div className="w-full bg-gray-200 h-1.5 rounded-full overflow-hidden">
                <div className="bg-emerald-500 h-full" style={{ width: `${Math.min(analytics?.active_workers || 0 * 10, 100)}%` }}></div>
              </div>
            </div>

            <div className="p-5 bg-blue-50/50 rounded-2xl border border-blue-100">
              <div className="flex justify-between items-center mb-1">
                <span className="text-xs font-bold text-gray-400 uppercase tracking-widest">Queue Status</span>
                <span className="text-2xl font-black text-blue-600">{analytics?.queue_size}</span>
              </div>
              <p className="text-xs text-blue-400 font-medium">Pending submissions in queue</p>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  )
}
