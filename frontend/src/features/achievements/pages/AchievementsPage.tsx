import { useQuery } from '@tanstack/react-query'
import { achievementsApi, type Achievement } from '../api/achievementsApi'
import { Trophy, Lock, Award, Target, Zap, TrendingUp, CheckCircle } from 'lucide-react'
import { useAuth } from '@/shared/hooks/useAuth'

const categoryIcons: Record<string, typeof Trophy> = {
    'getting-started': Award,
    'solving': Target,
    'difficulty': TrendingUp,
    'streak': Zap,
    'language': Trophy,
    'misc': Trophy,
}

const categoryNames: Record<string, string> = {
    'getting-started': 'Getting Started',
    'solving': 'Problem Solving',
    'difficulty': 'Difficulty Mastery',
    'streak': 'Streaks',
    'language': 'Language Mastery',
    'misc': 'Miscellaneous',
}

const categoryColors: Record<string, { bg: string; border: string; text: string; icon: string }> = {
    'getting-started': { bg: 'bg-blue-50', border: 'border-blue-200', text: 'text-blue-900', icon: 'text-blue-600' },
    'solving': { bg: 'bg-emerald-50', border: 'border-emerald-200', text: 'text-emerald-900', icon: 'text-emerald-600' },
    'difficulty': { bg: 'bg-purple-50', border: 'border-purple-200', text: 'text-purple-900', icon: 'text-purple-600' },
    'streak': { bg: 'bg-amber-50', border: 'border-amber-200', text: 'text-amber-900', icon: 'text-amber-600' },
    'language': { bg: 'bg-pink-50', border: 'border-pink-200', text: 'text-pink-900', icon: 'text-pink-600' },
    'misc': { bg: 'bg-slate-50', border: 'border-slate-200', text: 'text-slate-900', icon: 'text-slate-600' },
}

// Helper to calculate progress for an achievement
const calculateProgress = (achievement: Achievement, userStats: any): { current: number; target: number; percentage: number } => {
    const conditionValue = achievement.condition_value

    // Parse condition value
    let target = 1
    try {
        if (achievement.condition_type === 'difficulty_count') {
            const parsed = JSON.parse(conditionValue)
            target = parsed.count || 1
        } else if (achievement.condition_type === 'count' || achievement.condition_type === 'streak') {
            target = parseInt(conditionValue) || 1
        }
    } catch {
        target = parseInt(conditionValue) || 1
    }

    // Calculate current progress based on achievement type
    let current = 0

    if (achievement.slug === 'hello-world' || achievement.slug === 'first-blood') {
        current = userStats?.problems_solved || 0
    } else if (achievement.category === 'solving') {
        current = userStats?.problems_solved || 0
    } else if (achievement.slug.includes('bug-hunter') || achievement.slug.includes('speed-demon') || achievement.slug.includes('memory-leak')) {
        current = userStats?.total_submissions || 0
    } else {
        // Default to problems solved
        current = userStats?.problems_solved || 0
    }

    const percentage = Math.min((current / target) * 100, 100)

    return { current, target, percentage }
}

export function AchievementsPage() {
    const { user } = useAuth()

    const { data: allAchievements = [], isLoading: loadingAll } = useQuery({
        queryKey: ['achievements'],
        queryFn: achievementsApi.list,
    })

    const { data: userAchievements = [], isLoading: loadingUser } = useQuery({
        queryKey: ['my-achievements'],
        queryFn: achievementsApi.getMyAchievements,
        enabled: !!user,
    })

    const unlockedIds = new Set(userAchievements.map(ua => ua.achievement_id))

    // Group achievements by category
    const groupedAchievements = allAchievements.reduce((acc, achievement) => {
        const category = achievement.category || 'misc'
        if (!acc[category]) {
            acc[category] = []
        }
        acc[category].push(achievement)
        return acc
    }, {} as Record<string, Achievement[]>)

    if (loadingAll || loadingUser) {
        return (
            <div className="min-h-screen bg-white p-6">
                <div className="max-w-7xl mx-auto">
                    <div className="animate-pulse space-y-8">
                        <div className="h-12 bg-slate-200 rounded w-64"></div>
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                            {[...Array(6)].map((_, i) => (
                                <div key={i} className="h-40 bg-slate-200 rounded-lg"></div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        )
    }

    const totalAchievements = allAchievements.length
    const unlockedCount = userAchievements.length
    const totalXP = userAchievements.reduce((sum, ua) => sum + ua.achievement.xp_reward, 0)

    return (
        <div className="min-h-screen bg-white p-6">
            <div className="max-w-7xl mx-auto space-y-8">
                {/* Header */}
                <div className="text-center space-y-4">
                    <div className="inline-flex items-center justify-center p-3 bg-gradient-to-br from-purple-500 to-pink-500 rounded-2xl shadow-lg mb-4">
                        <Trophy className="w-8 h-8 text-white" />
                    </div>
                    <h1 className="text-5xl font-black text-slate-900">
                        Achievements
                    </h1>
                    <p className="text-lg text-slate-600 max-w-2xl mx-auto">
                        Unlock badges by solving problems and reaching milestones. Track your progress and see how close you are to your next achievement!
                    </p>
                </div>

                {/* Stats Overview */}
                {user && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
                        <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-2xl p-6 border-2 border-purple-200">
                            <div className="flex items-center gap-3 mb-3">
                                <div className="p-3 bg-white rounded-xl shadow-sm">
                                    <Trophy className="w-6 h-6 text-purple-600" />
                                </div>
                                <div>
                                    <p className="text-sm font-bold text-purple-700 uppercase tracking-wide">Unlocked</p>
                                    <p className="text-3xl font-black text-purple-900">
                                        {unlockedCount}<span className="text-lg text-purple-600">/{totalAchievements}</span>
                                    </p>
                                </div>
                            </div>
                            <div className="mt-4">
                                <div className="h-2 bg-white/50 rounded-full overflow-hidden">
                                    <div
                                        className="h-full bg-gradient-to-r from-purple-500 to-pink-500 transition-all duration-500"
                                        style={{ width: `${(unlockedCount / totalAchievements) * 100}%` }}
                                    />
                                </div>
                            </div>
                        </div>

                        <div className="bg-gradient-to-br from-amber-50 to-amber-100 rounded-2xl p-6 border-2 border-amber-200">
                            <div className="flex items-center gap-3">
                                <div className="p-3 bg-white rounded-xl shadow-sm">
                                    <Zap className="w-6 h-6 text-amber-600" />
                                </div>
                                <div>
                                    <p className="text-sm font-bold text-amber-700 uppercase tracking-wide">Total XP</p>
                                    <p className="text-3xl font-black text-amber-900">
                                        {totalXP.toLocaleString()}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-2xl p-6 border-2 border-blue-200">
                            <div className="flex items-center gap-3">
                                <div className="p-3 bg-white rounded-xl shadow-sm">
                                    <Target className="w-6 h-6 text-blue-600" />
                                </div>
                                <div>
                                    <p className="text-sm font-bold text-blue-700 uppercase tracking-wide">Progress</p>
                                    <p className="text-3xl font-black text-blue-900">
                                        {Math.round((unlockedCount / totalAchievements) * 100)}%
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {/* Achievements by Category */}
                <div className="space-y-12">
                    {Object.entries(groupedAchievements).map(([category, achievements]) => {
                        const Icon = categoryIcons[category] || Trophy
                        const categoryName = categoryNames[category] || category
                        const colors = categoryColors[category] || categoryColors['misc']
                        const unlockedInCategory = achievements.filter(a => unlockedIds.has(a.id)).length

                        return (
                            <div key={category} className="space-y-6">
                                <div className="flex items-center gap-4">
                                    <div className={`p-3 ${colors.bg} rounded-xl border-2 ${colors.border}`}>
                                        <Icon className={`w-6 h-6 ${colors.icon}`} />
                                    </div>
                                    <div>
                                        <h2 className="text-3xl font-black text-slate-900">
                                            {categoryName}
                                        </h2>
                                        <p className="text-sm font-semibold text-slate-500">
                                            {unlockedInCategory} of {achievements.length} unlocked
                                        </p>
                                    </div>
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                                    {achievements.map(achievement => {
                                        const isUnlocked = unlockedIds.has(achievement.id)
                                        const userAchievement = userAchievements.find(ua => ua.achievement_id === achievement.id)
                                        const progress = calculateProgress(achievement, user?.stats)

                                        return (
                                            <div
                                                key={achievement.id}
                                                className={`
                          relative overflow-hidden rounded-2xl p-6 border-2 transition-all duration-300 hover:shadow-lg
                          ${isUnlocked
                                                        ? `${colors.bg} ${colors.border} shadow-md`
                                                        : 'bg-white border-slate-200 hover:border-slate-300'
                                                    }
                        `}
                                            >
                                                {/* Achievement Icon */}
                                                <div className="flex items-start justify-between mb-4">
                                                    <div className={`
                            p-3 rounded-xl transition-all
                            ${isUnlocked
                                                            ? 'bg-gradient-to-br from-purple-500 to-pink-500 shadow-lg'
                                                            : 'bg-slate-100 border-2 border-slate-200'
                                                        }
                          `}>
                                                        {isUnlocked ? (
                                                            <Trophy className="w-6 h-6 text-white" />
                                                        ) : (
                                                            <Lock className="w-6 h-6 text-slate-400" />
                                                        )}
                                                    </div>

                                                    {isUnlocked && (
                                                        <div className="flex items-center gap-1 px-3 py-1 bg-white/80 backdrop-blur-sm rounded-full border border-purple-200">
                                                            <CheckCircle className="w-4 h-4 text-emerald-500" />
                                                            <span className="text-xs font-bold text-slate-700">Unlocked</span>
                                                        </div>
                                                    )}
                                                </div>

                                                {/* Achievement Info */}
                                                <div className="space-y-3">
                                                    <h3 className={`text-lg font-bold ${isUnlocked ? colors.text : 'text-slate-900'}`}>
                                                        {achievement.name}
                                                    </h3>

                                                    <p className={`text-sm ${isUnlocked ? 'text-slate-700' : 'text-slate-600'}`}>
                                                        {achievement.description}
                                                    </p>

                                                    {/* XP Reward */}
                                                    <div className="flex items-center gap-2 pt-2">
                                                        <Zap className={`w-4 h-4 ${isUnlocked ? 'text-amber-500' : 'text-slate-400'}`} />
                                                        <span className={`text-sm font-bold ${isUnlocked ? 'text-amber-600' : 'text-slate-500'}`}>
                                                            {achievement.xp_reward} XP
                                                        </span>
                                                    </div>

                                                    {/* Progress Bar (for locked achievements) */}
                                                    {!isUnlocked && user && (
                                                        <div className="pt-3 space-y-2">
                                                            <div className="flex items-center justify-between text-xs font-semibold">
                                                                <span className="text-slate-600">Progress</span>
                                                                <span className="text-slate-900">
                                                                    {progress.current} / {progress.target}
                                                                </span>
                                                            </div>
                                                            <div className="h-2 bg-slate-100 rounded-full overflow-hidden border border-slate-200">
                                                                <div
                                                                    className="h-full bg-gradient-to-r from-purple-500 to-pink-500 transition-all duration-500"
                                                                    style={{ width: `${progress.percentage}%` }}
                                                                />
                                                            </div>
                                                            <p className="text-xs font-medium text-slate-500">
                                                                {progress.percentage >= 100
                                                                    ? '🎉 Ready to unlock!'
                                                                    : `${Math.round(progress.percentage)}% complete`
                                                                }
                                                            </p>
                                                        </div>
                                                    )}

                                                    {/* Unlock Date (for unlocked achievements) */}
                                                    {isUnlocked && userAchievement && (
                                                        <div className="pt-2 border-t border-slate-200">
                                                            <p className="text-xs font-semibold text-slate-500">
                                                                Unlocked on {new Date(userAchievement.unlocked_at).toLocaleDateString('en-US', {
                                                                    month: 'short',
                                                                    day: 'numeric',
                                                                    year: 'numeric'
                                                                })}
                                                            </p>
                                                        </div>
                                                    )}
                                                </div>
                                            </div>
                                        )
                                    })}
                                </div>
                            </div>
                        )
                    })}
                </div>
            </div>
        </div>
    )
}
