import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Typography,
    Button,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
    Chip,
    IconButton,
    Stack,
    TextField,
    MenuItem,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    CircularProgress,
    Tooltip,
} from '@mui/material';
import {
    Add as AddIcon,
    Delete as DeleteIcon,
    Settings as ManageIcon,
    FilterList as FilterIcon,
} from '@mui/icons-material';
import { adminProblemApi } from '../lib/api/admin';
import type { Problem } from '../types';
import toast from 'react-hot-toast';

const ProblemsList: React.FC = () => {
    const navigate = useNavigate();
    const [problems, setProblems] = useState<Problem[]>([]);
    const [loading, setLoading] = useState(true);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(20);

    // Filters
    const [statusFilter, setStatusFilter] = useState('');
    const [difficultyFilter, setDifficultyFilter] = useState('');

    // Delete confirmation
    const [deleteId, setDeleteId] = useState<number | null>(null);
    const [deleting, setDeleting] = useState(false);

    useEffect(() => {
        fetchProblems();
    }, [page, rowsPerPage, statusFilter, difficultyFilter]);

    const fetchProblems = async () => {
        setLoading(true);
        try {
            const filters: any = {};
            if (statusFilter) filters.status = statusFilter;
            if (difficultyFilter) filters.difficulty = difficultyFilter;

            const response = await adminProblemApi.v2List(page + 1, rowsPerPage, filters);
            console.log(response.data.data);
            setProblems(response.data.data ?? []);
            setTotal(response.data.total);
        } catch (error) {
            toast.error('Failed to load problems');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async () => {
        if (!deleteId) return;
        setDeleting(true);
        try {
            await adminProblemApi.v2Delete(deleteId);
            toast.success('Problem deleted successfully');
            fetchProblems();
        } catch (error) {
            toast.error('Failed to delete problem');
        } finally {
            setDeleting(false);
            setDeleteId(null);
        }
    };

    const getDifficultyColor = (difficulty: string) => {
        switch (difficulty) {
            case 'easy': return 'success';
            case 'medium': return 'warning';
            case 'hard': return 'error';
            default: return 'default';
        }
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'published': return 'success';
            case 'draft': return 'default';
            default: return 'default';
        }
    };

    return (
        <Box sx={{ p: 4 }}>
            <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="h4" component="h1" fontWeight="bold">
                    Problem Management
                </Typography>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => navigate('/problems/create')}
                >
                    Create New Problem
                </Button>
            </Box>

            <Paper sx={{ mb: 3, p: 2 }}>
                <Stack direction="row" spacing={2} alignItems="center">
                    <FilterIcon color="action" />
                    <TextField
                        select
                        size="small"
                        label="Status"
                        value={statusFilter}
                        onChange={(e) => setStatusFilter(e.target.value)}
                        sx={{ width: 150 }}
                    >
                        <MenuItem value="">All Status</MenuItem>
                        <MenuItem value="draft">Draft</MenuItem>
                        <MenuItem value="published">Published</MenuItem>
                    </TextField>

                    <TextField
                        select
                        size="small"
                        label="Difficulty"
                        value={difficultyFilter}
                        onChange={(e) => setDifficultyFilter(e.target.value)}
                        sx={{ width: 150 }}
                    >
                        <MenuItem value="">All Difficulty</MenuItem>
                        <MenuItem value="easy">Easy</MenuItem>
                        <MenuItem value="medium">Medium</MenuItem>
                        <MenuItem value="hard">Hard</MenuItem>
                    </TextField>

                    <Button
                        variant="text"
                        onClick={() => { setStatusFilter(''); setDifficultyFilter(''); }}
                        size="small"
                    >
                        Clear Filters
                    </Button>
                </Stack>
            </Paper>

            <TableContainer component={Paper} elevation={2}>
                <Table>
                    <TableHead sx={{ bgcolor: '#f5f5f5' }}>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 'bold' }}>Title</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }}>Difficulty</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }}>Classification</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }}>Status</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }}>Validation</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }}>Submissions</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }}>Acceptance</TableCell>
                            <TableCell sx={{ fontWeight: 'bold' }} align="right">Actions</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {loading && problems.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} align="center" sx={{ py: 10 }}>
                                    <CircularProgress />
                                </TableCell>
                            </TableRow>
                        ) : problems.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} align="center" sx={{ py: 10 }}>
                                    No problems found
                                </TableCell>
                            </TableRow>
                        ) : (
                            problems.map((prob) => (
                                <TableRow key={prob.id} hover>
                                    <TableCell>
                                        <Typography variant="subtitle2" fontWeight="bold">
                                            {prob.title}
                                        </Typography>
                                        <Typography variant="caption" color="textSecondary">
                                            {prob.slug}
                                        </Typography>
                                    </TableCell>
                                    <TableCell>
                                        <Chip
                                            label={prob.difficulty.toUpperCase()}
                                            size="small"
                                            color={getDifficultyColor(prob.difficulty) as any}
                                            variant="outlined"
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <Stack direction="row" spacing={0.5} flexWrap="wrap">
                                            {prob.tags?.map(tag => (
                                                <Chip key={tag.id} label={tag.name} size="small" variant="outlined" sx={{ fontSize: '0.65rem', height: 20 }} />
                                            ))}
                                            {prob.categories?.map(cat => (
                                                <Chip key={cat.id} label={cat.name} size="small" color="secondary" variant="outlined" sx={{ fontSize: '0.65rem', height: 20 }} />
                                            ))}
                                        </Stack>
                                    </TableCell>
                                    <TableCell>
                                        <Chip
                                            label={prob.status.toUpperCase()}
                                            size="small"
                                            color={getStatusColor(prob.status) as any}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <Chip
                                            label={prob.validation_status.toUpperCase()}
                                            size="small"
                                            variant="outlined"
                                            color={prob.validation_status === 'validated' ? 'success' : 'warning'}
                                        />
                                    </TableCell>
                                    <TableCell>{prob.total_submissions}</TableCell>
                                    <TableCell>
                                        {prob.acceptance_rate.toFixed(1)}%
                                    </TableCell>
                                    <TableCell align="right">
                                        <Stack direction="row" spacing={1} justifyContent="flex-end">
                                            <Tooltip title="Manage Problem">
                                                <IconButton
                                                    size="small"
                                                    color="primary"
                                                    onClick={() => navigate(`/problems/${prob.id}/manage`)}
                                                >
                                                    <ManageIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Delete">
                                                <IconButton
                                                    size="small"
                                                    color="error"
                                                    onClick={() => setDeleteId(prob.id)}
                                                >
                                                    <DeleteIcon fontSize="small" />
                                                </IconButton>
                                            </Tooltip>
                                        </Stack>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
                <TablePagination
                    rowsPerPageOptions={[10, 20, 50]}
                    component="div"
                    count={total}
                    rowsPerPage={rowsPerPage}
                    page={page}
                    onPageChange={(_, newPage) => setPage(newPage)}
                    onRowsPerPageChange={(e) => {
                        setRowsPerPage(parseInt(e.target.value, 10));
                        setPage(0);
                    }}
                />
            </TableContainer>

            {/* Delete Confirmation */}
            <Dialog open={!!deleteId} onClose={() => setDeleteId(null)}>
                <DialogTitle>Confirm Delete</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Are you sure you want to delete this problem? This action cannot be undone.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDeleteId(null)}>Cancel</Button>
                    <Button
                        onClick={handleDelete}
                        color="error"
                        variant="contained"
                        disabled={deleting}
                    >
                        {deleting ? 'Deleting...' : 'Delete'}
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    );
};

export default ProblemsList;
