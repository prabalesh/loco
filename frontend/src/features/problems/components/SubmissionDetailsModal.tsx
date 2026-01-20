import { motion, AnimatePresence } from 'framer-motion'
import { X, CheckCircle2, XCircle, Clock, MemoryStick, Code, Calendar } from 'lucide-react'
import type { Submission } from '../types'

interface SubmissionDetailsModalProps {
    submission: Submission | null
    isOpen: boolean
    onClose: () => void
}

export const SubmissionDetailsModal = ({ submission, isOpen, onClose }: SubmissionDetailsModalProps) => {
    if (!submission) return null

    return (
        <AnimatePresence>
            {isOpen && (
                <>
                    {/* Backdrop */}
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={onClose}
                        className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
                    />

                    {/* Modal */}
                    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95, y: 20 }}
                            animate={{ opacity: 1, scale: 1, y: 0 }}
                            exit={{ opacity: 0, scale: 0.95, y: 20 }}
                            transition={{ type: "spring", damping: 25, stiffness: 300 }}
                            className="bg-white rounded-3xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-hidden"
                            onClick={(e) => e.stopPropagation()}
                        >
                            {/* Header */}
                            <div className={`relative p-8 border-b border-gray-200 ${submission.status === 'Accepted'
                                ? 'bg-gradient-to-r from-emerald-50 to-emerald-100/50'
                                : 'bg-gradient-to-r from-rose-50 to-rose-100/50'
                                }`}>
                                <button
                                    onClick={onClose}
                                    className="absolute top-4 right-4 p-2 rounded-full hover:bg-white/50 transition-colors"
                                >
                                    <X className="h-5 w-5 text-gray-600" />
                                </button>

                                <div className="flex items-center gap-4 mb-4">
                                    <div className={`h-16 w-16 rounded-2xl flex items-center justify-center shadow-lg ${submission.status === 'Accepted'
                                        ? 'bg-emerald-500 text-white'
                                        : 'bg-rose-500 text-white'
                                        }`}>
                                        {submission.status === 'Accepted' ? (
                                            <CheckCircle2 className="h-8 w-8" />
                                        ) : (
                                            <XCircle className="h-8 w-8" />
                                        )}
                                    </div>
                                    <div>
                                        <h2 className={`text-3xl font-black mb-1 ${submission.status === 'Accepted' ? 'text-emerald-700' : 'text-rose-700'
                                            }`}>
                                            {submission.status}
                                        </h2>
                                        <p className="text-gray-600 font-medium">
                                            Submission #{submission.id}
                                        </p>
                                    </div>
                                </div>

                                {/* Stats Grid */}
                                <div className="grid grid-cols-4 gap-4 mt-6">
                                    <StatCard
                                        icon={<CheckCircle2 className="h-5 w-5" />}
                                        label="Test Cases"
                                        value={`${submission.passed_test_cases}/${submission.total_test_cases}`}
                                        color="blue"
                                    />
                                    <StatCard
                                        icon={<Clock className="h-5 w-5" />}
                                        label="Runtime"
                                        value={`${submission.runtime}ms`}
                                        color="purple"
                                    />
                                    <StatCard
                                        icon={<MemoryStick className="h-5 w-5" />}
                                        label="Memory"
                                        value={`${submission.memory}KB`}
                                        color="orange"
                                    />
                                    <StatCard
                                        icon={<Calendar className="h-5 w-5" />}
                                        label="Submitted"
                                        value={new Date(submission.created_at).toLocaleDateString(undefined, {
                                            month: 'short',
                                            day: 'numeric'
                                        })}
                                        color="green"
                                    />
                                </div>
                            </div>

                            {/* Body */}
                            <div className="p-8 overflow-y-auto max-h-[calc(90vh-300px)] custom-scrollbar">
                                {/* Language Info */}
                                <div className="mb-6">
                                    <h3 className="text-sm font-bold uppercase tracking-widest text-gray-500 mb-3 flex items-center gap-2">
                                        <Code className="h-4 w-4" />
                                        Language
                                    </h3>
                                    <div className="bg-gray-50 rounded-2xl p-4 border border-gray-200">
                                        <p className="font-mono text-sm text-gray-900">
                                            Language: <span className="font-bold">{submission.language.name}</span>
                                        </p>
                                    </div>
                                </div>

                                {/* Submitted Code */}
                                <div className="mb-6">
                                    <h3 className="text-sm font-bold uppercase tracking-widest text-gray-500 mb-3 flex items-center gap-2">
                                        <Code className="h-4 w-4" />
                                        Submitted Code
                                    </h3>
                                    <div className="bg-gradient-to-br from-gray-900 to-black rounded-2xl p-6 border border-gray-800 shadow-xl overflow-hidden">
                                        <div className="flex items-center gap-2 mb-4 pb-3 border-b border-gray-800">
                                            <div className="flex gap-1.5">
                                                <div className="w-3 h-3 rounded-full bg-rose-500 shadow-lg shadow-rose-500/50" />
                                                <div className="w-3 h-3 rounded-full bg-amber-500 shadow-lg shadow-amber-500/50" />
                                                <div className="w-3 h-3 rounded-full bg-emerald-500 shadow-lg shadow-emerald-500/50" />
                                            </div>
                                        </div>
                                        <pre className="text-sm font-mono text-gray-300 overflow-x-auto custom-scrollbar-dark leading-relaxed">
                                            {submission.function_code}
                                        </pre>
                                    </div>
                                </div>

                                {/* Error Message (if any) */}
                                {submission.error_message && (
                                    <div className="mb-6">
                                        <h3 className="text-sm font-bold uppercase tracking-widest text-gray-500 mb-3 flex items-center gap-2">
                                            <XCircle className="h-4 w-4" />
                                            Error Message
                                        </h3>
                                        <div className="bg-gradient-to-br from-rose-950 to-black rounded-2xl p-6 border border-rose-900 shadow-xl">
                                            <pre className="text-sm font-mono text-rose-300 whitespace-pre-wrap leading-relaxed">
                                                {submission.error_message}
                                            </pre>
                                        </div>
                                    </div>
                                )}

                                {/* Timestamp */}
                                <div className="text-center pt-4 border-t border-gray-200">
                                    <p className="text-xs text-gray-500">
                                        Submitted on {new Date(submission.created_at).toLocaleString(undefined, {
                                            weekday: 'long',
                                            year: 'numeric',
                                            month: 'long',
                                            day: 'numeric',
                                            hour: '2-digit',
                                            minute: '2-digit',
                                            second: '2-digit'
                                        })}
                                    </p>
                                </div>
                            </div>
                        </motion.div>
                    </div>
                </>
            )}
        </AnimatePresence>
    )
}

interface StatCardProps {
    icon: React.ReactNode
    label: string
    value: string
    color: 'blue' | 'purple' | 'orange' | 'green'
}

const StatCard = ({ icon, label, value, color }: StatCardProps) => {
    const colorClasses = {
        blue: 'bg-blue-50 text-blue-600 border-blue-200',
        purple: 'bg-purple-50 text-purple-600 border-purple-200',
        orange: 'bg-orange-50 text-orange-600 border-orange-200',
        green: 'bg-emerald-50 text-emerald-600 border-emerald-200',
    }

    return (
        <div className={`${colorClasses[color]} rounded-xl p-4 border shadow-sm`}>
            <div className="flex items-center justify-center mb-2">
                {icon}
            </div>
            <div className="text-center">
                <p className="text-xs font-medium opacity-70 mb-1">{label}</p>
                <p className="text-lg font-black">{value}</p>
            </div>
        </div>
    )
}
