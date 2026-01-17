import { motion } from 'framer-motion'
import { FileText, Trophy, History } from 'lucide-react'

interface ProblemTabsProps {
    activeTab: 'description' | 'result' | 'submissions'
    onTabChange: (tab: 'description' | 'result' | 'submissions') => void
}

export const ProblemTabs = ({ activeTab, onTabChange }: ProblemTabsProps) => {
    return (
        <nav className="flex border-b border-gray-200 px-2 bg-gradient-to-b from-white to-gray-50">
            <TabButton
                active={activeTab === 'description'}
                onClick={() => onTabChange('description')}
                icon={<FileText className="h-4 w-4" />}
                label="Description"
            />
            <TabButton
                active={activeTab === 'result'}
                onClick={() => onTabChange('result')}
                icon={<Trophy className="h-4 w-4" />}
                label="Result"
            />
            <TabButton
                active={activeTab === 'submissions'}
                onClick={() => onTabChange('submissions')}
                icon={<History className="h-4 w-4" />}
                label="Submissions"
            />
        </nav>
    )
}

const TabButton = ({ 
    active, 
    onClick, 
    icon, 
    label 
}: { 
    active: boolean
    onClick: () => void
    icon: React.ReactNode
    label: string 
}) => (
    <button
        onClick={onClick}
        className={`flex items-center gap-2 px-6 py-3.5 text-xs font-bold uppercase tracking-wider transition-all duration-200 relative ${
            active 
                ? 'text-blue-600 bg-gradient-to-b from-blue-50 to-white' 
                : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50/80'
        }`}
    >
        <motion.span 
            className="transition-all duration-200"
            animate={{ scale: active ? 1.1 : 1, opacity: active ? 1 : 0.7 }}
        >
            {icon}
        </motion.span>
        {label}
        {active && (
            <motion.div
                layoutId="activeTabUnderline"
                className="absolute bottom-0 left-0 right-0 h-0.5 bg-gradient-to-r from-blue-400 via-blue-600 to-blue-400 shadow-lg shadow-blue-500/50"
                transition={{ type: "spring", stiffness: 380, damping: 30 }}
            />
        )}
    </button>
)
