import { useNavigate, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import { XCircle, ChevronLeft } from 'lucide-react'
import { problemsApi } from '../api/problems'
import { submissionsApi, type RunCodeResult } from '../api/submissions'
import { Button } from '@/shared/components/ui/Button'
import type { Problem, ProblemLanguage, Submission } from '../types'
import { ProblemHeader } from '../components/ProblemHeader'
import { ProblemTabs } from '../components/ProblemTabs'
import { DescriptionTab } from '../components/DescriptionTab'
import { ResultTab } from '../components/ResultTab'
import { SubmissionsTab } from '../components/SubmissionsTab'
import { CodeEditor } from '../components/CodeEditor'
import { SubmissionDetailsModal } from '../components/SubmissionDetailsModal'
import { Skeleton } from '@/shared/components/ui/Skeleton'

const ProblemDetailSkeleton = () => (
    <div className="flex flex-col h-[calc(100vh-64px)] bg-gray-50">
        {/* Header Skeleton */}
        <div className="h-16 bg-white border-b border-gray-200 px-6 flex items-center justify-between">
            <div className="flex items-center gap-4">
                <Skeleton className="h-8 w-8 rounded-lg" />
                <Skeleton className="h-6 w-48" />
            </div>
            <div className="flex items-center gap-3">
                <Skeleton className="h-9 w-24 rounded-lg" />
                <Skeleton className="h-9 w-24 rounded-lg" />
            </div>
        </div>

        <main className="flex-1 flex overflow-hidden">
            {/* Left Panel Skeleton */}
            <div className="w-1/2 flex flex-col border-r border-gray-200 bg-white">
                <div className="flex border-b border-gray-100">
                    <Skeleton className="h-12 w-32 m-2" />
                    <Skeleton className="h-12 w-32 m-2" />
                    <Skeleton className="h-12 w-32 m-2" />
                </div>
                <div className="p-8 space-y-6">
                    <Skeleton className="h-10 w-3/4" />
                    <div className="space-y-3">
                        <Skeleton className="h-4 w-full" />
                        <Skeleton className="h-4 w-full" />
                        <Skeleton className="h-4 w-5/6" />
                    </div>
                    <Skeleton className="h-32 w-full rounded-xl" />
                </div>
            </div>

            {/* Right Panel Skeleton (Editor) */}
            <div className="flex-1 flex flex-col bg-[#1e1e1e]">
                <div className="h-11 border-b border-white/5 flex items-center justify-between px-4">
                    <Skeleton className="h-7 w-32 bg-white/10" />
                    <Skeleton className="h-7 w-20 bg-white/10" />
                </div>
                <div className="flex-1 p-4">
                    <Skeleton className="h-full w-full bg-white/5 rounded-md" />
                </div>
            </div>
        </main>
    </div>
)

export const ProblemDetailPage = () => {
    const { slug } = useParams<{ slug: string }>()
    const navigate = useNavigate()
    const queryClient = useQueryClient()
    const [activeTab, setActiveTab] = useState<'description' | 'result' | 'submissions'>('description')
    const [selectedLang, setSelectedLang] = useState<number | null>(null)
    const [code, setCode] = useState('')
    const [leftWidth, setLeftWidth] = useState(50) // Percentage
    const [isResizing, setIsResizing] = useState(false)
    const [pollingId, setPollingId] = useState<number | null>(null)
    const [runResult, setRunResult] = useState<RunCodeResult | null>(null)
    const [isRunning, setIsRunning] = useState(false)
    const [selectedSubmission, setSelectedSubmission] = useState<Submission | null>(null)
    const [isModalOpen, setIsModalOpen] = useState(false)

    // Fetch Problem
    const { data: problem, isLoading: isProblemLoading } = useQuery({
        queryKey: ['problem', slug],
        queryFn: () => problemsApi.get(slug!).then(res => (res.data as any).data),
        enabled: !!slug
    })

    // Fetch Languages
    const { data: languages } = useQuery({
        queryKey: ['problem-languages', problem?.id],
        queryFn: () => problemsApi.getLanguages(problem!.id).then(res => (res.data as any).data),
        enabled: !!problem?.id
    })

    // Fetch Sample Test Cases
    const { data: sampleTestCases } = useQuery({
        queryKey: ['sample-test-cases', problem?.id],
        queryFn: () => problemsApi.getSampleTestCases(problem!.id).then(res => (res.data as any).data),
        enabled: !!problem?.id
    })

    // Set default language and code (from localStorage or first available)
    useEffect(() => {
        if (languages && languages.length > 0) {
            if (!selectedLang) {
                const firstLang = languages[0]
                const savedCode = localStorage.getItem(`loco-code-${problem?.id}-${firstLang.language_id}`)
                setSelectedLang(firstLang.language_id)
                setCode(savedCode || firstLang.function_code || '')
            }
        }
    }, [languages, selectedLang, problem?.id])

    // Save code to localStorage
    useEffect(() => {
        if (problem?.id && selectedLang && code) {
            localStorage.setItem(`loco-code-${problem.id}-${selectedLang}`, code)
        }
    }, [code, selectedLang, problem?.id])

    const startResizing = (e: React.MouseEvent) => {
        setIsResizing(true)
        e.preventDefault()
    }

    const stopResizing = () => {
        setIsResizing(false)
    }

    const resize = (e: MouseEvent) => {
        if (isResizing) {
            const newWidth = (e.clientX / window.innerWidth) * 100
            if (newWidth > 20 && newWidth < 80) {
                setLeftWidth(newWidth)
            }
        }
    }

    useEffect(() => {
        if (isResizing) {
            window.addEventListener('mousemove', resize)
            window.addEventListener('mouseup', stopResizing)
        } else {
            window.removeEventListener('mousemove', resize)
            window.removeEventListener('mouseup', stopResizing)
        }
        return () => {
            window.removeEventListener('mousemove', resize)
            window.removeEventListener('mouseup', stopResizing)
        }
    }, [isResizing])

    const currentLang = languages?.find((l: ProblemLanguage) => l.language_id === selectedLang)

    // Run Code Mutation (doesn't create submission)
    const runCodeMutation = useMutation({
        mutationFn: ({ pId, lId, code }: { pId: number, lId: number, code: string }) =>
            submissionsApi.runCode(pId, lId, code).then(res => (res.data as any).data),
        onSuccess: (data: RunCodeResult) => {
            setRunResult(data)
            setActiveTab('result')
            setIsRunning(false)
            if (data.status === 'Accepted') {
                toast.success('All test cases passed!', { id: 'run-result' })
            } else {
                toast.error(`Failed: ${data.status}`, { id: 'run-result' })
            }
        },
        onError: (err: any) => {
            setIsRunning(false)
            toast.error(err.response?.data?.message || 'Run failed')
        }
    })

    // Submission Mutation
    const submitMutation = useMutation({
        mutationFn: ({ pId, lId, code }: { pId: number, lId: number, code: string }) =>
            submissionsApi.submit(pId, lId, code).then(res => (res.data as any).data),
        onSuccess: (data: Submission) => {
            setPollingId(data.id)
            setRunResult(null) // Clear run result when submitting
            setActiveTab('result')
            toast.loading('Evaluating...', { id: 'evaluating' })
        },
        onError: (err: any) => {
            toast.error(err.response?.data?.message || 'Submission failed')
        }
    })

    // Polling for submission status
    const { data: submissionResult } = useQuery({
        queryKey: ['submission', pollingId],
        queryFn: () => submissionsApi.get(pollingId!).then(res => (res.data as any).data),
        enabled: !!pollingId,
        refetchInterval: (query) => {
            const status = query.state.data?.status
            if (status && status !== 'Pending') {
                if (status === 'Accepted') {
                    toast.success('Accepted!', { id: 'evaluating' })
                } else {
                    toast.error(`Failed: ${status}`, { id: 'evaluating' })
                }
                setPollingId(null)
                queryClient.invalidateQueries({ queryKey: ['user-submissions'] })
                return false
            }
            return 3000 // Increased from 1000ms to 3000ms
        }
    })

    // Submissions History
    const { data: submissionsHistory } = useQuery({
        queryKey: ['user-submissions', problem?.id],
        queryFn: () => submissionsApi.list(problem!.id, 1, 10).then(res => (res.data as any).data.data),
        enabled: !!problem?.id && activeTab === 'submissions'
    })

    const handleLanguageChange = (langId: number) => {
        const lang = languages?.find((l: ProblemLanguage) => l.language_id === langId)
        if (lang) {
            setSelectedLang(langId)
            const savedCode = localStorage.getItem(`loco-code-${problem?.id}-${langId}`)
            setCode(savedCode || lang.function_code || '')
            toast.success(`Switched to ${lang.language_name}`, { duration: 2000 })
        }
    }

    const handleRun = () => {
        if (!problem || !selectedLang) return
        setIsRunning(true)
        setRunResult(null)
        runCodeMutation.mutate({ pId: (problem as Problem).id, lId: selectedLang, code })
    }

    const handleSubmit = () => {
        if (!problem || !selectedLang) return
        submitMutation.mutate({ pId: (problem as Problem).id, lId: selectedLang, code })
    }

    const handleResetCode = () => {
        setCode(currentLang?.function_code || '')
        toast.success('Code reset to default', { duration: 2000 })
    }

    const handleViewSubmission = (submission: Submission) => {
        setSelectedSubmission(submission)
        setIsModalOpen(true)
    }

    const handleCloseModal = () => {
        setIsModalOpen(false)
        setTimeout(() => setSelectedSubmission(null), 300)
    }

    if (isProblemLoading) {
        return <ProblemDetailSkeleton />
    }

    if (!problem) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="text-center">
                    <XCircle className="h-16 w-16 text-rose-500 mx-auto mb-4" />
                    <h2 className="text-2xl font-bold text-gray-900 mb-2">Problem not found</h2>
                    <p className="text-gray-600 mb-6">The problem you're looking for doesn't exist.</p>
                    <Button onClick={() => navigate('/problems')}>
                        <ChevronLeft className="h-4 w-4 mr-2" />
                        Back to Problems
                    </Button>
                </div>
            </div>
        )
    }

    return (
        <div className="flex flex-col h-[calc(100vh-64px)] bg-gray-50">
            <ProblemHeader
                problem={problem}
                onBack={() => navigate('/problems')}
                onRun={handleRun}
                onSubmit={handleSubmit}
                isSubmitting={submitMutation.isPending || runCodeMutation.isPending || !!pollingId || isRunning}
            />

            <main className="flex-1 flex overflow-hidden relative">
                {/* Left Section: Context */}
                <section
                    className="flex flex-col border-r border-gray-200 bg-white shadow-sm overflow-hidden"
                    style={{ width: `${leftWidth}%` }}
                >
                    <ProblemTabs activeTab={activeTab} onTabChange={setActiveTab} />

                    <div className="flex-1 overflow-y-auto p-8 custom-scrollbar">
                        {activeTab === 'description' && <DescriptionTab problem={problem} />}
                        {activeTab === 'result' && (
                            <ResultTab
                                submissionResult={submissionResult}
                                pollingId={pollingId}
                                runResult={runResult}
                                isRunning={isRunning || runCodeMutation.isPending}
                                sampleTestCases={sampleTestCases}
                            />
                        )}
                        {activeTab === 'submissions' && (
                            <SubmissionsTab
                                submissionsHistory={submissionsHistory}
                                problemId={problem.id}
                                onViewSubmission={handleViewSubmission}
                            />
                        )}
                    </div>
                </section>

                {/* Resize Handle */}
                <div
                    onMouseDown={startResizing}
                    className={`w-1 hover:w-1.5 transition-all cursor-col-resize bg-gray-200 hover:bg-blue-400 z-50 flex items-center justify-center ${isResizing ? 'bg-blue-500 w-1.5' : ''
                        }`}
                />

                {/* Right Section: Code Editor */}
                <div className="flex-1 flex flex-col min-w-0">
                    <CodeEditor
                        languages={languages}
                        selectedLang={selectedLang}
                        currentLang={currentLang}
                        code={code}
                        onLanguageChange={handleLanguageChange}
                        onCodeChange={setCode}
                        onResetCode={handleResetCode}
                    />
                </div>
            </main>

            {/* Submission Details Modal */}
            <SubmissionDetailsModal
                submission={selectedSubmission}
                isOpen={isModalOpen}
                onClose={handleCloseModal}
            />
        </div>
    )
}