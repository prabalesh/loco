import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Box,
    Card,
    CardContent,
    Typography,
    Tabs,
    Tab,
    Button,
    Stack,
    CircularProgress,
    IconButton,
    Divider,
    Alert,
    Chip,
    Paper,
} from '@mui/material';
import {
    ArrowBack as ArrowBackIcon,
    CheckCircle as CheckCircleIcon,
    Publish as PublishIcon,
    Edit as EditIcon,
} from '@mui/icons-material';
import { adminProblemApi } from '../lib/api/admin';
import { ReferenceSolutionValidator } from '../components/v2/ReferenceSolutionValidator';
import toast from 'react-hot-toast';

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`problem-tabpanel-${index}`}
            aria-labelledby={`problem-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ pt: 3 }}>
                    {children}
                </Box>
            )}
        </div>
    );
}

export const ProblemManagement: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [problem, setProblem] = useState<any>(null);
    const [loading, setLoading] = useState(true);
    const [publishing, setPublishing] = useState(false);
    const [activeTab, setActiveTab] = useState(0);

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

    const handlePublish = async () => {
        setPublishing(true);
        try {
            await adminProblemApi.v2Publish(id!);
            toast.success('Problem published successfully!');
            fetchProblem(); // Refresh to see updated status
        } catch (error: any) {
            const message = error.response?.data?.data?.message || 'Failed to publish problem';
            toast.error(message);
        } finally {
            setPublishing(false);
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
                <Button startIcon={<ArrowBackIcon />} onClick={() => navigate('/admin/problems')} sx={{ mt: 2 }}>
                    Back to Problems
                </Button>
            </Box>
        );
    }

    return (
        <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
            <Box sx={{ mb: 4, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <Stack direction="row" spacing={2} alignItems="center">
                    <IconButton onClick={() => navigate('/admin/problems')} size="small">
                        <ArrowBackIcon />
                    </IconButton>
                    <Box>
                        <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}>
                            {problem.title}
                            {problem.status === 'published' && (
                                <CheckCircleIcon color="success" />
                            )}
                        </Typography>
                        <Typography variant="body2" color="textSecondary">
                            Target ID: {problem.id} • Slug: {problem.slug}
                        </Typography>
                    </Box>
                </Stack>

                <Stack direction="row" spacing={2}>
                    <Button
                        variant="outlined"
                        startIcon={<EditIcon />}
                        onClick={() => navigate(`/admin/problems/${problem.id}/edit`)}
                    >
                        Edit Base Problem
                    </Button>
                    {problem.validation_status === 'validated' && problem.status !== 'published' && (
                        <Button
                            variant="contained"
                            color="success"
                            startIcon={<PublishIcon />}
                            onClick={handlePublish}
                            loading={publishing}
                        >
                            Publish Problem
                        </Button>
                    )}
                </Stack>
            </Box>

            <Card variant="outlined">
                <CardContent sx={{ p: 0 }}>
                    <Box sx={{ borderBottom: 1, borderColor: 'divider', px: 2 }}>
                        <Tabs value={activeTab} onChange={(_, newValue) => setActiveTab(newValue)}>
                            <Tab label="Overview" />
                            <Tab label="Validation" />
                            <Tab label="Test Cases" />
                            <Tab label="JSON Review" />
                        </Tabs>
                    </Box>

                    <Box sx={{ px: 3, pb: 2 }}>
                        <TabPanel value={activeTab} index={0}>
                            <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: 4 }}>
                                <Box>
                                    <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>Status Info</Typography>
                                    <Stack spacing={1}>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                            <Typography variant="body2" color="textSecondary">Publication Status:</Typography>
                                            <Chip label={problem.status} size="small" color={problem.status === 'published' ? "success" : "default"} />
                                        </Box>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                            <Typography variant="body2" color="textSecondary">Validation Status:</Typography>
                                            <Chip label={problem.validation_status} size="small" color={problem.validation_status === 'validated' ? "success" : "warning"} />
                                        </Box>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                            <Typography variant="body2" color="textSecondary">Visibility:</Typography>
                                            <Typography variant="body2" sx={{ fontWeight: 'medium' }}>{problem.visibility}</Typography>
                                        </Box>
                                    </Stack>
                                </Box>

                                <Box>
                                    <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>Technical Details</Typography>
                                    <Stack spacing={1}>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                            <Typography variant="body2" color="textSecondary">Function Name:</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>{problem.function_name}</Typography>
                                        </Box>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                            <Typography variant="body2" color="textSecondary">Return Type:</Typography>
                                            <Typography variant="body2" sx={{ fontFamily: 'monospace' }}>{problem.return_type}</Typography>
                                        </Box>
                                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                                            <Typography variant="body2" color="textSecondary">Complexity:</Typography>
                                            <Typography variant="body2">{problem.expected_time_complexity || 'N/A'}</Typography>
                                        </Box>
                                    </Stack>
                                </Box>
                            </Box>

                            <Box sx={{ mt: 4 }}>
                                <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'bold' }}>Description Preview</Typography>
                                <Paper variant="outlined" sx={{ p: 2, bgcolor: '#fafafa', maxHeight: 300, overflow: 'auto' }}>
                                    <Typography variant="body2" sx={{ whiteSpace: 'pre-wrap' }}>
                                        {problem.description}
                                    </Typography>
                                </Paper>
                            </Box>
                        </TabPanel>

                        <TabPanel value={activeTab} index={1}>
                            <ReferenceSolutionValidator problemId={problem.id} />
                        </TabPanel>

                        <TabPanel value={activeTab} index={2}>
                            <Typography variant="body1" color="textSecondary">
                                Test case management for existing problems coming soon.
                                Currently, test cases are defined during initial problem creation.
                            </Typography>
                        </TabPanel>

                        <TabPanel value={activeTab} index={3}>
                            <Paper variant="outlined" sx={{ p: 2, bgcolor: '#1e1e1e', color: '#d4d4d4' }}>
                                <pre style={{ margin: 0, fontSize: '0.8rem', overflow: 'auto' }}>
                                    {JSON.stringify(problem, null, 2)}
                                </pre>
                            </Paper>
                        </TabPanel>
                    </Box>
                </CardContent>
            </Card>
        </Box>
    );
};

export default ProblemManagement;
