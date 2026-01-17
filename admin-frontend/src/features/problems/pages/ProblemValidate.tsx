import { useParams, useNavigate } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    Card,
    Button,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Chip,
    Stack,
    CircularProgress,
    Alert,
    Typography,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Divider,
    Box,
    AlertTitle,
    CardContent,
    CardHeader
} from '@mui/material'
import {
    CheckCircle as CheckCircleIcon,
    Cancel as CancelIcon,
    Sync as SyncIcon,
    ArrowBack as ArrowBackIcon,
    RocketLaunch as RocketIcon,
    Visibility as EyeIcon,
    Info as InfoIcon
} from '@mui/icons-material'
import toast from 'react-hot-toast'
import Editor from '@monaco-editor/react'
import { adminProblemLanguagesApi, adminProblemApi, adminSubmissionsApi } from '../../../lib/api/admin'
import { ProblemStepper } from '../components/ProblemStepper'

export default function ProblemValidate() {
    const { problemId } = useParams<{ problemId: string }>()
    const navigate = useNavigate()
    const queryClient = useQueryClient()

    // Track polling for each language
    const [pollingStatus, setPollingStatus] = useState<Record<number, number | null>>({})
    // Store validation errors for the current session
    const [validationErrors, setValidationErrors] = useState<Record<number, string | null>>({})
    const [errorDetail, setErrorDetail] = useState<{ language: string, error: string } | null>(null)

    // Preview state
    const [previewCode, setPreviewCode] = useState<string | null>(null)
    const [previewLanguage, setPreviewLanguage] = useState<string | null>(null)
    const [isPreviewVisible, setIsPreviewVisible] = useState(false)
    const [isPreviewLoading, setIsPreviewLoading] = useState(false)

    // Fetch problem languages
    const { data: problemLanguagesData, isLoading: problemLanguagesLoading } = useQuery({
        queryKey: ['problem-languages', problemId],
        queryFn: () => adminProblemLanguagesApi.getAll(String(problemId)),
    })

    // Fetch problem details for publishing (optional, for validation check)
    const { data: _problemData } = useQuery({
        queryKey: ['problem', problemId],
        queryFn: () => adminProblemApi.getById(String(problemId)),
        select: (res) => res.data.data
    })

    const problemLanguages = problemLanguagesData?.data?.data || []

    // Validate mutation
    const validateMutation = useMutation({
        mutationFn: (languageId: number) => adminProblemLanguagesApi.validate(String(problemId), languageId),
        onSuccess: (res, languageId) => {
            toast.success('Validation started')
            const submissionId = res.data.data.id
            setPollingStatus(prev => ({ ...prev, [languageId]: submissionId }))
        },
        onError: () => {
            toast.error('Failed to start validation')
        },
    })

    // Publish mutation
    const publishMutation = useMutation({
        mutationFn: () => adminProblemApi.publish(String(problemId)),
        onSuccess: () => {
            toast.success('Problem published successfully!')
            navigate('/problems')
        },
        onError: () => {
            toast.error('Failed to publish problem')
        },
    })

    // Polling logic
    useEffect(() => {
        const activePolling = Object.entries(pollingStatus).filter(([_, subId]) => subId !== null)

        if (activePolling.length === 0) return

        const pollInterval = setInterval(async () => {
            for (const [langId, subId] of activePolling) {
                try {
                    const response = await adminSubmissionsApi.getById(Number(problemId), Number(subId))
                    const data = response.data.data
                    const status = data.status

                    if (status !== 'Pending') {
                        // Finished
                        setPollingStatus(prev => ({ ...prev, [langId]: null }))
                        queryClient.invalidateQueries({ queryKey: ['problem-languages', problemId] })

                        if (status === 'Accepted') {
                            toast.success(`Validation passed! (${data.passed_test_cases}/${data.total_test_cases} passed)`)
                            setValidationErrors(prev => ({ ...prev, [langId]: null }))
                        } else {
                            const errorMsg = data.error_message || status
                            const passRate = data.total_test_cases > 0 ? ` (${data.passed_test_cases}/${data.total_test_cases} passed)` : ''
                            setValidationErrors(prev => ({ ...prev, [langId]: errorMsg }))
                            toast.error(`Validation failed: ${status}${passRate}`)
                        }
                    }
                } catch (error) {
                    console.error('Polling error:', error)
                }
            }
        }, 2000)

        return () => clearInterval(pollInterval)
    }, [pollingStatus, problemId, queryClient])

    const handleValidate = (languageId: number) => {
        validateMutation.mutate(languageId)
    }

    const handlePreview = async (languageId: number, languageName: string) => {
        setIsPreviewLoading(true)
        setPreviewLanguage(languageName)
        try {
            const res = await adminProblemLanguagesApi.preview(String(problemId), languageId)
            setPreviewCode(res.data.data.combined_code)
            setIsPreviewVisible(true)
        } catch (error) {
            toast.error('Failed to fetch code preview')
        } finally {
            setIsPreviewLoading(false)
        }
    }

    const allValidated = problemLanguages.length > 0 && problemLanguages.every(pl => pl.is_validated)

    if (problemLanguagesLoading) {
        return (
            <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
                <CircularProgress size={60} />
            </Box>
        )
    }

    return (
        <Stack spacing={4} p={4}>
            {/* Header */}
            <Box textAlign="center">
                <Typography variant="h4" fontWeight="bold" gutterBottom>
                    Validate & Publish
                </Typography>
                <ProblemStepper currentStep={4} model="publish" problemId={problemId || "create"} />
            </Box>

            <Card variant="outlined">
                <CardHeader title="Language Validation" />
                <Divider />
                <CardContent>
                    <Alert severity="info" sx={{ mb: 4 }}>
                        <AlertTitle>Final Step</AlertTitle>
                        Each language must be validated by running its solution code against all test cases. Once all languages are validated, you can publish the problem.
                    </Alert>

                    <TableContainer component={Paper} variant="outlined">
                        <Table>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Language</TableCell>
                                    <TableCell>Status</TableCell>
                                    <TableCell>Test Cases</TableCell>
                                    <TableCell>Last Validated</TableCell>
                                    <TableCell>Actions</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {problemLanguages.map((record: any) => {
                                    const isPolling = pollingStatus[record.language_id] !== null && pollingStatus[record.language_id] !== undefined

                                    return (
                                        <TableRow key={record.language_id}>
                                            <TableCell>
                                                <Typography variant="body2" fontWeight="bold">{record.language_name}</Typography>
                                                <Typography variant="caption" color="text.secondary">{record.language_version}</Typography>
                                            </TableCell>
                                            <TableCell>
                                                {isPolling ? (
                                                    <Chip icon={<SyncIcon sx={{ animation: 'spin 2s linear infinite' }} />} label="Validating..." color="info" variant="outlined" size="small" />
                                                ) : record.is_validated ? (
                                                    <Chip icon={<CheckCircleIcon />} label="Validated" color="success" size="small" />
                                                ) : (validationErrors[record.language_id] || (record.last_validation_status && record.last_validation_status !== 'Accepted')) ? (
                                                    <Chip icon={<CancelIcon />} label={validationErrors[record.language_id] ? 'Failed' : record.last_validation_status || 'Validation Failed'} color="error" size="small" />
                                                ) : (
                                                    <Chip label="Not Validated" size="small" />
                                                )}
                                            </TableCell>
                                            <TableCell>
                                                {isPolling ? (
                                                    <Typography variant="body2" color="text.secondary">Evaluating...</Typography>
                                                ) : (record.is_validated || record.last_validation_status) ? (
                                                    <Typography
                                                        variant="body2"
                                                        color={(record.last_pass_count === record.last_total_count) ? 'success.main' : 'error.main'}
                                                    >
                                                        {record.last_pass_count || 0} / {record.last_total_count || 0} passed
                                                    </Typography>
                                                ) : '-'}
                                            </TableCell>
                                            <TableCell>
                                                {record.is_validated ? new Date(record.updated_at).toLocaleString() : '-'}
                                            </TableCell>
                                            <TableCell>
                                                <Stack direction="row" spacing={1}>
                                                    <Button
                                                        variant="outlined"
                                                        size="small"
                                                        startIcon={<EyeIcon />}
                                                        onClick={() => handlePreview(record.language_id, record.language_name)}
                                                        disabled={isPreviewLoading && previewLanguage === record.language_name}
                                                    >
                                                        Preview
                                                    </Button>
                                                    <Button
                                                        variant="contained"
                                                        size="small"
                                                        onClick={() => handleValidate(record.language_id)}
                                                        disabled={pollingStatus[record.language_id] !== null && pollingStatus[record.language_id] !== undefined}
                                                    >
                                                        Validate
                                                    </Button>
                                                    {(validationErrors[record.language_id] || record.last_validation_error) && (
                                                        <Button
                                                            color="error"
                                                            size="small"
                                                            startIcon={<InfoIcon />}
                                                            onClick={() => setErrorDetail({
                                                                language: record.language_name,
                                                                error: validationErrors[record.language_id] || record.last_validation_error || 'Unknown error'
                                                            })}
                                                        >
                                                            Error
                                                        </Button>
                                                    )}
                                                </Stack>
                                            </TableCell>
                                        </TableRow>
                                    )
                                })}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>

            {/* Navigation & Publish */}
            <Box display="flex" justifyContent="space-between" alignItems="center">
                <Button
                    variant="outlined"
                    startIcon={<ArrowBackIcon />}
                    onClick={() => navigate(`/problems/${problemId}/languages`)}
                >
                    Back to Languages
                </Button>

                <Stack direction="row" spacing={2} alignItems="center">
                    {!allValidated && problemLanguages.length > 0 && (
                        <Typography variant="body2" color="warning.main">
                            Validation required for all languages before publishing
                        </Typography>
                    )}
                    <Button
                        variant="contained"
                        size="large"
                        startIcon={<RocketIcon />}
                        onClick={() => publishMutation.mutate()}
                        disabled={!allValidated || publishMutation.isPending}
                        color={allValidated ? "success" : "primary"}
                    >
                        Publish Problem
                    </Button>
                </Stack>
            </Box>

            {/* Preview Dialog */}
            <Dialog
                open={isPreviewVisible}
                onClose={() => setIsPreviewVisible(false)}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle>Code Preview - {previewLanguage}</DialogTitle>
                <DialogContent dividers>
                    <Box sx={{ height: 500, border: 1, borderColor: 'divider', borderRadius: 1, overflow: 'hidden' }}>
                        <Editor
                            height="100%"
                            language={previewLanguage?.toLowerCase() || 'text'}
                            value={previewCode || ''}
                            options={{
                                readOnly: true,
                                minimap: { enabled: false },
                                fontSize: 14,
                                scrollBeyondLastLine: false,
                            }}
                        />
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setIsPreviewVisible(false)}>Close</Button>
                </DialogActions>
            </Dialog>

            {/* Error Detail Dialog */}
            <Dialog
                open={!!errorDetail}
                onClose={() => setErrorDetail(null)}
                maxWidth="md"
                fullWidth
            >
                <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <CancelIcon color="error" />
                    Validation Error - {errorDetail?.language}
                </DialogTitle>
                <DialogContent dividers>
                    <Paper
                        variant="outlined"
                        sx={{ p: 2, bgcolor: 'grey.50', maxHeight: 400, overflow: 'auto' }}
                    >
                        <Typography
                            component="pre"
                            variant="body2"
                            sx={{ color: 'error.main', fontFamily: 'monospace', whiteSpace: 'pre-wrap', m: 0 }}
                        >
                            {errorDetail?.error}
                        </Typography>
                    </Paper>
                    <Alert severity="error" sx={{ mt: 2 }}>
                        <AlertTitle>Validation Failed</AlertTitle>
                        The solution code failed against one or more test cases. Please review the solution code and try again.
                    </Alert>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setErrorDetail(null)}>Close</Button>
                </DialogActions>
            </Dialog>
        </Stack>
    )
}
