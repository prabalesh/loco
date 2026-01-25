import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Box,
    Typography,
    Button,
    Stack,
    CircularProgress,
    IconButton,
    Divider,
    Alert,
    Chip,
    Paper,
    Stepper,
    Step,
    StepLabel,
} from '@mui/material';
import {
    ArrowBack as ArrowBackIcon,
    Publish as PublishIcon,
    NavigateNext as NextIcon,
    NavigateBefore as BeforeIcon,
} from '@mui/icons-material';
import { adminProblemApi } from '../lib/api/admin';
import { ReferenceSolutionValidator } from '../components/v2/ReferenceSolutionValidator';
import { SignatureStep } from '../features/problems/components/wizard/SignatureStep';
import { BoilerplateStep } from '../features/problems/components/wizard/BoilerplateStep';
import { GeneralInfoStep } from '../features/problems/components/wizard/GeneralInfoStep';
import { TestCasesStep } from '../features/problems/components/wizard/TestCasesStep';
import toast from 'react-hot-toast';

const STEPS = [
    'Problem Details',
    'Function Signature',
    'Boilerplates',
    'Test Cases',
    'Verification',
];

const ProblemManagement: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [problem, setProblem] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [publishing, setPublishing] = useState(false);
    const [activeStep, setActiveStep] = useState(0);
    const [isDirty, setIsDirty] = useState(false);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (id) {
            fetchProblem();
        }
    }, [id]);

    const fetchProblem = async () => {
        try {
            const response = await adminProblemApi.v2GetById(id!);
            const p = response.data.data as any;
            // Normalize data for steps
            const normalizedProblem = {
                ...p,
                selected_languages: p.boilerplates?.map((b: any) => b.language?.language_id || b.language_id) || p.selected_languages || [],
                return_type: p.return_type === 'int' ? 'integer' :
                    p.return_type === 'int[]' ? 'integer_array' :
                        p.return_type === 'bool' ? 'boolean' : p.return_type,
                parameters: p.parameters?.map((param: any) => ({
                    ...param,
                    type: param.type === 'int' ? 'integer' :
                        param.type === 'int[]' ? 'integer_array' :
                            param.type === 'bool' ? 'boolean' : param.type
                })) || [],
                test_cases: p.test_cases?.map((tc: any) => ({
                    ...tc,
                    input: typeof tc.input === 'string' ? tc.input : JSON.stringify(tc.input),
                    expected_output: typeof tc.expected_output === 'string' ? tc.expected_output : JSON.stringify(tc.expected_output)
                })) || [],
                tag_ids: p.tags?.map((t: any) => t.id) || [],
                category_ids: p.categories?.map((c: any) => c.id) || []
            };
            setProblem(normalizedProblem);
            setIsDirty(false);
        } catch (error) {
            toast.error('Failed to load problem');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handlePublish = async () => {
        setPublishing(true);
        try {
            await adminProblemApi.v2Publish(id!);
            toast.success('Problem published successfully!');
            fetchProblem();
        } catch (error: any) {
            const message = error.response?.data?.data?.message || 'Failed to publish problem';
            toast.error(message);
        } finally {
            setPublishing(false);
        }
    };

    const handleProblemChange = (newData: Partial<any>) => {
        setProblem((prev: any) => ({ ...prev, ...newData }));
        setIsDirty(true);
    };

    const handleSaveChanges = async () => {
        setSaving(true);
        try {
            // Re-format test cases back to JSON for API
            const formattedTestCases = problem.test_cases.map((tc: any) => ({
                ...tc,
                input: JSON.parse(tc.input),
                expected_output: JSON.parse(tc.expected_output)
            }));

            await adminProblemApi.update(id!, {
                title: problem.title,
                slug: problem.slug,
                description: problem.description,
                difficulty: problem.difficulty,
                function_name: problem.function_name,
                return_type: problem.return_type,
                parameters: problem.parameters,
                validation_type: problem.validation_type,
                selected_languages: problem.selected_languages,
                status: problem.status,
                is_active: problem.is_active,
                constraints: problem.constraints,
                input_format: problem.input_format,
                output_format: problem.output_format,
                time_limit: problem.time_limit,
                memory_limit: problem.memory_limit,
                validator_type: problem.validator_type || 'exact_match',
                test_cases: formattedTestCases,
                tag_ids: problem.tag_ids || [],
                category_ids: problem.category_ids || [],
            });
            toast.success('Problem updated successfully!');
            setIsDirty(false);
        } catch (error) {
            toast.error('Failed to update problem');
        } finally {
            setSaving(false);
        }
    };

    const handleNext = () => {
        setActiveStep((prev) => prev + 1);
    };

    const handleBack = () => {
        setActiveStep((prev) => prev - 1);
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

    const renderStepContent = (step: number) => {
        switch (step) {
            case 0:
                return (
                    <Box sx={{ mt: 2 }}>
                        <GeneralInfoStep data={problem} onChange={handleProblemChange} />
                    </Box>
                );
            case 1:
                return (
                    <Box sx={{ mt: 2 }}>
                        <SignatureStep data={problem} onChange={handleProblemChange} />
                    </Box>
                );
            case 2:
                return (
                    <Box sx={{ mt: 2 }}>
                        <BoilerplateStep
                            data={problem}
                            onChange={handleProblemChange}
                            onRefresh={() => fetchProblem()}
                        />
                    </Box>
                );
            case 3:
                return (
                    <Box sx={{ mt: 2 }}>
                        <TestCasesStep
                            data={problem}
                            onChange={handleProblemChange}
                            onSave={handleSaveChanges}
                            saving={saving}
                        />
                    </Box>
                );
            case 4:
                return (
                    <Box sx={{ mt: 2 }}>
                        <ReferenceSolutionValidator
                            problemId={problem.id}
                            supportedLanguages={problem.selected_languages}
                            onValidationSuccess={() => fetchProblem()}
                        />
                        {problem.status !== 'published' && problem.validation_status === 'validated' && (
                            <Alert severity="success" sx={{ mt: 3, borderRadius: 2 }}>
                                This problem is validated and ready for publication.
                            </Alert>
                        )}
                    </Box>
                );
            default:
                return null;
        }
    };

    return (
        <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
            {/* Header */}
            <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <Stack direction="row" spacing={2} alignItems="center">
                    <IconButton onClick={() => navigate('/problems')} size="small" sx={{ bgcolor: '#f5f5f5' }}>
                        <ArrowBackIcon />
                    </IconButton>
                    <Box>
                        <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
                            Update Problem
                        </Typography>
                        <Stack direction="row" spacing={1} alignItems="center">
                            <Typography variant="body2" color="textSecondary">
                                {problem.title} (ID: {problem.id})
                            </Typography>
                            <Chip
                                label={problem.status.toUpperCase()}
                                size="small"
                                color={problem.status === 'published' ? 'success' : 'default'}
                                sx={{ borderRadius: '4px', fontWeight: 'bold' }}
                            />
                        </Stack>
                    </Box>
                </Stack>

                <Stack direction="row" spacing={2}>
                    {isDirty && (
                        <Button
                            variant="contained"
                            color="primary"
                            onClick={handleSaveChanges}
                            disabled={saving}
                            sx={{ borderRadius: 2, px: 3 }}
                        >
                            {saving ? 'Saving...' : 'Save Changes'}
                        </Button>
                    )}
                    {problem.status !== 'published' && (
                        <Button
                            variant="contained"
                            color="success"
                            startIcon={<PublishIcon />}
                            onClick={handlePublish}
                            disabled={publishing || problem.validation_status !== 'validated'}
                            sx={{ borderRadius: 2, px: 3 }}
                        >
                            {publishing ? 'Publishing...' : 'Publish Problem'}
                        </Button>
                    )}
                </Stack>
            </Box>

            <Paper sx={{ p: 4, mb: 4, borderRadius: 3, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 5 }}>
                    {STEPS.map((label) => (
                        <Step key={label}>
                            <StepLabel>{label}</StepLabel>
                        </Step>
                    ))}
                </Stepper>

                <Box sx={{ minHeight: 400, mb: 4 }}>
                    {renderStepContent(activeStep)}
                </Box>

                <Divider sx={{ mb: 3 }} />

                <Stack direction="row" justifyContent="space-between">
                    <Button
                        disabled={activeStep === 0 || saving}
                        onClick={handleBack}
                        startIcon={<BeforeIcon />}
                        variant="outlined"
                        sx={{ borderRadius: 2 }}
                    >
                        Back
                    </Button>
                    <Button
                        disabled={activeStep === STEPS.length - 1 || saving}
                        onClick={handleNext}
                        endIcon={<NextIcon />}
                        variant="contained"
                        sx={{ borderRadius: 2, px: 4 }}
                    >
                        Next
                    </Button>
                </Stack>
            </Paper>
        </Box>
    );
};

export default ProblemManagement;
