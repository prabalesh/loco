import { motion } from 'framer-motion'
import { Play, CheckCircle2, XCircle } from 'lucide-react'
import type { Submission } from '../types'
import type { RunCodeResult } from '../api/submissions'
import { QueueStatusBanner } from './QueueStatusBanner'

interface ResultTabProps {
    submissionResult?: Submission
    pollingId: number | null
    runResult?: RunCodeResult | null
    isRunning?: boolean
}

export const ResultTab = ({ submissionResult, pollingId, runResult, isRunning }: ResultTabProps) => {
    // Show loading state for either submission polling or run execution
    if ((pollingId || isRunning) && !submissionResult && !runResult) {
        return (
            <div>
                <QueueStatusBanner show={!!pollingId} />
                <motion.div
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="flex flex-col items-center justify-center py-32 text-gray-400"
                >
                    <div className="relative">
                        <div className="absolute inset-0 bg-blue-100 rounded-full blur-3xl opacity-20 animate-pulse" />
                        <Play className="h-20 w-20 mb-6 opacity-20 relative" />
                    </div>
                    <p className="text-lg font-semibold opacity-50">Ready to test your solution?</p>
                    <p className="text-sm opacity-40 mt-2">Submit your code to see the results here.</p>
                </motion.div>
            </div>
        )
    }

    // Show run result if available (takes priority over submission result for display)
    const displayResult = runResult || submissionResult
    const isRunResult = !!runResult && !submissionResult

    if (displayResult && !pollingId && !isRunning) {
        return (
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4 }}
                className="space-y-6"
            >
                <div className={`p-8 rounded-3xl border-2 shadow-xl flex items-center justify-between relative overflow-hidden ${displayResult.status === 'Accepted'
                    ? 'bg-gradient-to-br from-emerald-50 to-emerald-100/50 border-emerald-200'
                    : 'bg-gradient-to-br from-rose-50 to-rose-100/50 border-rose-200'
                    }`}>
                    <div className="absolute top-0 right-0 w-64 h-64 bg-white rounded-full blur-3xl opacity-50" />
                    <div className="relative z-10">
                        <h2 className={`text-3xl font-black flex items-center gap-3 mb-3 ${displayResult.status === 'Accepted' ? 'text-emerald-700' : 'text-rose-700'
                            }`}>
                            <motion.div
                                initial={{ scale: 0 }}
                                animate={{ scale: 1 }}
                                transition={{ type: "spring", stiffness: 200, damping: 10 }}
                            >
                                {displayResult.status === 'Accepted' ? (
                                    <CheckCircle2 className="h-10 w-10" />
                                ) : (
                                    <XCircle className="h-10 w-10" />
                                )}
                            </motion.div>
                            {displayResult.status}
                        </h2>
                        <p className="text-gray-700 font-medium text-lg">
                            You passed <span className="text-gray-900 font-black">{displayResult.passed_test_cases}</span> out of <span className="text-gray-900 font-black">{displayResult.total_test_cases}</span> test cases.
                            {isRunResult && <span className="text-sm text-gray-500 ml-2">(Public test cases only)</span>}
                        </p>
                    </div>
                    {!isRunResult && 'runtime' in displayResult && 'memory' in displayResult && (
                        <div className="flex gap-4 relative z-10">
                            <motion.div
                                initial={{ scale: 0.9, opacity: 0 }}
                                animate={{ scale: 1, opacity: 1 }}
                                transition={{ delay: 0.1 }}
                                className="bg-white/70 backdrop-blur-sm px-6 py-4 rounded-2xl border border-white shadow-lg text-center"
                            >
                                <div className="text-xs text-gray-600 uppercase font-bold tracking-widest mb-2">Runtime</div>
                                <div className="text-2xl font-black bg-gradient-to-br from-gray-900 to-gray-700 bg-clip-text text-transparent">
                                    {displayResult.runtime}<span className="text-sm font-normal ml-1 text-gray-600">ms</span>
                                </div>
                            </motion.div>
                            <motion.div
                                initial={{ scale: 0.9, opacity: 0 }}
                                animate={{ scale: 1, opacity: 1 }}
                                transition={{ delay: 0.2 }}
                                className="bg-white/70 backdrop-blur-sm px-6 py-4 rounded-2xl border border-white shadow-lg text-center"
                            >
                                <div className="text-xs text-gray-600 uppercase font-bold tracking-widest mb-2">Memory</div>
                                <div className="text-2xl font-black bg-gradient-to-br from-gray-900 to-gray-700 bg-clip-text text-transparent">
                                    {displayResult.memory}<span className="text-sm font-normal ml-1 text-gray-600">KB</span>
                                </div>
                            </motion.div>
                        </div>
                    )}
                </div>

                {displayResult.error_message && (
                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.3 }}
                        className="bg-gradient-to-br from-black to-gray-900 rounded-3xl p-8 overflow-hidden shadow-2xl border border-gray-800"
                    >
                        <div className="flex items-center justify-between mb-6 pb-4 border-b border-gray-800">
                            <div className="flex items-center gap-3">
                                <div className="flex gap-1.5">
                                    <div className="w-3 h-3 rounded-full bg-rose-500 shadow-lg shadow-rose-500/50" />
                                    <div className="w-3 h-3 rounded-full bg-amber-500 shadow-lg shadow-amber-500/50" />
                                    <div className="w-3 h-3 rounded-full bg-emerald-500 shadow-lg shadow-emerald-500/50" />
                                </div>
                                <h3 className="text-gray-400 text-xs font-mono uppercase tracking-widest font-bold">Error Log</h3>
                            </div>
                        </div>
                        <div className="max-h-[300px] overflow-y-auto custom-scrollbar-dark">
                            <pre className="text-rose-400 font-mono text-sm whitespace-pre-wrap leading-relaxed">
                                {displayResult.error_message}
                            </pre>
                        </div>
                    </motion.div>
                )}

                {/* Test Case Results */}
                <div className="space-y-4">
                    <h3 className="text-xl font-bold text-gray-800">Test Cases</h3>
                    <div className="grid gap-4">
                        {((isRunResult ? (displayResult as RunCodeResult).results : (displayResult as Submission).test_case_results) || []).map((result, index) => (
                            <motion.div
                                key={index}
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.1 * index }}
                                className={`p-4 rounded-xl border ${result.status === 'Passed'
                                    ? 'bg-emerald-50 border-emerald-100'
                                    : 'bg-rose-50 border-rose-100'
                                    }`}
                            >
                                <div className="flex items-center justify-between mb-2">
                                    <div className="flex items-center gap-2">
                                        {result.status === 'Passed' ? (
                                            <CheckCircle2 className="h-5 w-5 text-emerald-600" />
                                        ) : (
                                            <XCircle className="h-5 w-5 text-rose-600" />
                                        )}
                                        <span className={`font-semibold ${result.status === 'Passed' ? 'text-emerald-700' : 'text-rose-700'
                                            }`}>
                                            Test Case {index + 1}
                                        </span>
                                        {!result.is_sample && (
                                            <span className="text-xs px-2 py-0.5 rounded-full bg-gray-200 text-gray-600 font-medium">
                                                Hidden
                                            </span>
                                        )}
                                    </div>
                                    <span className={`text-sm font-medium ${result.status === 'Passed' ? 'text-emerald-600' : 'text-rose-600'
                                        }`}>
                                        {result.status}
                                    </span>
                                </div>

                                {result.is_sample && (
                                    <div className="grid grid-cols-3 gap-4 mt-3 text-sm">
                                        <div className="bg-white/50 p-3 rounded-lg">
                                            <div className="text-xs font-bold text-gray-500 uppercase mb-1">Input</div>
                                            <code className="font-mono text-gray-800 whitespace-pre-wrap break-all">{result.input}</code>
                                        </div>
                                        <div className="bg-white/50 p-3 rounded-lg">
                                            <div className="text-xs font-bold text-gray-500 uppercase mb-1">Expected</div>
                                            <code className="font-mono text-gray-800 whitespace-pre-wrap break-all">{result.expected_output}</code>
                                        </div>
                                        <div className={`p-3 rounded-lg ${result.status === 'Passed' ? 'bg-emerald-100/50' : 'bg-rose-100/50'
                                            }`}>
                                            <div className="text-xs font-bold text-gray-500 uppercase mb-1">Actual</div>
                                            <code className="font-mono text-gray-800 whitespace-pre-wrap break-all">{result.actual_output}</code>
                                        </div>
                                    </div>
                                )}
                            </motion.div>
                        ))}
                    </div>
                </div>
            </motion.div>
        )
    }

    return null
}
