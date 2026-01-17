import { useState } from 'react'
import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { format } from 'date-fns'
import {
    FileCode2,
    Calendar,
    CheckCircle2,
    XCircle,
    Clock,
    AlertTriangle,
    ArrowRight,
    Code2
} from 'lucide-react'
import { submissionsApi } from '@/features/problems/api/submissions'
import { Card } from '@/shared/components/ui/Card'
import { ROUTES } from '@/shared/constants/routes'
import type { SubmissionStatus } from '@/features/problems/types'
import { Skeleton } from '@/shared/components/ui/Skeleton'

const SubmissionsSkeleton = () => (
    <div className="min-h-screen bg-gray-50 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex items-center justify-between mb-8">
                <div className="space-y-2">
                    <Skeleton className="h-9 w-64" />
                    <Skeleton className="h-5 w-96" />
                </div>
                <Skeleton className="w-12 h-12 rounded-2xl" />
            </div>

            <Card className="overflow-hidden border-0 shadow-lg bg-white">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                {[1, 2, 3, 4, 5].map((i) => (
                                    <th key={i} className="px-6 py-4"><Skeleton className="h-4 w-20" /></th>
                                ))}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                            {[1, 2, 3, 4, 5, 6].map((i) => (
                                <tr key={i}>
                                    {[1, 2, 3, 4, 5].map((j) => (
                                        <td key={j} className="px-6 py-4"><Skeleton className="h-5 w-full" /></td>
                                    ))}
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </Card>
        </div>
    </div>
)

export const SubmissionsPage = () => {
    const [page, setPage] = useState(1)
    const limit = 10

    const { data: submissionsData, isLoading } = useQuery({
        queryKey: ['user-submissions', page],
        queryFn: () => submissionsApi.listUserSubmissions(page, limit),
        placeholderData: keepPreviousData,
    })

    // Helper to render status badge (reused from Problem Detail but simplified)
    const renderStatus = (status: SubmissionStatus) => {
        switch (status) {
            case 'Accepted':
                return (
                    <span className="inline-flex items-center text-green-600 bg-green-50 px-2.5 py-0.5 rounded-full text-xs font-medium border border-green-100">
                        <CheckCircle2 className="w-3 h-3 mr-1" />
                        Accepted
                    </span>
                )
            case 'Pending':
            case 'Processing':
                return (
                    <span className="inline-flex items-center text-blue-600 bg-blue-50 px-2.5 py-0.5 rounded-full text-xs font-medium border border-blue-100">
                        <Clock className="w-3 h-3 mr-1" />
                        {status}
                    </span>
                )
            case 'Wrong Answer':
                return (
                    <span className="inline-flex items-center text-red-600 bg-red-50 px-2.5 py-0.5 rounded-full text-xs font-medium border border-red-100">
                        <XCircle className="w-3 h-3 mr-1" />
                        Wrong Answer
                    </span>
                )
            default:
                return (
                    <span className="inline-flex items-center text-orange-600 bg-orange-50 px-2.5 py-0.5 rounded-full text-xs font-medium border border-orange-100">
                        <AlertTriangle className="w-3 h-3 mr-1" />
                        {status}
                    </span>
                )
        }
    }

    if (isLoading) {
        return <SubmissionsSkeleton />
    }

    const submissions = submissionsData?.data.data.data || []
    const total = submissionsData?.data.total || 0
    const totalPages = Math.ceil(total / limit)

    return (
        <div className="min-h-screen bg-gray-50 py-12">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.5 }}
                >
                    <div className="flex items-center justify-between mb-8">
                        <div>
                            <h1 className="text-3xl font-bold text-gray-900">My Submissions</h1>
                            <p className="text-gray-600 mt-1">Track your coding journey and submission history</p>
                        </div>
                        <div className="w-12 h-12 bg-blue-100 rounded-2xl flex items-center justify-center">
                            <FileCode2 className="w-6 h-6 text-blue-600" />
                        </div>
                    </div>

                    <Card className="overflow-hidden border-0 shadow-lg bg-white/80 backdrop-blur-sm">
                        <div className="overflow-x-auto">
                            {submissions.length === 0 ? (
                                <div className="text-center py-20">
                                    <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                                        <Code2 className="w-8 h-8 text-gray-400" />
                                    </div>
                                    <h3 className="text-lg font-medium text-gray-900 mb-2">No submissions yet</h3>
                                    <p className="text-gray-500 mb-6">Start solving problems to see your history here!</p>
                                    <Link to={ROUTES.PROBLEMS} className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                                        Browse Problems
                                        <ArrowRight className="w-4 h-4 ml-2" />
                                    </Link>
                                </div>
                            ) : (
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50/50">
                                        <tr>
                                            <th scope="col" className="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Problem
                                            </th>
                                            <th scope="col" className="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Status
                                            </th>
                                            <th scope="col" className="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Language
                                            </th>
                                            <th scope="col" className="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Runtime
                                            </th>
                                            <th scope="col" className="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                                Date
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {submissions.map((submission) => (
                                            <tr key={submission.id} className="hover:bg-gray-50/80 transition-colors">
                                                <td className="px-6 py-4 whitespace-nowrap">
                                                    <Link to={`/problems/${submission.problem?.slug}`} className="text-sm font-medium text-blue-600 hover:text-blue-800 hover:underline">
                                                        {submission.problem?.title || `Problem #${submission.problem_id}`}
                                                    </Link>
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap">
                                                    {renderStatus(submission.status)}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                    {submission.language.name}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-mono">
                                                    {submission.runtime > 0 ? `${submission.runtime}ms` : '-'}
                                                </td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                    <div className="flex items-center">
                                                        <Calendar className="w-3 h-3 mr-1.5 text-gray-400" />
                                                        {format(new Date(submission.created_at), 'MMM d, yyyy HH:mm')}
                                                    </div>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </div>

                        {/* Pagination */}
                        {totalPages > 1 && (
                            <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
                                <button
                                    onClick={() => setPage(p => Math.max(1, p - 1))}
                                    disabled={page === 1}
                                    className="px-3 py-1 border rounded text-sm disabled:opacity-50"
                                >
                                    Previous
                                </button>
                                <span className="text-sm text-gray-500">
                                    Page {page} of {totalPages}
                                </span>
                                <button
                                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                                    disabled={page === totalPages}
                                    className="px-3 py-1 border rounded text-sm disabled:opacity-50"
                                >
                                    Next
                                </button>
                            </div>
                        )}
                    </Card>
                </motion.div>
            </div>
        </div>
    )
}
