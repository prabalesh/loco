import { motion } from 'framer-motion'
import { Play, CheckCircle2, XCircle, Clock, Loader2, Zap, Database, Cpu } from 'lucide-react'
import type { Submission, TestCase } from '../types'

interface ResultTabProps {
    submissionResult?: Submission
    pollingId: number | null
    runResult?: Submission | null
    isRunning?: boolean
    sampleTestCases?: TestCase[]
    pollingType: 'run' | 'submit' | null
}

export const ResultTab = ({ submissionResult, pollingId, runResult, isRunning, sampleTestCases, pollingType }: ResultTabProps) => {
    // Determine what to display
    const isPending = isRunning || !!pollingId || (submissionResult?.status === 'Pending' || submissionResult?.status === 'Processing')
    const hasData = !!runResult || (!!submissionResult && submissionResult.status !== 'Pending' && submissionResult.status !== 'Processing')

    // Priority: 
    // 1. If currently polling or running AND we don't have final data yet -> Show Pending State
    // 2. If we have final data (run or submission) -> Show Result State
    // 3. Otherwise -> Show Initial "Not Run" State

    // 1. Pending/Loading State
    if (isPending && !hasData) {
        return (
            <div className="space-y-6">
                <div className="flex flex-col items-center justify-center py-20 px-4 text-center">
                    <div className="relative mb-10">
                        {/* Outer Glow */}
                        <motion.div
                            animate={{
                                scale: [1, 1.2, 1],
                                opacity: [0.1, 0.3, 0.1]
                            }}
                            transition={{ duration: 2, repeat: Infinity }}
                            className="absolute -inset-8 bg-blue-500 rounded-full blur-3xl"
                        />

                        {/* Spinning Loader */}
                        <div className="relative w-24 h-24">
                            <svg className="w-full h-full" viewBox="0 0 100 100">
                                <circle
                                    className="text-blue-100"
                                    strokeWidth="6"
                                    stroke="currentColor"
                                    fill="transparent"
                                    r="38"
                                    cx="50"
                                    cy="50"
                                />
                                <motion.circle
                                    className="text-blue-600"
                                    strokeWidth="6"
                                    strokeDasharray="238.76"
                                    strokeDashoffset="100"
                                    strokeLinecap="round"
                                    stroke="currentColor"
                                    fill="transparent"
                                    r="38"
                                    cx="50"
                                    cy="50"
                                    animate={{ strokeDashoffset: [238.76, 0] }}
                                    transition={{ duration: 2, repeat: Infinity, ease: "linear" }}
                                />
                            </svg>
                            <div className="absolute inset-0 flex items-center justify-center">
                                {pollingId ? (
                                    <Clock className="h-8 w-8 text-blue-600 animate-pulse" />
                                ) : (
                                    <Play className="h-8 w-8 text-blue-600 animate-pulse fill-blue-600/20" />
                                )}
                            </div>
                        </div>
                    </div>

                    <h3 className="text-2xl font-black text-gray-900 mb-3 tracking-tight">
                        {pollingType === 'submit' ? 'Evaluating Submission...' : 'Running Code...'}
                    </h3>
                    <div className="max-w-md space-y-2">
                        <p className="text-gray-500 font-medium">
                            {pollingId
                                ? "Your code is being tested against all test cases in the cloud."
                                : "Executing your solution against sample test cases."}
                        </p>
                        <div className="flex items-center justify-center gap-1 text-xs font-bold text-gray-400 uppercase tracking-widest mt-4">
                            <Loader2 className="h-3 w-3 animate-spin" />
                            <span>Processing...</span>
                        </div>
                    </div>
                </div>

                {/* Faded background test cases to maintain layout weight */}
                <div className="space-y-4 opacity-30 select-none pointer-events-none">
                    <h3 className="text-xl font-bold text-gray-400">Sample Test Cases</h3>
                    <div className="grid gap-4">
                        {(sampleTestCases || []).slice(0, 2).map((_, index) => (
                            <div key={index} className="p-6 rounded-2xl border border-gray-100 bg-white shadow-sm">
                                <div className="h-4 w-32 bg-gray-100 rounded animate-pulse mb-4" />
                                <div className="grid grid-cols-2 gap-4">
                                    <div className="h-12 bg-gray-50 rounded-xl animate-pulse" />
                                    <div className="h-12 bg-gray-50 rounded-xl animate-pulse" />
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        )
    }

    // 2. Result State
    if (hasData) {
        const displayResult = runResult || submissionResult!
        const isRunResult = displayResult.is_run_only
        const submissionData = displayResult as Submission
        const status = displayResult.status
        const isAccepted = status === 'Accepted' // Handle both casing versions

        return (
            <motion.div
                initial={{ opacity: 0, scale: 0.98 }}
                animate={{ opacity: 1, scale: 1 }}
                className="space-y-8"
            >
                {/* Result Header Card */}
                <div className={`relative overflow-hidden p-8 rounded-[2rem] border-2 shadow-2xl transition-all duration-500 ${isAccepted
                    ? 'bg-gradient-to-br from-emerald-50 to-emerald-100 border-emerald-200 shadow-emerald-200/20'
                    : 'bg-gradient-to-br from-rose-50 to-rose-100 border-rose-200 shadow-rose-200/20'
                    }`}>
                    {/* Background Decorative Element */}
                    <div className={`absolute -top-24 -right-24 w-64 h-64 rounded-full blur-[80px] opacity-40 ${isAccepted ? 'bg-emerald-400' : 'bg-rose-400'
                        }`} />

                    <div className="relative flex flex-col md:flex-row md:items-center justify-between gap-8">
                        <div className="flex-1 space-y-4">
                            <div className="flex items-center gap-4">
                                <motion.div
                                    initial={{ scale: 0.5, rotate: -45 }}
                                    animate={{ scale: 1, rotate: 0 }}
                                    className={`p-3 rounded-2xl ${isAccepted ? 'bg-emerald-600 text-white' : 'bg-rose-600 text-white'}`}
                                >
                                    {isAccepted ? <CheckCircle2 className="h-8 w-8" /> : <XCircle className="h-8 w-8" />}
                                </motion.div>
                                <div>
                                    <h2 className={`text-4xl font-black tracking-tight ${isAccepted ? 'text-emerald-900' : 'text-rose-900'}`}>
                                        {status}
                                    </h2>
                                    <p className={`font-bold mt-1 ${isAccepted ? 'text-emerald-700/70' : 'text-rose-700/70'}`}>
                                        {isRunResult ? 'Public Tests Only' : 'Final Evaluation'}
                                    </p>
                                </div>
                            </div>

                            <div className="flex flex-wrap items-center gap-x-6 gap-y-2">
                                <div className="flex items-center gap-2">
                                    <span className="text-gray-500 text-sm font-bold uppercase tracking-widest">Score:</span>
                                    <span className={`text-xl font-black ${isAccepted ? 'text-emerald-600' : 'text-rose-600'}`}>
                                        {displayResult.passed_test_cases} / {displayResult.total_test_cases}
                                    </span>
                                </div>
                                <div className="h-4 w-[2px] bg-gray-200 hidden md:block" />
                                <div className="flex items-center gap-2">
                                    <span className="text-gray-500 text-sm font-bold uppercase tracking-widest">Accuracy:</span>
                                    <span className="text-xl font-black text-gray-800">
                                        {Math.round((displayResult.passed_test_cases / displayResult.total_test_cases) * 100)}%
                                    </span>
                                </div>
                            </div>
                        </div>

                        {!isRunResult && (
                            <div className="grid grid-cols-2 gap-4">
                                <div className="bg-white/80 backdrop-blur-md p-5 rounded-2xl border border-white shadow-sm flex flex-col items-center justify-center min-w-[120px]">
                                    <Cpu className="h-5 w-5 text-indigo-500 mb-2" />
                                    <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-1">Runtime</div>
                                    <div className="text-2xl font-black text-gray-900">{submissionData.runtime || 0}ms</div>
                                </div>
                                <div className="bg-white/80 backdrop-blur-md p-5 rounded-2xl border border-white shadow-sm flex flex-col items-center justify-center min-w-[120px]">
                                    <Database className="h-5 w-5 text-amber-500 mb-2" />
                                    <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-1">Memory</div>
                                    <div className="text-2xl font-black text-gray-900">{submissionData.memory || 0}KB</div>
                                </div>
                            </div>
                        )}
                    </div>
                </div>

                {/* Error Box if any */}
                {displayResult.error_message && (
                    <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="bg-gray-950 rounded-3xl overflow-hidden shadow-2xl border border-gray-800"
                    >
                        <div className="px-6 py-4 bg-gray-900 border-b border-gray-800 flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <div className="flex gap-1.5 mr-2">
                                    <div className="w-2.5 h-2.5 rounded-full bg-rose-500" />
                                    <div className="w-2.5 h-2.5 rounded-full bg-amber-500" />
                                    <div className="w-2.5 h-2.5 rounded-full bg-emerald-500" />
                                </div>
                                <span className="text-xs font-mono text-gray-400 uppercase tracking-widest font-black">Execution Error</span>
                            </div>
                        </div>
                        <div className="p-6 overflow-x-auto">
                            <pre className="text-rose-400 font-mono text-sm leading-relaxed whitespace-pre-wrap">
                                {displayResult.error_message}
                            </pre>
                        </div>
                    </motion.div>
                )}

                {/* Detailed Test Results */}
                <div className="space-y-6">
                    <div className="flex items-center justify-between">
                        <h3 className="text-2xl font-black text-gray-900 tracking-tight">Detailed Test Cases</h3>
                        <div className="flex items-center gap-2 text-sm font-bold text-gray-500 px-4 py-2 bg-gray-100 rounded-full">
                            <Zap className="h-4 w-4 text-amber-500 fill-amber-500" />
                            <span>{displayResult.passed_test_cases} Passed</span>
                        </div>
                    </div>

                    <div className="grid gap-4">
                        {(
                            displayResult.test_case_results || []
                        )
                            .sort((a, b) => (b.is_sample ? 1 : 0) - (a.is_sample ? 1 : 0))
                            .map((result, index) => {
                                const passed = result.status === 'passed' || result.status === 'Passed'
                                return (
                                    <motion.div
                                        key={index}
                                        initial={{ opacity: 0, x: -10 }}
                                        animate={{ opacity: 1, x: 0 }}
                                        transition={{ delay: index * 0.05 }}
                                        className={`group rounded-2xl border-2 transition-all duration-300 overflow-hidden ${passed
                                            ? 'bg-white border-emerald-100 hover:border-emerald-300 hover:shadow-lg hover:shadow-emerald-500/5'
                                            : 'bg-white border-rose-100 hover:border-rose-300 hover:shadow-lg hover:shadow-rose-500/5'
                                            }`}
                                    >
                                        <div className="p-6">
                                            <div className="flex items-center justify-between mb-6">
                                                <div className="flex items-center gap-4">
                                                    <div className={`w-10 h-10 rounded-xl flex items-center justify-center font-black ${passed ? 'bg-emerald-100 text-emerald-600' : 'bg-rose-100 text-rose-600'
                                                        }`}>
                                                        {index + 1}
                                                    </div>
                                                    <div>
                                                        <div className="flex items-center gap-2">
                                                            <span className="font-black text-gray-900">Test Case {index + 1}</span>
                                                            {result.is_sample && (
                                                                <span className="text-[10px] font-black uppercase tracking-widest bg-blue-100 text-blue-600 px-2 py-0.5 rounded-md">Sample</span>
                                                            )}
                                                            {!result.is_sample && (
                                                                <span className="text-[10px] font-black uppercase tracking-widest bg-gray-100 text-gray-500 px-2 py-0.5 rounded-md">Hidden</span>
                                                            )}
                                                        </div>
                                                    </div>
                                                </div>
                                                <div className={`px-4 py-1.5 rounded-xl text-xs font-black uppercase tracking-widest ${passed ? 'bg-emerald-50 text-emerald-600' : 'bg-rose-50 text-rose-600'
                                                    }`}>
                                                    {result.status}
                                                </div>
                                            </div>

                                            {result.is_sample ? (
                                                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                                    <div className="space-y-2">
                                                        <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Input</div>
                                                        <div className="bg-gray-50 p-3 rounded-xl border border-gray-100 font-mono text-sm text-gray-700 whitespace-pre-wrap break-all min-h-[44px]">
                                                            {result.input}
                                                        </div>
                                                    </div>
                                                    <div className="space-y-2">
                                                        <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Expected</div>
                                                        <div className="bg-gray-50 p-3 rounded-xl border border-gray-100 font-mono text-sm text-gray-700 whitespace-pre-wrap break-all min-h-[44px]">
                                                            {result.expected_output}
                                                        </div>
                                                    </div>
                                                    <div className="space-y-2">
                                                        <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Actual</div>
                                                        <div className={`p-3 rounded-xl border font-mono text-sm whitespace-pre-wrap break-all min-h-[44px] ${passed ? 'bg-emerald-50 border-emerald-100 text-emerald-700' : 'bg-rose-50 border-rose-100 text-rose-700 font-bold'
                                                            }`}>
                                                            {result.actual_output || 'No output'}
                                                        </div>
                                                    </div>
                                                </div>
                                            ) : (
                                                <div className="py-4 px-6 bg-gray-50 rounded-2xl border border-dashed border-gray-200 text-center">
                                                    <p className="text-gray-400 text-sm font-medium italic">Detailed inputs/outputs are hidden for this test case</p>
                                                </div>
                                            )}
                                        </div>
                                    </motion.div>
                                )
                            })}
                    </div>
                </div>
            </motion.div>
        )
    }

    // 3. Final Fallback: Initial State (Not Run)
    return (
        <div className="space-y-8">
            <div className="flex flex-col items-center justify-center py-20 bg-gradient-to-b from-gray-50/50 to-white rounded-[2.5rem] border-2 border-dashed border-gray-200">
                <div className="w-20 h-20 rounded-3xl bg-white shadow-xl flex items-center justify-center mb-6 overflow-hidden relative group">
                    <motion.div
                        animate={{ scale: [1, 1.1, 1] }}
                        transition={{ duration: 3, repeat: Infinity }}
                        className="absolute inset-0 bg-blue-500/5 group-hover:bg-blue-500/10 transition-colors"
                    />
                    <Play className="h-10 w-10 text-blue-500 fill-blue-500/20 relative z-10" />
                </div>
                <h3 className="text-2xl font-black text-gray-900 mb-2">Ready to test?</h3>
                <p className="text-gray-500 font-medium text-center max-w-sm px-6">
                    Run your solution against sample test cases or submit it for a full evaluation.
                </p>
            </div>

            <div className="space-y-6">
                <h3 className="text-2xl font-black text-gray-900 tracking-tight">Sample Test Cases</h3>
                <div className="grid gap-4">
                    {(sampleTestCases || []).map((tc, index) => (
                        <div key={tc.id} className="p-7 rounded-3xl border border-gray-200 bg-white shadow-sm hover:shadow-md transition-all duration-300">
                            <div className="flex items-center justify-between mb-5">
                                <div className="flex items-center gap-4">
                                    <div className="w-10 h-10 rounded-xl bg-gray-100 flex items-center justify-center text-gray-700 font-black text-lg">
                                        {index + 1}
                                    </div>
                                    <span className="font-black text-gray-800 text-lg">Test Case {index + 1}</span>
                                </div>
                                <span className="text-[10px] font-black uppercase tracking-widest text-gray-400 bg-gray-50 border border-gray-100 px-4 py-1.5 rounded-[10px]">
                                    Not Run
                                </span>
                            </div>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div className="space-y-2.5">
                                    <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Sample Input</div>
                                    <div className="bg-gray-50 p-4 rounded-2xl border border-gray-100">
                                        <code className="font-mono text-gray-700 text-sm whitespace-pre-wrap break-all">
                                            {tc.input}
                                        </code>
                                    </div>
                                </div>
                                <div className="space-y-2.5">
                                    <div className="text-[10px] font-black text-gray-400 uppercase tracking-widest ml-1">Expected Output</div>
                                    <div className="bg-gray-50 p-4 rounded-2xl border border-gray-100">
                                        <code className="font-mono text-gray-700 text-sm whitespace-pre-wrap break-all">
                                            {tc.expected_output}
                                        </code>
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    )
}
