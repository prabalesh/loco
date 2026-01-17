import { motion } from 'framer-motion'
import type { Problem } from '../types'

interface DescriptionTabProps {
    problem: Problem
}

export const DescriptionTab = ({ problem }: DescriptionTabProps) => {
    return (
        <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
            className="prose prose-blue max-w-none"
        >
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

            {problem.input_format && (
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
            )}
        </motion.div>
    )
}
