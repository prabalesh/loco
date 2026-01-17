import { Send, Play, ChevronLeft, Loader2 } from 'lucide-react'
import { Button } from '@/shared/components/ui/Button'
import type { Problem } from '../types'

interface ProblemHeaderProps {
    problem: Problem
    onBack: () => void
    onRun: () => void
    onSubmit: () => void
    isSubmitting: boolean
}

export const ProblemHeader = ({ problem, onBack, onRun, onSubmit, isSubmitting }: ProblemHeaderProps) => {
    return (
        <header className="bg-white/95 backdrop-blur-md border-b border-gray-200 px-6 py-3 flex items-center justify-between shadow-lg z-20">
            <div className="flex items-center gap-4">
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onBack}
                    className="text-gray-500 hover:text-gray-900 hover:bg-gray-100 transition-all duration-200"
                >
                    <ChevronLeft className="h-5 w-5 mr-1" />
                    Back
                </Button>
                <div className="h-6 w-px bg-gray-300" />
                <div className="flex items-center gap-3">
                    <h1 className="text-lg font-bold text-gray-900 truncate max-w-[300px]">
                        {problem.id}. {problem.title}
                    </h1>
                    <span className={`text-xs font-bold px-3 py-1 rounded-full capitalize shadow-sm ${problem.difficulty === 'easy'
                            ? 'text-emerald-700 bg-gradient-to-r from-emerald-50 to-emerald-100 border border-emerald-200' :
                            problem.difficulty === 'medium'
                                ? 'text-amber-700 bg-gradient-to-r from-amber-50 to-amber-100 border border-amber-200' :
                                'text-rose-700 bg-gradient-to-r from-rose-50 to-rose-100 border border-rose-200'
                        }`}>
                        {problem.difficulty}
                    </span>
                </div>
            </div>

            <div className="flex items-center gap-3">
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={onRun}
                    disabled={isSubmitting}
                    className="shadow-md bg-gradient-to-r from-gray-100 to-gray-50 hover:from-gray-200 hover:to-gray-100 text-gray-700 hover:text-gray-900 px-6 transition-all duration-200 hover:shadow-lg border border-gray-200"
                >
                    {isSubmitting ? (
                        <Loader2 className="h-4 w-4 animate-spin mr-2" />
                    ) : (
                        <Play className="h-4 w-4 mr-2" />
                    )}
                    Run
                </Button>
                <Button
                    variant="primary"
                    size="sm"
                    onClick={onSubmit}
                    disabled={isSubmitting}
                    className="shadow-lg bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-700 hover:to-blue-600 px-8 transition-all duration-200 hover:shadow-xl hover:scale-105"
                >
                    {isSubmitting ? (
                        <Loader2 className="h-4 w-4 animate-spin mr-2" />
                    ) : (
                        <Send className="h-4 w-4 mr-2" />
                    )}
                    Submit
                </Button>
            </div>
        </header>
    )
}
