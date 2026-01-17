import { motion } from 'framer-motion'
import { History, CheckCircle2, XCircle } from 'lucide-react'
import type { Submission } from '../types'

interface SubmissionsTabProps {
    submissionsHistory?: Submission[]
    problemId: number
    onViewSubmission: (submission: Submission) => void
}

export const SubmissionsTab = ({ submissionsHistory, problemId, onViewSubmission }: SubmissionsTabProps) => {
    const filteredSubmissions = submissionsHistory?.filter((s: Submission) => s.problem_id === problemId) || []

    if (filteredSubmissions.length === 0) {
        return (
            <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="flex flex-col items-center justify-center py-32 text-gray-400"
            >
                <div className="relative">
                    <div className="absolute inset-0 bg-gray-200 rounded-full blur-3xl opacity-30" />
                    <History className="h-20 w-20 mb-6 opacity-20 relative" />
                </div>
                <p className="text-lg font-semibold opacity-50">No submissions yet</p>
                <p className="text-sm opacity-40 mt-2">Your previous attempts will appear here.</p>
            </motion.div>
        )
    }

    return (
        <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="space-y-4"
        >
            {filteredSubmissions.map((sub: Submission, index: number) => (
                <motion.div
                    key={sub.id}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.05 }}
                    onClick={() => onViewSubmission(sub)}
                    className="group flex items-center justify-between p-6 rounded-3xl border border-gray-200 hover:border-blue-300 hover:bg-gradient-to-r hover:from-blue-50/50 hover:to-transparent transition-all duration-300 shadow-sm hover:shadow-md cursor-pointer"
                >
                    <div className="flex items-center gap-5">
                        <div className={`h-12 w-12 rounded-2xl flex items-center justify-center shadow-sm ${sub.status === 'Accepted'
                            ? 'bg-gradient-to-br from-emerald-50 to-emerald-100 text-emerald-600 border border-emerald-200'
                            : 'bg-gradient-to-br from-rose-50 to-rose-100 text-rose-600 border border-rose-200'
                            }`}>
                            {sub.status === 'Accepted' ? <CheckCircle2 className="h-6 w-6" /> : <XCircle className="h-6 w-6" />}
                        </div>
                        <div>
                            <div className={`font-black text-lg mb-1 ${sub.status === 'Accepted' ? 'text-emerald-700' : 'text-rose-700'
                                }`}>
                                {sub.status}
                            </div>
                            <div className="text-xs text-gray-500 font-medium">
                                {new Date(sub.created_at).toLocaleString(undefined, {
                                    month: 'short',
                                    day: 'numeric',
                                    year: 'numeric',
                                    hour: '2-digit',
                                    minute: '2-digit'
                                })}
                            </div>
                        </div>
                    </div>
                    <div className="text-right">
                        <div className="text-base font-black text-gray-900 mb-1">
                            {sub.passed_test_cases} / {sub.total_test_cases} <span className="text-xs text-gray-500 font-normal">Passed</span>
                        </div>
                        <div className="text-xs text-blue-600 font-bold opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                            View Details →
                        </div>
                    </div>
                </motion.div>
            ))}
        </motion.div>
    )
}