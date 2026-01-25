import React, { useEffect, useState } from 'react';
import {
    Box,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Typography,
    CircularProgress,
    Pagination,
    Chip,
    Button
} from '@mui/material';
import { Link, useNavigate } from 'react-router-dom';
import { adminSubmissionsApi } from '../lib/api/admin';
import toast from 'react-hot-toast';

const SubmissionsList: React.FC = () => {
    const [submissions, setSubmissions] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const limit = 10;
    const navigate = useNavigate();

    const fetchSubmissions = async () => {
        setLoading(true);
        try {
            const response = await adminSubmissionsApi.listSubmissions(page, limit);
            setSubmissions(response.data.data);
            setTotalPages(Math.ceil(response.data.total / limit));
        } catch (error) {
            toast.error('Failed to fetch submissions');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchSubmissions();
    }, [page]);

    const handlePageChange = (_: React.ChangeEvent<unknown>, value: number) => {
        setPage(value);
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'Accepted':
                return 'success';
            case 'Wrong Answer':
                return 'error';
            case 'Time Limit Exceeded':
                return 'warning';
            case 'Runtime Error':
                return 'error';
            case 'Compilation Error':
                return 'error';
            default:
                return 'default';
        }
    };

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom component="div" sx={{ mb: 4, fontWeight: 'bold' }}>
                All Submissions
            </Typography>

            <Paper sx={{ width: '100%', mb: 2, borderRadius: 2, boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}>
                {loading ? (
                    <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
                        <CircularProgress />
                    </Box>
                ) : (
                    <>
                        <TableContainer>
                            <Table aria-label="submissions table">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>ID</TableCell>
                                        <TableCell>User</TableCell>
                                        <TableCell>Problem</TableCell>
                                        <TableCell>Language</TableCell>
                                        <TableCell>Status</TableCell>
                                        <TableCell align="right">Runtime (ms)</TableCell>
                                        <TableCell align="right">Memory (KB)</TableCell>
                                        <TableCell align="right">Created At</TableCell>
                                        <TableCell align="right">Action</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {submissions.map((row) => (
                                        <TableRow key={row.id}>
                                            <TableCell>{row.id}</TableCell>
                                            <TableCell>
                                                <Link to={`/users/${row.user_id}`} style={{ textDecoration: 'none', color: 'inherit' }} className="hover:underline text-blue-600">
                                                    {row.user?.username || row.user_id}
                                                </Link>
                                            </TableCell>
                                            <TableCell>
                                                <Link to={`/problems/${row.problem_id}/manage`} style={{ textDecoration: 'none', color: 'inherit' }} className="hover:underline text-blue-600">
                                                    {row.problem?.title || row.problem_id}
                                                </Link>
                                            </TableCell>
                                            <TableCell>
                                                <Chip label={row.language?.name || row.language_id} size="small" variant="outlined" />
                                            </TableCell>
                                            <TableCell>
                                                <Chip
                                                    label={row.status}
                                                    color={getStatusColor(row.status) as any}
                                                    size="small"
                                                />
                                            </TableCell>
                                            <TableCell align="right">{row.runtime}</TableCell>
                                            <TableCell align="right">{row.memory}</TableCell>
                                            <TableCell align="right">{new Date(row.created_at).toLocaleString()}</TableCell>
                                            <TableCell align="right">
                                                {/* Placeholder for detail view if implemented later */}
                                                <Button size="small" disabled>View</Button>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                    {submissions.length === 0 && (
                                        <TableRow>
                                            <TableCell colSpan={9} align="center" sx={{ py: 3 }}>
                                                No submissions found
                                            </TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        </TableContainer>
                        <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
                            <Pagination
                                count={totalPages}
                                page={page}
                                onChange={handlePageChange}
                                color="primary"
                            />
                        </Box>
                    </>
                )}
            </Paper>
        </Box>
    );
};

export default SubmissionsList;
