import { useState } from 'react'
import Editor from '@monaco-editor/react'
import {
    Button,
    Card,
    CardContent,
    CardHeader,
    Alert,
    AlertTitle,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Box,
    Typography,
    Divider,
    Stack,
    Paper
} from '@mui/material'
import { Send as SendIcon, Refresh as RefreshIcon } from '@mui/icons-material'
import toast from 'react-hot-toast'

interface Language {
    id: number
    name: string
    language_id: string
    extension: string
}

interface SubmissionResult {
    id: number
    status: string
    runtime?: number
    memory?: number
    error_message?: string
    created_at: string
}

interface ProblemCodeTesterProps {
    problemId: number
    languages: Language[]
    onSubmit?: (result: SubmissionResult) => void
}

export const ProblemCodeTester = ({ problemId, languages, onSubmit }: ProblemCodeTesterProps) => {
    const [code, setCode] = useState('')
    const [selectedLanguage, setSelectedLanguage] = useState<number | null>(null)
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [result, setResult] = useState<SubmissionResult | null>(null)
    const [submissions, setSubmissions] = useState<SubmissionResult[]>([])

    const selectedLang = languages.find(l => l.id === selectedLanguage)

    const handleSubmit = async () => {
        if (!selectedLanguage || !code.trim()) {
            toast.error('Please select a language and write some code')
            return
        }

        setIsSubmitting(true)
        setResult(null)

        try {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/admin/problems/${problemId}/submit`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('admin_access_token')}`
                },
                body: JSON.stringify({
                    language_id: selectedLanguage,
                    code: code
                })
            })

            if (!response.ok) {
                throw new Error('Submission failed')
            }

            const data = await response.json()
            setResult(data)

            if (onSubmit) {
                onSubmit(data)
            }

            // Poll for result
            pollSubmissionResult(data.id)
            toast.success('Code submitted successfully!')
        } catch (error) {
            toast.error('Failed to submit code')
            console.error(error)
        } finally {
            setIsSubmitting(false)
        }
    }

    const pollSubmissionResult = async (submissionId: number) => {
        const maxAttempts = 20
        let attempts = 0

        const poll = setInterval(async () => {
            try {
                const response = await fetch(`${import.meta.env.VITE_API_URL}/submissions/${submissionId}`, {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('admin_access_token')}`
                    }
                })

                if (response.ok) {
                    const data = await response.json()
                    if (data.status !== 'Pending') {
                        setResult(data)
                        clearInterval(poll)

                        if (data.status === 'Accepted') {
                            toast.success('All test cases passed!')
                        } else {
                            toast.error(`Submission ${data.status}`)
                        }
                    }
                }

                attempts++
                if (attempts >= maxAttempts) {
                    clearInterval(poll)
                    toast.error('Submission timed out')
                }
            } catch (error) {
                clearInterval(poll)
                console.error('Polling error:', error)
            }
        }, 1000)
    }

    const fetchSubmissions = async () => {
        try {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/admin/problems/${problemId}/submissions?limit=10`,
                {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('admin_access_token')}`
                    }
                }
            )

            if (response.ok) {
                const data = await response.json()
                setSubmissions(data.data || [])
            }
        } catch (error) {
            console.error('Failed to fetch submissions:', error)
        }
    }

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'Accepted': return 'success'
            case 'Pending': return 'info'
            case 'Wrong Answer': return 'warning'
            default: return 'error'
        }
    }

    return (
        <Stack spacing={4}>
            <Card variant="outlined">
                <CardHeader title="Test Solution" />
                <Divider />
                <CardContent>
                    <Stack spacing={3}>
                        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                            <FormControl fullWidth sx={{ maxWidth: 200 }}>
                                <InputLabel id="language-select-label">Select Language</InputLabel>
                                <Select
                                    labelId="language-select-label"
                                    value={selectedLanguage || ''}
                                    label="Select Language"
                                    onChange={(e) => setSelectedLanguage(Number(e.target.value))}
                                >
                                    {languages.map(lang => (
                                        <MenuItem key={lang.id} value={lang.id}>
                                            {lang.name}
                                        </MenuItem>
                                    ))}
                                </Select>
                            </FormControl>

                            <Box sx={{ flexGrow: 1 }} />

                            <Button
                                variant="contained"
                                startIcon={<SendIcon />}
                                onClick={handleSubmit}
                                disabled={isSubmitting || !selectedLanguage || !code.trim()}
                            >
                                Submit Code
                            </Button>
                        </Box>

                        <Box sx={{ border: 1, borderColor: 'divider', borderRadius: 1, overflow: 'hidden' }}>
                            <Editor
                                height="400px"
                                language={selectedLang?.language_id || 'python'}
                                value={code}
                                onChange={(value) => setCode(value || '')}
                                theme="vs-dark"
                                options={{
                                    minimap: { enabled: false },
                                    fontSize: 14,
                                    lineNumbers: 'on',
                                    scrollBeyondLastLine: false,
                                    automaticLayout: true,
                                }}
                            />
                        </Box>

                        {result && (
                            <Alert
                                severity={getStatusColor(result.status) as any}
                                variant="outlined"
                            >
                                <AlertTitle>Status: {result.status}</AlertTitle>
                                <Box sx={{ mt: 1 }}>
                                    {result.runtime && <Typography variant="body2">Runtime: {result.runtime}ms</Typography>}
                                    {result.memory && <Typography variant="body2">Memory: {result.memory}KB</Typography>}
                                    {result.error_message && (
                                        <Box
                                            component="pre"
                                            sx={{
                                                p: 1,
                                                bgcolor: 'background.default',
                                                borderRadius: 1,
                                                fontSize: '0.875rem',
                                                overflow: 'auto',
                                                mt: 1
                                            }}
                                        >
                                            {result.error_message}
                                        </Box>
                                    )}
                                </Box>
                            </Alert>
                        )}
                    </Stack>
                </CardContent>
            </Card>

            <Card variant="outlined">
                <CardHeader
                    title="Recent Submissions"
                    action={
                        <Button size="small" onClick={fetchSubmissions} startIcon={<RefreshIcon />}>
                            Refresh
                        </Button>
                    }
                />
                <Divider />
                <CardContent>
                    {submissions.length === 0 ? (
                        <Typography color="text.secondary" align="center" py={2}>
                            No submissions yet
                        </Typography>
                    ) : (
                        <Stack spacing={1}>
                            {submissions.map((sub) => (
                                <Paper
                                    key={sub.id}
                                    sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between', bgcolor: 'grey.50' }}
                                    variant="outlined"
                                >
                                    <Box>
                                        <Typography
                                            fontWeight="medium"
                                            color={`${getStatusColor(sub.status)}.main`}
                                            component="span"
                                        >
                                            {sub.status}
                                        </Typography>
                                        <Typography variant="caption" color="text.secondary" sx={{ ml: 1 }}>
                                            {new Date(sub.created_at).toLocaleString()}
                                        </Typography>
                                    </Box>
                                    {sub.runtime && (
                                        <Typography variant="body2" color="text.secondary">
                                            {sub.runtime}ms
                                        </Typography>
                                    )}
                                </Paper>
                            ))}
                        </Stack>
                    )}
                </CardContent>
            </Card>
        </Stack>
    )
}
