import { useState, useRef, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { ChevronDown, Check } from 'lucide-react'
import type { Boilerplate } from '../types'

interface LanguageDropdownProps {
    boilerplates: Boilerplate[]
    selectedLang: number | null
    onLanguageChange: (langId: number) => void
}

export const LanguageDropdown = ({ boilerplates, selectedLang, onLanguageChange }: LanguageDropdownProps) => {
    const [isOpen, setIsOpen] = useState(false)
    const dropdownRef = useRef<HTMLDivElement>(null)

    const currentLang = boilerplates?.find((b: Boilerplate) => b.language_id === selectedLang)

    console.log(currentLang)

    // Close dropdown when clicking outside
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
                setIsOpen(false)
            }
        }

        document.addEventListener('mousedown', handleClickOutside)
        return () => document.removeEventListener('mousedown', handleClickOutside)
    }, [])

    const handleSelect = (langId: number) => {
        onLanguageChange(langId)
        setIsOpen(false)
    }

    if (!boilerplates || boilerplates.length === 0) {
        return (
            <div className="bg-gray-800/80 rounded-xl px-4 py-2.5 border border-white/5">
                <span className="text-xs text-gray-500 font-bold uppercase tracking-widest italic">
                    No languages configured
                </span>
            </div>
        )
    }

    return (
        <div ref={dropdownRef} className="relative">
            <motion.button
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                onClick={() => setIsOpen(!isOpen)}
                className="bg-gray-800/80 hover:bg-gray-700/80 rounded-xl px-4 py-2.5 border border-white/5 shadow-lg flex items-center gap-3 min-w-[180px] transition-all duration-200"
            >
                <div className="flex-1 text-left">
                    <div className="text-white text-sm font-bold">
                        {currentLang?.language.name || 'Select Language'}
                    </div>
                </div>
                <motion.div
                    animate={{ rotate: isOpen ? 180 : 0 }}
                    transition={{ duration: 0.2 }}
                >
                    <ChevronDown className="h-4 w-4 text-gray-400" />
                </motion.div>
            </motion.button>

            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0, y: -10, scale: 0.95 }}
                        animate={{ opacity: 1, y: 0, scale: 1 }}
                        exit={{ opacity: 0, y: -10, scale: 0.95 }}
                        transition={{ duration: 0.15 }}
                        className="absolute top-full mt-2 left-0 right-0 bg-gray-800 rounded-xl border border-white/10 shadow-2xl overflow-hidden z-50"
                    >
                        <div className="max-h-[300px] overflow-y-auto custom-scrollbar-dark">
                            {boilerplates.map((boilerplate: Boilerplate, index: number) => (
                                <motion.button
                                    key={boilerplate.language_id}
                                    initial={{ opacity: 0, x: -10 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    transition={{ delay: index * 0.03 }}
                                    onClick={() => handleSelect(boilerplate.language_id)}
                                    className={`w-full px-4 py-3 flex items-center justify-between hover:bg-gray-700/50 transition-all duration-150 ${selectedLang === boilerplate.language_id ? 'bg-blue-600/20' : ''
                                        }`}
                                >
                                    <div className="text-left flex-1">
                                        <div className={`text-sm font-bold ${selectedLang === boilerplate.language_id ? 'text-blue-400' : 'text-white'
                                            }`}>
                                            {boilerplate.language.name}
                                        </div>
                                        <div className="text-gray-500 text-[10px] font-medium uppercase tracking-wider mt-0.5">
                                            {boilerplate.language.version}
                                        </div>
                                    </div>
                                    {selectedLang === boilerplate.language_id && (
                                        <motion.div
                                            initial={{ scale: 0 }}
                                            animate={{ scale: 1 }}
                                            transition={{ type: "spring", stiffness: 500, damping: 25 }}
                                        >
                                            <Check className="h-5 w-5 text-blue-400" />
                                        </motion.div>
                                    )}
                                </motion.button>
                            ))}
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    )
}
