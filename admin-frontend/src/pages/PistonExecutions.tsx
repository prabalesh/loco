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
    Stack,
    IconButton,
    Collapse,
    Chip,
    Card,
    CardContent
} from '@mui/material';
import {
    KeyboardArrowDown as KeyboardArrowDownIcon,
    KeyboardArrowUp as KeyboardArrowUpIcon
} from '@mui/icons-material';
import { adminAnalyticsApi } from '../lib/api/admin';
import { type PistonExecution } from '../types';
import toast from 'react-hot-toast';
import { Link } from 'react-router-dom';

const Row: React.FC<{ row: PistonExecution }> = ({ row }) => {
    const [open, setOpen] = useState(false);

    const getStatus = (row: PistonExecution) => {
        let response: any = {};
        if (typeof row.response === 'string') {
            try {
                response = JSON.parse(row.response);
            } catch (e) {
                return { label: 'Invalid JSON', color: 'error' };
            }
        } else {
            response = row.response;
        }

        if (response?.run?.code !== 0) {
            return { label: 'Error', color: 'error' };
        }
        return { label: 'Success', color: 'success' };
    };

    const status = getStatus(row);

    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">
                    {row.id}
                </TableCell>
                <TableCell align="right">
                    <Link to={`/problems/${row.problem_id}/manage`} style={{ textDecoration: 'none', color: 'inherit' }} className="hover:underline text-blue-600">
                        {row.problem?.title || row.problem_id}
                    </Link>
                </TableCell>
                <TableCell align="right">
                    {row.submission_id ? (
                        <Link to={`/submissions`} style={{ textDecoration: 'none', color: 'inherit' }} className="hover:underline text-blue-600">
                            {row.submission_id}
                        </Link>
                    ) : '-'}
                </TableCell>
                <TableCell align="right">
                    <Chip label={row.language} size="small" color="primary" variant="outlined" />
                </TableCell>
                <TableCell align="right">
                    <Chip label={status.label} size="small" color={status.color as any} />
                </TableCell>
                <TableCell align="right">{new Date(row.created_at).toLocaleString()}</TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={7}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 2 }}>
                            <Typography variant="h6" gutterBottom component="div">
                                Execution Details
                            </Typography>
                            <Stack direction="row" spacing={2} sx={{ mb: 2 }}>
                                <Card variant="outlined" sx={{ flex: 1 }}>
                                    <CardContent>
                                        <Typography color="textSecondary" gutterBottom>
                                            Response Output
                                        </Typography>
                                        <pre style={{ overflow: 'auto', maxHeight: '200px' }}>
                                            {JSON.stringify(row.response, null, 2)}
                                        </pre>
                                    </CardContent>
                                </Card>
                            </Stack>
                            <Stack direction="row" spacing={2}>
                                <Card variant="outlined" sx={{ flex: 1 }}>
                                    <CardContent>
                                        <Typography color="textSecondary" gutterBottom>
                                            Code
                                        </Typography>
                                        <pre style={{ overflow: 'auto', maxHeight: '200px' }}>
                                            {row.code}
                                        </pre>
                                    </CardContent>
                                </Card>
                                <Card variant="outlined" sx={{ flex: 1 }}>
                                    <CardContent>
                                        <Typography color="textSecondary" gutterBottom>
                                            Stdin
                                        </Typography>
                                        <pre style={{ overflow: 'auto', maxHeight: '200px' }}>
                                            {row.stdin}
                                        </pre>
                                    </CardContent>
                                </Card>
                            </Stack>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
};

const PistonExecutions: React.FC = () => {
    const [executions, setExecutions] = useState<PistonExecution[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const limit = 10;

    const fetchExecutions = async () => {
        setLoading(true);
        try {
            const response = await adminAnalyticsApi.listPistonExecutions(page, limit);
            setExecutions(response.data.data);
            setTotalPages(Math.ceil(response.data.total / limit));
        } catch (error) {
            toast.error('Failed to fetch Piston executions');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchExecutions();
    }, [page]);

    const handlePageChange = (_: React.ChangeEvent<unknown>, value: number) => {
        setPage(value);
    };

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom component="div" sx={{ mb: 4, fontWeight: 'bold' }}>
                Piston Executions
            </Typography>

            <Paper sx={{ width: '100%', mb: 2, borderRadius: 2, boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}>
                {loading ? (
                    <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
                        <CircularProgress />
                    </Box>
                ) : (
                    <>
                        <TableContainer>
                            <Table aria-label="collapsible table">
                                <TableHead>
                                    <TableRow>
                                        <TableCell />
                                        <TableCell>ID</TableCell>
                                        <TableCell align="right">Problem ID</TableCell>
                                        <TableCell align="right">Submission ID</TableCell>
                                        <TableCell align="right">Language</TableCell>
                                        <TableCell align="right">Version</TableCell>
                                        <TableCell align="right">Created At</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {executions.map((row) => (
                                        <Row key={row.id} row={row} />
                                    ))}
                                    {executions.length === 0 && (
                                        <TableRow>
                                            <TableCell colSpan={7} align="center" sx={{ py: 3 }}>
                                                No executions found
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

export default PistonExecutions;
