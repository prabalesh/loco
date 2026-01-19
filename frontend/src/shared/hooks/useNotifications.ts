import { useEffect, createElement } from 'react'
import { toast } from 'react-hot-toast'
import { useNavigate } from 'react-router-dom'
import { useAuth } from './useAuth'
import { CONFIG } from '../constants/config'
import { Trophy, ArrowRight } from 'lucide-react'

export function useNotifications() {
    const { user } = useAuth()
    const navigate = useNavigate()

    useEffect(() => {
        if (!user) return

        const streamUrl = `${CONFIG.API_BASE_URL}/notifications/stream`
        const eventSource = new EventSource(streamUrl, { withCredentials: true })

        eventSource.onopen = () => {
            console.log('SSE connection established to:', streamUrl)
        }

        eventSource.onmessage = (event) => {
            try {
                const notification = JSON.parse(event.data)

                if (notification.type === 'achievement_unlocked') {
                    const achievement = notification.data

                    toast.custom((t) => (
                        createElement('div', {
                            className: `${t.visible ? 'animate-enter' : 'animate-leave'} max-w-md w-full bg-white shadow-2xl rounded-2xl pointer-events-auto flex ring-1 ring-black ring-opacity-5 overflow-hidden cursor-pointer hover:scale-[1.02] transition-transform duration-200 border-2 border-purple-100`,
                            onClick: () => {
                                toast.dismiss(t.id)
                                navigate('/achievements')
                            }
                        },
                            createElement('div', { className: 'flex-1 w-0 p-5' },
                                createElement('div', { className: 'flex items-start' },
                                    createElement('div', { className: 'flex-shrink-0 pt-0.5' },
                                        createElement('div', { className: 'p-3 bg-gradient-to-br from-purple-500 to-pink-500 rounded-xl shadow-lg' },
                                            createElement(Trophy, { className: 'h-6 w-6 text-white' })
                                        )
                                    ),
                                    createElement('div', { className: 'ml-4 flex-1' },
                                        createElement('p', { className: 'text-sm font-black text-purple-600 uppercase tracking-wider mb-1' },
                                            "New Achievement!"
                                        ),
                                        createElement('p', { className: 'text-xl font-black text-slate-900 leading-tight' },
                                            achievement.name
                                        ),
                                        createElement('p', { className: 'mt-1 text-sm font-medium text-slate-600' },
                                            achievement.description
                                        ),
                                        createElement('div', { className: 'mt-3 flex items-center gap-1 text-xs font-bold text-amber-600 bg-amber-50 px-2 py-1 rounded-lg w-fit' },
                                            createElement('span', null, `+${achievement.xp_reward} XP`)
                                        )
                                    )
                                )
                            ),
                            createElement('div', { className: 'flex border-l border-slate-100' },
                                createElement('button', {
                                    className: 'w-full border border-transparent rounded-none rounded-r-lg p-4 flex items-center justify-center text-sm font-bold text-purple-600 hover:bg-slate-50 transition-colors',
                                },
                                    createElement(ArrowRight, { className: 'w-5 h-5' })
                                )
                            )
                        )
                    ), {
                        duration: 8000,
                        position: 'top-right'
                    })
                }
            } catch (err) {
                console.error('Failed to parse notification:', err)
            }
        }

        eventSource.onerror = (err) => {
            console.error('SSE connection error:', err)
            // Browser automatically attempts to reconnect for SSE
        }

        return () => {
            eventSource.close()
        }
    }, [user])
}
