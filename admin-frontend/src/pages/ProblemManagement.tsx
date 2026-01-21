import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Button,
    Stack,
    CircularProgress,
    IconButton,
    Divider,
    Alert,
    AlertTitle,
    Chip,
    Paper,
    Stepper,
    Step,
    StepLabel,
    Grid,
    List,
    ListItem,
    ListItemText,
} from '@mui/material';
import {
    ArrowBack as ArrowBackIcon,
    Publish as PublishIcon,
    Code as CodeIcon,
    Assignment as AssignmentIcon,
    Refresh as RefreshIcon,
} from '@mui/icons-material';
import { adminProblemApi } from '../lib/api/admin';
import { ReferenceSolutionValidator } from '../components/v2/ReferenceSolutionValidator';
import toast from 'react-hot-toast';

const STEPS = [
    'Problem Created',
    'Boilerplates Generated',
    'Solution Validated',
    'Published',
];

const ProblemManagement: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [problem, setProblem] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [publishing, setPublishing] = useState(false);
    const [generatingBoilerplates, setGeneratingBoilerplates] = useState(false);

    useEffect(() => {
        if (id) {
            fetchProblem();
        }
    }, [id]);

    const fetchProblem = async () => {
        try {
            const response = await adminProblemApi.v2GetById(id!);
            setProblem(response.data.data);
        } catch (error) {
            toast.error('Failed to load problem');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleRegenerateBoilerplates = async () => {
        setGeneratingBoilerplates(true);
        try {
            await adminProblemApi.v2RegenerateBoilerplates(id!);
            toast.success('Boilerplates regenerated successfully!');
            fetchProblem(); // Refresh stats
        } catch (error: any) {
            toast.error('Failed to regenerate boilerplates');
        } finally {
            setGeneratingBoilerplates(false);
        }
    };

    const handlePublish = async () => {
        setPublishing(true);
        try {
            await adminProblemApi.v2Publish(id!);
            toast.success('Problem published successfully!');
            fetchProblem(); // Refresh
        } catch (error: any) {
            const message = error.response?.data?.data?.message || 'Failed to publish problem';
            toast.error(message);
        } finally {
            setPublishing(false);
        }
    };

    const getActiveStep = () => {
        if (!problem) return 0;
        if (problem.status === 'published') return 3;
        if (problem.validation_status === 'validated') return 2;
        if (problem.boilerplates && problem.boilerplates.length > 0) return 1;
        return 0;
    };

    const getNextAction = () => {
        const step = getActiveStep();
        switch (step) {
            case 0: return "Generate boilerplates to proceed.";
            case 1: return "Validate reference solution to ensure correctness.";
            case 2: return "Ready to publish! Click the publish button.";
            case 3: return "Problem is live.";
            default: return "";
        }
    };

    if (loading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 10 }}>
                <CircularProgress />
            </Box>
        );
    }

    if (!problem) {
        return (
            <Box sx={{ p: 3 }}>
                <Alert severity="error">Problem not found</Alert>
                <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/problems')} sx={{ mt: 2 }}>
                    Back to Problems
                </Button>
            </Box>
        );
    }

    const activeStep = getActiveStep();

    return (
        <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
            {/* Header */}
            <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <Stack direction="row" spacing={2} alignItems="center">
                    <IconButton onClick={() => navigate('/problems')} size="small">
                        <ArrowBackIcon />
                    </IconButton>
                    <Box>
                        <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
                            {problem.title}
                        </Typography>
                        <Stack direction="row" spacing={1} alignItems="center">
                            <Typography variant="body2" color="textSecondary">
                                ID: {problem.id} • Slug: {problem.slug}
                            </Typography>
                            <Chip
                                label={problem.status.toUpperCase()}
                                size="small"
                                color={problem.status === 'published' ? 'success' : 'default'}
                            />
                        </Stack>
                    </Box>
                </Stack>

                <Stack direction="row" spacing={2}>
                    {problem.validation_status === 'validated' && problem.status !== 'published' && (
                        <Button
                            variant="contained"
                            color="success"
                            startIcon={<PublishIcon />}
                            onClick={handlePublish}
                            disabled={publishing}
                        >
                            {publishing ? 'Publishing...' : 'Publish Problem'}
                        </Button>
                    )}
                </Stack>
            </Box>

            {/* Workflow Steps */}
            <Paper sx={{ p: 4, mb: 4 }}>
                <Stepper activeStep={activeStep}>
                    {STEPS.map((label) => (
                        <Step key={label}>
                            <StepLabel>{label}</StepLabel>
                        </Step>
                    ))}
                </Stepper>
                <Alert severity={activeStep === 3 ? "success" : "info"} sx={{ mt: 3 }}>
                    <AlertTitle>Next Action</AlertTitle>
                    {getNextAction()}
                </Alert>
            </Paper>

            <Grid container spacing={4}>
                {/* Left Column: Details */}
                <Grid size={{ xs: 12, md: 4 }}>
                    <Card variant="outlined" sx={{ mb: 3 }}>
                        <CardContent>
                            <Typography variant="h6" gutterBottom>
                                Details
                            </Typography>
                            <Divider sx={{ mb: 2 }} />
                            <Stack spacing={2}>
                                <Box>
                                    <Typography variant="caption" color="textSecondary">Difficulty</Typography>
                                    <Typography variant="body1" sx={{ textTransform: 'capitalize' }}>{problem.difficulty}</Typography>
                                </Box>
                                <Box>
                                    <Typography variant="caption" color="textSecondary">Function Name</Typography>
                                    <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>{problem.function_name}()</Typography>
                                </Box>
                                <Box>
                                    <Typography variant="caption" color="textSecondary">Return Type</Typography>
                                    <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>{problem.return_type}</Typography>
                                </Box>
                                <Box>
                                    <Typography variant="caption" color="textSecondary">Expected Complexity</Typography>
                                    <Typography variant="body1">Time: {problem.expected_time_complexity || 'N/A'}</Typography>
                                    <Typography variant="body1">Space: {problem.expected_space_complexity || 'N/A'}</Typography>
                                </Box>
                                <Divider />
                                <Box>
                                    <Typography variant="caption" color="textSecondary">Stats</Typography>
                                    <Typography variant="body1">Submissions: {problem.total_submissions}</Typography>
                                    <Typography variant="body1">Acceptance Rate: {problem.acceptance_rate.toFixed(1)}%</Typography>
                                </Box>
                            </Stack>
                        </CardContent>
                    </Card>

                    <Card variant="outlined">
                        <CardContent sx={{ p: 0 }}>
                            <Box sx={{ p: 2 }}>
                                <Typography variant="h6">Test Cases</Typography>
                            </Box>
                            <Divider />
                            <List dense>
                                {problem.test_cases?.map((tc: any, index: number) => (
                                    <ListItem key={tc.id}>
                                        <ListItemText
                                            primary={`Test Case ${index + 1}`}
                                            secondary={tc.is_sample ? "Sample (Public)" : "Hidden"}
                                        />
                                        {tc.is_sample ? (
                                            <Chip label="Public" size="small" color="primary" />
                                        ) : (
                                            <Chip label="Private" size="small" variant="outlined" />
                                        )}
                                    </ListItem>
                                ))}
                            </List>
                        </CardContent>
                    </Card>
                </Grid>

                {/* Right Column: Validation */}
                <Grid size={{ xs: 12, md: 8 }}>
                    <ReferenceSolutionValidator
                        problemId={problem.id}
                        onValidationSuccess={() => fetchProblem()}
                    />

                    {/* Boilerplates Status */}
                    <Card variant="outlined" sx={{ mt: 3 }}>
                        <CardContent>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                                <Typography variant="h6" display="flex" alignItems="center" gap={1}>
                                    <CodeIcon fontSize="small" /> Boilerplates
                                </Typography>
                                <Button
                                    size="small"
                                    startIcon={generatingBoilerplates ? <CircularProgress size={16} /> : <RefreshIcon />}
                                    onClick={handleRegenerateBoilerplates}
                                    disabled={generatingBoilerplates}
                                >
                                    {generatingBoilerplates ? 'Generating...' : 'Regenerate All'}
                                </Button>
                            </Box>
                            <Divider sx={{ mb: 2 }} />
                            {problem.boilerplates?.length > 0 ? (
                                <Stack direction="row" spacing={1} flexWrap="wrap">
                                    {problem.boilerplates.map((bp: any) => (
                                        <Chip
                                            key={bp.id}
                                            label={bp.language?.name || `ID: ${bp.language_id}`}
                                            size="small"
                                            variant="outlined"
                                            sx={{ mb: 1 }}
                                        />
                                    ))}
                                </Stack>
                            ) : (
                                <Typography color="textSecondary">No boilerplates generated yet.</Typography>
                            )}
                        </CardContent>
                    </Card>

                    {/* Description Card */}
                    <Card variant="outlined" sx={{ mt: 3 }}>
                        <CardContent>
                            <Typography variant="h6" gutterBottom display="flex" alignItems="center" gap={1}>
                                <AssignmentIcon fontSize="small" /> Description
                            </Typography>
                            <Divider sx={{ mb: 2 }} />
                            <Paper variant="outlined" sx={{ p: 2, bgcolor: '#fcfcfc', maxHeight: 300, overflow: 'auto' }}>
                                <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                                    {problem.description}
                                </Typography>
                            </Paper>
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </Box>
    );
};

export default ProblemManagement;
