import { motion } from 'framer-motion'
import { Link } from 'react-router-dom'
import type { ProblemResponse } from '../types'
import { calculateAcceptanceRate } from '@/lib/utils'

interface DescriptionTabProps {
    problem: ProblemResponse
}

export const DescriptionTab = ({ problem }: DescriptionTabProps) => {
    return (
        <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
            className="prose prose-blue max-w-none"
        >
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
                <div className="bg-gray-50 rounded-xl p-4 border border-gray-100 shadow-sm">
                    <p className="text-sm font-medium text-gray-500 uppercase tracking-wider mb-1">Acceptance Rate</p>
                    <p className="text-2xl font-bold text-gray-900">{calculateAcceptanceRate(problem.total_accepted, problem.total_submissions).toFixed(1)}%</p>
                </div>
                <div className="bg-gray-50 rounded-xl p-4 border border-gray-100 shadow-sm">
                    <p className="text-sm font-medium text-gray-500 uppercase tracking-wider mb-1">Total Submissions</p>
                    <p className="text-2xl font-bold text-gray-900">{problem.total_submissions.toLocaleString()}</p>
                </div>
                <div className="bg-gray-50 rounded-xl p-4 border border-gray-100 shadow-sm">
                    <p className="text-sm font-medium text-gray-500 uppercase tracking-wider mb-1">Total Accepted</p>
                    <p className="text-2xl font-bold text-gray-900">{problem.total_accepted.toLocaleString()}</p>
                </div>
            </div>

            {/* <div className="flex flex-wrap gap-2 mb-8">
                {problem.categories?.map(cat => (
                    <span key={cat.id} className="px-3 py-1 bg-blue-50 text-blue-600 rounded-lg text-xs font-bold uppercase tracking-wider border border-blue-100 shadow-sm">
                        {cat.name}
                    </span>
                ))}
                {problem.tags?.map(tag => (
                    <span key={tag.id} className="px-3 py-1 bg-gray-50 text-gray-600 rounded-lg text-xs font-medium border border-gray-100 shadow-sm">
                        {tag.name}
                    </span>
                ))}
            </div> */}

            {problem.description && (
                <div className="mb-8">
                    <h2 className="text-xl font-bold mb-4 text-gray-900 border-b border-gray-200 pb-2">
                        Description
                    </h2>
                    <div
                        className="text-gray-700 leading-relaxed"
                        dangerouslySetInnerHTML={{ __html: problem.description }}
                    />
                </div>
            )}

            {/* {problem.input_format && (
                <div className="mb-8">
                    <h2 className="text-xl font-bold mb-4 text-gray-900 border-b border-gray-200 pb-2">
                        Input Format
                    </h2>
                    <div
                        className="text-gray-700 leading-relaxed"
                        dangerouslySetInnerHTML={{ __html: problem.input_format }}
                    />
                </div>
            )}

            {problem.output_format && (
                <div className="mb-8">
                    <h2 className="text-xl font-bold mb-4 text-gray-900 border-b border-gray-200 pb-2">
                        Output Format
                    </h2>
                    <div
                        className="text-gray-700 leading-relaxed"
                        dangerouslySetInnerHTML={{ __html: problem.output_format }}
                    />
                </div>
            )}

            {problem.constraints && (
                <div className="mb-8">
                    <h2 className="text-xl font-bold mb-4 text-gray-900 border-b border-gray-200 pb-2">
                        Constraints
                    </h2>
                    <div
                        className="text-gray-700 leading-relaxed"
                        dangerouslySetInnerHTML={{ __html: problem.constraints }}
                    />
                </div>
            )} */}

            {problem.creator && (
                <div className="mt-12 pt-8 border-t border-gray-100">
                    <div className="flex items-center justify-between bg-gray-50/50 p-4 rounded-2xl border border-gray-100">
                        <div className="flex items-center gap-3">
                            <div className="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold">
                                {problem.creator.username.charAt(0).toUpperCase()}
                            </div>
                            <div>
                                <p className="text-xs font-bold text-gray-400 uppercase tracking-widest">Problem Author</p>
                                <Link
                                    to={`/users/${problem.creator.username}`}
                                    className="text-sm font-bold text-gray-900 hover:text-blue-600 transition-colors"
                                >
                                    @{problem.creator.username}
                                </Link>
                            </div>
                        </div>
                        <Link
                            to={`/users/${problem.creator.username}`}
                            className="text-xs font-bold text-blue-600 bg-blue-50 px-3 py-1.5 rounded-lg hover:bg-blue-100 transition-colors"
                        >
                            View Profile
                        </Link>
                    </div>
                </div>
            )}
        </motion.div>
    )
}
