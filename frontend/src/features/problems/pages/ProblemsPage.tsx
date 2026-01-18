import { Search, ChevronRight, Star } from 'lucide-react'
import { problemsApi } from '../api/problems'
import { Input } from '@/shared/components/ui/Input'
import { Card } from '@/shared/components/ui/Card'
import type { Difficulty } from '../types'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { Skeleton } from '@/shared/components/ui/Skeleton'

const DIFFICULTY_COLORS: Record<Difficulty, string> = {
    easy: 'text-emerald-500 bg-emerald-500/10',
    medium: 'text-amber-500 bg-amber-500/10',
    hard: 'text-rose-500 bg-rose-500/10',
}

const ProblemsSkeleton = () => (
    <div className="space-y-4">
        {[1, 2, 3, 4, 5].map((i) => (
            <Card key={i} className="p-6 border-gray-100">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-6">
                        <Skeleton className="hidden sm:block h-12 w-12 rounded-xl" />
                        <div className="space-y-2">
                            <Skeleton className="h-7 w-64" />
                            <div className="flex gap-4">
                                <Skeleton className="h-5 w-16 rounded-full" />
                                <Skeleton className="h-5 w-32" />
                            </div>
                        </div>
                    </div>
                </div>
            </Card>
        ))}
    </div>
)

export const ProblemsPage = () => {
    const [search, setSearch] = useState('')
    const [difficulty, setDifficulty] = useState<string>('')
    const [page] = useState(1)

    const { data, isLoading } = useQuery({
        queryKey: ['problems', { page, search, difficulty }],
        queryFn: () => problemsApi.list({ page, search, difficulty: difficulty || undefined }),
    })

    const problems = data?.data.data || []

    return (
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-12">
                <div>
                    <h1 className="text-4xl font-bold text-gray-900 mb-2">Algorithm Challenges</h1>
                    <p className="text-gray-600 italic">Sharpen your skills with our curated set of problems.</p>
                </div>
                <div className="flex flex-col sm:flex-row gap-4">
                    <div className="relative">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
                        <Input
                            type="text"
                            placeholder="Search problems..."
                            className="pl-10 min-w-[300px]"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                        />
                    </div>
                    <select
                        className="px-4 py-2 rounded-lg border border-gray-200 bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                        value={difficulty}
                        onChange={(e) => setDifficulty(e.target.value)}
                    >
                        <option value="">All Difficulties</option>
                        <option value="easy">Easy</option>
                        <option value="medium">Medium</option>
                        <option value="hard">Hard</option>
                    </select>
                </div>
            </div>

            {isLoading ? (
                <ProblemsSkeleton />
            ) : (
                <div className="space-y-4">
                    {problems.map((problem) => (
                        <Link
                            key={problem.id}
                            to={`/problems/${problem.slug}`}
                            className="block group"
                        >
                            <Card className="p-6 transition-all duration-300 hover:shadow-xl hover:-translate-y-1 border-gray-100 group-hover:border-blue-200">
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-6">
                                        <div className="hidden sm:flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600">
                                            <Star className="h-6 w-6" />
                                        </div>
                                        <div>
                                            <h3 className="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors">
                                                {problem.id}. {problem.title}
                                            </h3>
                                            <div className="flex items-center gap-4 mt-1">
                                                <span className={`px-3 py-0.5 rounded-full text-xs font-semibold capitalize ${DIFFICULTY_COLORS[problem.difficulty]}`}>
                                                    {problem.difficulty}
                                                </span>
                                                <span className="text-sm text-gray-500 flex items-center gap-1">
                                                    Acceptance: {problem.acceptance_rate.toFixed(1)}%
                                                </span>
                                                {problem.creator && (
                                                    <span className="text-sm text-gray-400">
                                                        by <Link
                                                            to={`/users/${problem.creator.username}`}
                                                            onClick={(e) => e.stopPropagation()}
                                                            className="hover:text-blue-500 font-medium transition-colors"
                                                        >
                                                            @{problem.creator.username}
                                                        </Link>
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-4">
                                        {problem.user_status === 'solved' && (
                                            <span className="flex items-center gap-1.5 px-3 py-1 bg-emerald-100 text-emerald-700 rounded-full text-xs font-bold uppercase tracking-wider animate-in fade-in zoom-in duration-300">
                                                <div className="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse" />
                                                Solved
                                            </span>
                                        )}
                                        {problem.user_status === 'attempted' && (
                                            <span className="flex items-center gap-1.5 px-3 py-1 bg-amber-100 text-amber-700 rounded-full text-xs font-bold uppercase tracking-wider animate-in fade-in zoom-in duration-300">
                                                Attempted
                                            </span>
                                        )}
                                        <div className="text-gray-300 group-hover:text-blue-500 transition-colors transform group-hover:translate-x-1 duration-300">
                                            <ChevronRight className="h-6 w-6" />
                                        </div>
                                    </div>
                                </div>
                            </Card>
                        </Link>
                    ))}

                    {problems.length === 0 && (
                        <div className="text-center py-20 bg-gray-50 rounded-2xl border-2 border-dashed border-gray-200">
                            <Search className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                            <h3 className="text-lg font-medium text-gray-900">No problems found</h3>
                            <p className="text-gray-500">Try adjusting your filters or search query.</p>
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}
