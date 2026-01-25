import React, { useState, useEffect, useRef } from 'react';
import {
    Card,
    CardHeader,
    CardContent,
    TextField,
    MenuItem,
    Button,
    Box,
    Typography,
    Alert,
    AlertTitle,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Chip,
    Stack,
    CircularProgress,
    Divider,
    Paper,
} from '@mui/material';
import {
    ExpandMore as ExpandMoreIcon,
    CheckCircle as CheckCircleIcon,
    Cancel as CancelIcon,
    PlayArrow as PlayIcon,
} from '@mui/icons-material';
import Editor from '@monaco-editor/react';
import { adminProblemApi, adminCodeGenApi, adminSubmissionsApi } from '../../lib/api/admin';
import toast from 'react-hot-toast';
import { Refresh as RefreshIcon } from '@mui/icons-material';
import { filterEditorLanguage } from '../../utils/utils';

interface ReferenceSolutionValidatorProps {
    problemId: number;
    supportedLanguages: string[];
    onValidationSuccess?: () => void;
}

const ALL_LANGUAGES = [
    { value: 'python', label: 'Python' },
    { value: 'javascript', label: 'JavaScript' },
    { value: 'java', label: 'Java' },
    { value: 'cpp', label: 'C++' },
    { value: 'c++', label: 'C++' },
    { value: 'c', label: 'C' },
    { value: 'go', label: 'Go' },
];

export const ReferenceSolutionValidator: React.FC<ReferenceSolutionValidatorProps> = ({ problemId, supportedLanguages, onValidationSuccess }) => {
    const [code, setCode] = useState('');
    const [language, setLanguage] = useState(supportedLanguages[0] || 'python');
    const [validating, setValidating] = useState(false);
    const [fetchingStub, setFetchingStub] = useState(false);
    const [validationResult, setValidationResult] = useState<any>(null);
    const [validationStatus, setValidationStatus] = useState<any>(null);

    const [submissionId, setSubmissionId] = useState<number | null>(null);
    const pollingIntervalRef = useRef<any>(null);

    useEffect(() => {
        fetchValidationStatus();
        return () => {
            if (pollingIntervalRef.current) {
                clearInterval(pollingIntervalRef.current);
            }
        };
    }, [problemId]);

    useEffect(() => {
        if (supportedLanguages.length > 0 && !supportedLanguages.includes(language)) {
            setLanguage(supportedLanguages[0]);
        }
    }, [supportedLanguages]);

    useEffect(() => {
        handleFetchStub();
    }, [language, problemId]);

    // Polling effect
    useEffect(() => {
        if (submissionId) {
            setValidating(true);
            pollingIntervalRef.current = setInterval(checkSubmissionStatus, 3000);
        } else {
            if (pollingIntervalRef.current) {
                clearInterval(pollingIntervalRef.current);
                pollingIntervalRef.current = null;
            }
            setValidating(false);
        }

        return () => {
            if (pollingIntervalRef.current) {
                clearInterval(pollingIntervalRef.current);
            }
        };
    }, [submissionId]);

    const checkSubmissionStatus = async () => {
        if (!submissionId) return;

        try {
            const response = await adminSubmissionsApi.getById(problemId, submissionId);
            const data = response.data.data;
            const status = data.status;

            if (status !== 'Pending' && status !== 'Processing') {
                // Completed
                setSubmissionId(null);

                // Construct validation result from submission data
                const result = {
                    is_valid: status === 'Accepted',
                    passed_tests: data.passed_test_cases,
                    total_tests: data.total_test_cases,
                    test_results: data.test_case_results,
                    error_message: data.error_message
                };

                setValidationResult(result);

                if (status === 'Accepted') {
                    toast.success('Reference solution validated! All test cases passed.');
                    if (onValidationSuccess) {
                        onValidationSuccess();
                    }
                } else {
                    toast.error(`Validation failed: ${status}`);
                }

                fetchValidationStatus();
            }
        } catch (error) {
            console.error('Polling error:', error);
        }
    };

    const handleFetchStub = async () => {
        setFetchingStub(true);
        try {
            const response = await adminCodeGenApi.v2GetStub(problemId, language);
            setCode(response.data.data.stub_code);
        } catch (error) {
            console.error('Failed to fetch stub:', error);
            // Don't toast error here as it might be expected if no boilerplate exists yet
        } finally {
            setFetchingStub(false);
        }
    };

    const fetchValidationStatus = async () => {
        try {
            const response = await adminProblemApi.v2GetValidationStatus(problemId);
            setValidationStatus(response.data.data);
        } catch (error) {
            console.error('Failed to fetch validation status:', error);
        }
    };

    const handleValidate = async () => {
        if (!code.trim()) {
            toast.error('Please write reference solution code');
            return;
        }

        setValidating(true);
        setValidationResult(null);
        try {
            const response = await adminProblemApi.v2Validate(problemId, {
                language_slug: language,
                code: code,
            });

            const data = response.data.data;
            if (data.submission_id) {
                setSubmissionId(data.submission_id);
                toast.success('Validation started...');
            } else {
                toast.error('Failed to start validation process');
                setValidating(false);
            }
        } catch (error: any) {
            const message = error.response?.data?.data?.message || 'Failed to validate reference solution';
            toast.error(message);
            console.error(error);
            setValidating(false);
        }
    };

    return (
        <Card variant="outlined" sx={{ mt: 3 }}>
            <CardHeader
                title="Reference Solution Validation"
                subheader="Ensure your problem is solvable and test cases are correct."
            />
            <Divider />
            <CardContent>
                {validationStatus && (
                    <Alert
                        severity={validationStatus.can_publish ? "success" : "warning"}
                        sx={{ mb: 4 }}
                        icon={validationStatus.can_publish ? <CheckCircleIcon /> : undefined}
                    >
                        <AlertTitle>
                            Status: {validationStatus.validation_status.toUpperCase()}
                        </AlertTitle>
                        {validationStatus.can_publish ? (
                            <Box>
                                <Typography variant="body2">
                                    Problem is validated and ready to publish.
                                </Typography>
                                {validationStatus.validated_languages.length > 0 && (
                                    <Stack direction="row" spacing={1} sx={{ mt: 1 }}>
                                        <Typography variant="caption" sx={{ fontWeight: 'bold' }}>Validated for:</Typography>
                                        {validationStatus.validated_languages.map((lang: string) => (
                                            <Chip key={lang} label={lang} size="small" color="success" variant="outlined" />
                                        ))}
                                    </Stack>
                                )}
                            </Box>
                        ) : (
                            <Typography variant="body2">
                                Problem must be validated before publishing. Submit and validate a reference solution below.
                            </Typography>
                        )}
                    </Alert>
                )}

                <Box sx={{ mb: 3 }}>
                    <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 2 }}>
                        <TextField
                            select
                            label="Select Language"
                            value={supportedLanguages.includes(language) ? language : (supportedLanguages[0] || '')}
                            onChange={(e) => setLanguage(e.target.value)}
                            sx={{ width: 200 }}
                            size="small"
                            disabled={supportedLanguages.length === 0}
                        >
                            {ALL_LANGUAGES.filter(lang => supportedLanguages.includes(lang.value)).map((option) => (
                                <MenuItem key={option.value} value={option.value}>
                                    {option.label}
                                </MenuItem>
                            ))}
                        </TextField>

                        <Button
                            startIcon={fetchingStub ? <CircularProgress size={16} color="inherit" /> : <RefreshIcon />}
                            onClick={handleFetchStub}
                            disabled={fetchingStub}
                            size="small"
                            variant="outlined"
                        >
                            {fetchingStub ? 'Loading...' : 'Load Stub'}
                        </Button>
                    </Stack>

                    <Box sx={{ border: '1px solid #ccc', borderRadius: 1, overflow: 'hidden', mb: 2 }}>
                        <Editor
                            height="400px"
                            language={filterEditorLanguage(language)}
                            theme="vs-dark"
                            value={code}
                            onChange={(value) => setCode(value || '')}
                            options={{
                                selectOnLineNumbers: true,
                                minimap: { enabled: false },
                                fontSize: 14,
                                scrollBeyondLastLine: false,
                            }}
                        />
                    </Box>

                    <Button
                        variant="contained"
                        color="primary"
                        size="large"
                        onClick={handleValidate}
                        disabled={validating}
                        startIcon={validating ? <CircularProgress size={20} color="inherit" /> : <PlayIcon />}
                    >
                        {validating ? 'Validating...' : 'Validate Reference Solution'}
                    </Button>
                </Box>

                {validationResult && (
                    <Box sx={{ mt: 4 }}>
                        <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            Validation Results:
                            {validationResult.is_valid ? (
                                <Chip icon={<CheckCircleIcon />} label="PASSED" color="success" />
                            ) : (
                                <Chip icon={<CancelIcon />} label="FAILED" color="error" />
                            )}
                        </Typography>

                        <Typography variant="body2" sx={{ mb: 2, fontWeight: 'medium' }}>
                            Passed: {validationResult.passed_tests} / {validationResult.total_tests} test cases
                        </Typography>

                        {validationResult.error_message && (
                            <Alert severity="error" sx={{ mb: 2 }}>
                                <AlertTitle>Execution Error</AlertTitle>
                                <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{validationResult.error_message}</pre>
                            </Alert>
                        )}

                        <Box>
                            {validationResult.test_results?.map((test: any, idx: number) => (
                                <Accordion key={idx} variant="outlined" sx={{ mb: 1 }}>
                                    <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                        <Stack direction="row" spacing={2} alignItems="center">
                                            <Typography variant="subtitle2">Test Case {idx + 1}</Typography>
                                            <Chip
                                                label={test.status}
                                                size="small"
                                                color={test.status === 'Passed' ? 'success' : 'error'}
                                                variant="outlined"
                                            />
                                            {test.is_sample && <Chip label="Sample" size="small" color="primary" variant="outlined" />}
                                        </Stack>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
                                            <Box>
                                                <Typography variant="caption" color="textSecondary" sx={{ fontWeight: 'bold' }}>Input</Typography>
                                                <Paper variant="outlined" sx={{ p: 1, bgcolor: '#f5f5f5', overflowX: 'auto' }}>
                                                    <pre style={{ margin: 0, fontSize: '0.75rem' }}>{test.input}</pre>
                                                </Paper>
                                            </Box>
                                            <Box>
                                                <Typography variant="caption" color="textSecondary" sx={{ fontWeight: 'bold' }}>Expected Output</Typography>
                                                <Paper variant="outlined" sx={{ p: 1, bgcolor: '#f5f5f5', overflowX: 'auto' }}>
                                                    <pre style={{ margin: 0, fontSize: '0.75rem' }}>{test.expected_output}</pre>
                                                </Paper>
                                            </Box>
                                        </Box>
                                        <Box sx={{ mt: 2 }}>
                                            <Typography variant="caption" color="textSecondary" sx={{ fontWeight: 'bold' }}>Actual Output</Typography>
                                            <Paper
                                                variant="outlined"
                                                sx={{
                                                    p: 1,
                                                    bgcolor: test.status === 'Passed' ? '#f6ffed' : '#fff1f0',
                                                    overflowX: 'auto'
                                                }}
                                            >
                                                <pre style={{ margin: 0, fontSize: '0.75rem' }}>{test.actual_output}</pre>
                                            </Paper>
                                        </Box>
                                        {test.error_message && (
                                            <Alert severity="error" sx={{ mt: 2 }}>
                                                {test.error_message}
                                            </Alert>
                                        )}
                                    </AccordionDetails>
                                </Accordion>
                            ))}
                        </Box>
                    </Box>
                )}
            </CardContent>
        </Card>
    );
};

export default ReferenceSolutionValidator;
