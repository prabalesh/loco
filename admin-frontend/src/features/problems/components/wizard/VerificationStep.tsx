import React, { useState } from 'react';
import {
    Stack,
    Typography,
    Box,
    Button,
    Paper,
    Alert,
    CircularProgress,
    TextField,
    MenuItem
} from '@mui/material';
import {
    PlayArrow as RunIcon,
    CheckCircle as SuccessIcon,
    Error as ErrorIcon
} from '@mui/icons-material';
import { adminProblemApi } from '../../../../lib/api/admin';
import toast from 'react-hot-toast';

interface VerificationStepProps {
    data: any;
    problemId: number | null;
    onProblemCreated: (id: number) => void;
}

export const VerificationStep: React.FC<VerificationStepProps> = ({ data, problemId, onProblemCreated }) => {
    const [verifying, setVerifying] = useState(false);
    const [creating, setCreating] = useState(false);
    const [judgeSolution, setJudgeSolution] = useState('');
    const [selectedLang, setSelectedLang] = useState(data.selected_languages[0] || 'python');
    const [results, setResults] = useState<any>(null);

    const handleCreateDraft = async () => {
        setCreating(true);
        try {
            const formattedTestCases = data.test_cases.map((tc: any) => ({
                ...tc,
                input: JSON.parse(tc.input),
                expected_output: JSON.parse(tc.expected_output)
            }));

            const problemData = {
                title: data.title,
                description: data.description,
                difficulty: data.difficulty,
                constraints: data.constraints,
                hints: data.hints,
                function_name: data.function_name,
                return_type: data.return_type,
                parameters: data.parameters,
                test_cases: formattedTestCases,
                validation_type: data.validation_type,
                status: 'draft',
                visibility: 'public',
                is_active: true,
                tag_ids: [],
                category_ids: []
            };

            const response = await adminProblemApi.v2Create(problemData);
            onProblemCreated(response.data.data.id);
            toast.success('Problem Draft Created!');
        } catch (err: any) {
            toast.error(err.message || 'Failed to create problem');
        } finally {
            setCreating(false);
        }
    };

    const handleVerify = async () => {
        if (!problemId) return;
        setVerifying(true);
        try {
            const resp = await adminProblemApi.v2Validate(problemId, {
                language_slug: selectedLang,
                code: judgeSolution
            });
            setResults(resp.data.data);
            if (resp.data.data.status === 'Accepted') {
                toast.success('Verification Successful! All test cases passed.');
            } else {
                toast.error(`Verification Failed: ${resp.data.data.status}`);
            }
        } catch (err: any) {
            toast.error(err.message || 'Validation failed');
        } finally {
            setVerifying(false);
        }
    };

    return (
        <Stack spacing={4}>
            <Box>
                <Typography variant="h6" color="primary" fontWeight="bold">Verification & Finalization</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                    Before publishing, you must verify the problem by providing a solution that passes all test cases.
                </Typography>

                {!problemId ? (
                    <Paper sx={{ p: 4, textAlign: 'center', bgcolor: '#f0f7ff', borderRadius: 2 }}>
                        <Typography variant="h6" gutterBottom>Final Step: Create Problem Draft</Typography>
                        <Typography variant="body2" sx={{ mb: 3 }}>
                            We'll save your configuration as a draft, allowing you to run the judge against your tests.
                        </Typography>
                        <Button
                            variant="contained"
                            size="large"
                            onClick={handleCreateDraft}
                            disabled={creating}
                            startIcon={creating ? <CircularProgress size={20} color="inherit" /> : null}
                            sx={{ px: 4 }}
                        >
                            {creating ? 'Initializing...' : 'Initialize Problem Draft'}
                        </Button>
                    </Paper>
                ) : (
                    <Stack spacing={4}>
                        <Alert severity="success" sx={{ borderRadius: 2 }}>Problem Draft Initialized (ID: {problemId})</Alert>

                        <Box>
                            <Typography variant="subtitle2" gutterBottom fontWeight="bold">Live Verification</Typography>
                            <Paper variant="outlined" sx={{ p: 3, borderRadius: 2 }}>
                                <Stack spacing={3}>
                                    <TextField
                                        select
                                        label="Select Language for Solution"
                                        size="small"
                                        value={data.selected_languages.includes(selectedLang) ? selectedLang : (data.selected_languages[0] || '')}
                                        onChange={(e) => setSelectedLang(e.target.value)}
                                        sx={{ width: 250 }}
                                        disabled={data.selected_languages.length === 0}
                                        error={data.selected_languages.length === 0}
                                        helperText={data.selected_languages.length === 0 ? "No supported languages selected" : ""}
                                    >
                                        {data.selected_languages.map((slug: string) => (
                                            <MenuItem key={slug} value={slug}>
                                                {slug.charAt(0).toUpperCase() + slug.slice(1)}
                                            </MenuItem>
                                        ))}
                                    </TextField>

                                    <TextField
                                        fullWidth
                                        multiline
                                        rows={12}
                                        label="Judge Solution Code"
                                        value={judgeSolution}
                                        onChange={(e) => setJudgeSolution(e.target.value)}
                                        placeholder="Write a solution that should pass all test cases..."
                                        sx={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                                    />

                                    <Button
                                        variant="contained"
                                        startIcon={verifying ? <CircularProgress size={20} color="inherit" /> : <RunIcon />}
                                        onClick={handleVerify}
                                        disabled={verifying || !judgeSolution}
                                        sx={{ alignSelf: 'flex-start', px: 4 }}
                                    >
                                        Run Verification Test
                                    </Button>
                                </Stack>
                            </Paper>
                        </Box>

                        {results && (
                            <Box>
                                <Typography variant="subtitle2" gutterBottom fontWeight="bold">Test Results</Typography>
                                <Paper
                                    variant="outlined"
                                    sx={{
                                        p: 3,
                                        borderRadius: 2,
                                        borderLeft: `6px solid ${results.status === 'Accepted' ? '#4caf50' : '#f44336'}`
                                    }}
                                >
                                    <Stack direction="row" spacing={3} alignItems="center">
                                        {results.status === 'Accepted' ?
                                            <SuccessIcon sx={{ fontSize: 40, color: '#4caf50' }} /> :
                                            <ErrorIcon sx={{ fontSize: 40, color: '#f44336' }} />
                                        }
                                        <Box>
                                            <Typography variant="h5" fontWeight="bold">{results.status}</Typography>
                                            <Typography variant="body1" color="text.secondary">
                                                Passed {results.passed_tests} out of {results.total_tests} test cases
                                            </Typography>
                                        </Box>
                                    </Stack>
                                </Paper>
                            </Box>
                        )}
                    </Stack>
                )}
            </Box>
        </Stack>
    );
};
