import React, { useState, useEffect } from 'react';
import {
    Box,
    Stepper,
    Step,
    StepLabel,
    Button,
    Typography,
    Paper,
    Container,
    Stack,
    IconButton,
    Divider
} from '@mui/material';
import {
    ArrowBack as ArrowBackIcon,
    NavigateNext as NextIcon,
    NavigateBefore as BeforeIcon,
    CheckCircle as CheckCircleIcon
} from '@mui/icons-material';
import { useNavigate, useSearchParams } from 'react-router-dom';
import toast from 'react-hot-toast';

// Steps
import { GeneralInfoStep } from './GeneralInfoStep';
import { BoilerplateStep } from './BoilerplateStep';
import { SignatureStep } from './SignatureStep';
import { TestCasesStep } from './TestCasesStep';
import { VerificationStep } from './VerificationStep';
import { adminProblemApi } from '../../../../lib/api/admin';

const STEPS = [
    'Problem Details',
    'Function Signature',
    'Boilerplates',
    'Test Cases',
    'Verify & Publish'
];

export const ProblemWizard: React.FC = () => {
    const [activeStep, setActiveStep] = useState(0);
    const [searchParams, setSearchParams] = useSearchParams();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [problemId, setProblemId] = useState<number | null>(null);

    // Shared state across steps
    const [formData, setFormData] = useState<any>({
        title: '',
        difficulty: 'medium',
        description: '',
        constraints: '',
        hints: '',
        slug: '',
        function_name: 'solve',
        return_type: 'integer',
        parameters: [{ name: 'nums', type: 'integer_array', is_custom: false }],
        test_cases: [{ input: '[]', expected_output: 'null', is_sample: true }],
        validation_type: 'EXACT',
        selected_languages: [],
        tag_ids: [],
        category_ids: []
    });

    useEffect(() => {
        const step = parseInt(searchParams.get('step') || '0');
        const id = searchParams.get('id');
        if (id) {
            const numId = parseInt(id);
            if (!problemId) {
                setProblemId(numId);
                fetchProblemData(numId);
            }
        }
        setActiveStep(step);
    }, [searchParams, problemId]);

    const fetchProblemData = async (id: number) => {
        try {
            const response = await adminProblemApi.v2GetById(id);
            const p = response.data.data as any;
            setFormData({
                title: p.title,
                difficulty: p.difficulty,
                description: p.description,
                constraints: p.constraints || '',
                hints: p.hints || '',
                slug: p.slug,
                function_name: p.function_name || 'solve',
                return_type: p.return_type || 'integer',
                parameters: p.parameters || [{ name: 'nums', type: 'integer_array', is_custom: false }],
                test_cases: p.test_cases?.map((tc: any) => ({
                    ...tc,
                    input: JSON.stringify(tc.input),
                    expected_output: JSON.stringify(tc.expected_output)
                })) || [{ input: '[]', expected_output: 'null', is_sample: true }],
                validation_type: p.validation_type || 'EXACT',
                selected_languages: p.boilerplates?.map((b: any) => b.language?.language_id || b.language_id) || [],
                boilerplates: p.boilerplates || [],
                tag_ids: p.tags?.map((t: any) => t.id) || [],
                category_ids: p.categories?.map((c: any) => c.id) || []
            });
        } catch (err) {
            console.error('Failed to fetch problem data', err);
        }
    };

    const updateFormData = (newData: Partial<any>) => {
        setFormData((prev: any) => ({ ...prev, ...newData }));
    };

    const handleNext = () => {
        const nextStep = activeStep + 1;
        setSearchParams({ step: nextStep.toString(), id: problemId?.toString() || '' });
    };

    const handleBack = () => {
        const prevStep = activeStep - 1;
        setSearchParams({ step: prevStep.toString(), id: problemId?.toString() || '' });
    };

    const handlePublish = async () => {
        if (!problemId) return;
        setLoading(true);
        try {
            await adminProblemApi.v2Publish(problemId.toString());
            toast.success('Problem Published Successfully!');
            navigate('/problems');
        } catch (err: any) {
            toast.error(err.message || 'Failed to publish problem');
        } finally {
            setLoading(false);
        }
    };

    const handleSaveTestCases = async () => {
        setLoading(true);
        try {
            const formattedTestCases = formData.test_cases.map((tc: any) => ({
                ...tc,
                input: typeof tc.input === 'string' ? JSON.parse(tc.input) : tc.input,
                expected_output: typeof tc.expected_output === 'string' ? JSON.parse(tc.expected_output) : tc.expected_output
            }));

            const problemData = {
                title: formData.title,
                difficulty: formData.difficulty,
                description: formData.description,
                constraints: formData.constraints,
                hints: formData.hints,
                function_name: formData.function_name,
                return_type: formData.return_type,
                parameters: formData.parameters,
                test_cases: formattedTestCases,
                validation_type: formData.validation_type,
                status: 'draft' as 'draft',
                visibility: 'public' as 'public',
                is_active: true,
                slug: formData.slug || '',
                time_limit: formData.time_limit || 1000,
                memory_limit: formData.memory_limit || 256,
                validator_type: formData.validator_type || 'exact_match',
                input_format: formData.input_format || '',
                output_format: formData.output_format || '',
                tag_ids: formData.tag_ids || [],
                category_ids: formData.category_ids || []
            };

            if (problemId) {
                await adminProblemApi.update(problemId.toString(), problemData);
                toast.success('Test cases saved successfully!');
            } else {
                const response = await adminProblemApi.v2Create(problemData);
                const newId = response.data.data.id;
                setProblemId(newId);
                setSearchParams({ step: '3', id: newId.toString() });
                toast.success('Problem draft created and test cases saved!');
            }
        } catch (err: any) {
            toast.error(err.message || 'Failed to save test cases');
        } finally {
            setLoading(false);
        }
    };

    const renderStepContent = (step: number) => {
        switch (step) {
            case 0:
                return <GeneralInfoStep data={formData} onChange={updateFormData} />;
            case 1:
                return <SignatureStep data={formData} onChange={updateFormData} />;
            case 2:
                return <BoilerplateStep
                    data={formData}
                    onChange={updateFormData}
                    onRefresh={() => problemId && fetchProblemData(problemId)}
                />;
            case 3:
                return <TestCasesStep
                    data={formData}
                    onChange={updateFormData}
                    onSave={handleSaveTestCases}
                    saving={loading}
                    problemId={problemId || undefined}
                />;
            case 4:
                return (
                    <VerificationStep
                        data={formData}
                        problemId={problemId}
                        onProblemCreated={(id) => {
                            setProblemId(id);
                            setSearchParams({ step: '4', id: id.toString() });
                        }}
                    />
                );
            default:
                return <Typography>Unknown step</Typography>;
        }
    };

    return (
        <Box sx={{ py: 4 }}>
            <Container maxWidth="lg">
                <Paper sx={{ p: 4, borderRadius: 2, boxShadow: '0 4px 20px rgba(0,0,0,0.08)' }}>
                    <Stack direction="row" spacing={2} alignItems="center" sx={{ mb: 4 }}>
                        <IconButton onClick={() => navigate('/problems')} size="small" sx={{ bgcolor: '#f5f5f5' }}>
                            <ArrowBackIcon />
                        </IconButton>
                        <Typography variant="h5" fontWeight="bold" color="text.primary">
                            {problemId ? 'Edit Problem' : 'Create New Problem'}
                        </Typography>
                    </Stack>

                    <Stepper activeStep={activeStep} alternativeLabel sx={{ mb: 5 }}>
                        {STEPS.map((label) => (
                            <Step key={label}>
                                <StepLabel>{label}</StepLabel>
                            </Step>
                        ))}
                    </Stepper>

                    <Box sx={{ minHeight: '400px', mb: 4 }}>
                        {renderStepContent(activeStep)}
                    </Box>

                    <Divider sx={{ mb: 3 }} />

                    <Stack direction="row" justifyContent="space-between">
                        <Button
                            disabled={activeStep === 0 || loading}
                            onClick={handleBack}
                            startIcon={<BeforeIcon />}
                            variant="outlined"
                        >
                            Back
                        </Button>
                        <Box sx={{ display: 'flex', gap: 2 }}>
                            {!problemId && (
                                <Button
                                    variant="outlined"
                                    color="info"
                                    onClick={handleSaveTestCases}
                                    disabled={loading}
                                    sx={{ px: 4 }}
                                >
                                    Save Draft
                                </Button>
                            )}
                            {activeStep === STEPS.length - 1 ? (
                                <Button
                                    variant="contained"
                                    color="success"
                                    onClick={handlePublish}
                                    startIcon={<CheckCircleIcon />}
                                    disabled={loading || !problemId}
                                    sx={{ px: 4 }}
                                >
                                    Publish Problem
                                </Button>
                            ) : (
                                <Button
                                    variant="contained"
                                    onClick={handleNext}
                                    endIcon={<NextIcon />}
                                    disabled={loading}
                                    sx={{ px: 4 }}
                                >
                                    Next
                                </Button>
                            )}
                        </Box>
                    </Stack>
                </Paper>
            </Container>
        </Box>
    );
};

export default ProblemWizard;
