import React, { useState } from 'react'
import {
    Box,
    Button,
    Card,
    CardContent,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Chip,
    Alert,
    LinearProgress,
    Stack,
    Divider
} from '@mui/material'
import { Upload, CheckCircle, XCircle, Clock } from 'lucide-react'
import { adminBulkApi } from '../lib/api/admin'
import { toast } from 'react-hot-toast'

interface ImportResult {
    total_submitted: number
    total_created: number
    total_failed: number
    created_problems: Array<{
        index: number
        title: string
        slug: string
        problem_id: number
        validation_status: string
    }>
    failed_problems: Array<{
        index: number
        title: string
        errors: string[]
        error_message: string
    }>
    processing_time_ms: number
}

export const BulkImport: React.FC = () => {
    const [isImporting, setIsImporting] = useState(false)
    const [result, setResult] = useState<ImportResult | null>(null)

    const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0]
        if (!file) return

        setIsImporting(true)
        setResult(null)

        try {
            const text = await file.text()
            const data = JSON.parse(text)

            const response = await adminBulkApi.import(data)
            setResult(response.data as unknown as ImportResult)

            const res = response.data as any
            if (res.total_failed === 0) {
                toast.success(`Successfully imported ${res.total_created} problems!`)
            } else {
                toast.error(`Imported ${res.total_created} problems, but ${res.total_failed} failed.`)
            }
        } catch (error: any) {
            console.error('Import failed:', error)
            toast.error(error?.response?.data?.error || 'Failed to import problems. Check JSON format.')
        } finally {
            setIsImporting(false)
            // Reset input
            event.target.value = ''
        }
    }

    return (
        <Box sx={{ p: 4, maxWidth: 1200, mx: 'auto' }}>
            <Typography variant="h4" sx={{ mb: 4, fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 2 }}>
                <Upload size={32} /> Bulk Import Problems
            </Typography>

            <Alert severity="info" sx={{ mb: 4 }}>
                <Typography variant="subtitle2" fontWeight="bold">Bulk Import Instructions:</Typography>
                <Typography variant="body2">
                    Upload a JSON file containing an array of problem definitions.
                    Maximum 100 problems per batch for synchronous processing.
                    Check the documentation for the required JSON structure.
                </Typography>
            </Alert>

            <Card sx={{ mb: 4, p: 2 }}>
                <CardContent>
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 4, border: '2px dashed #ccc', borderRadius: 2, bgcolor: '#fafafa' }}>
                        <Upload size={48} color="#666" style={{ marginBottom: 16 }} />
                        <Typography variant="h6" gutterBottom>Upload Problem Batch</Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>Accepts .json files</Typography>

                        <Button
                            variant="contained"
                            component="label"
                            disabled={isImporting}
                            startIcon={<Upload size={20} />}
                            sx={{ px: 4, py: 1.5, borderRadius: 2, textTransform: 'none', fontWeight: 'bold' }}
                        >
                            {isImporting ? 'Importing...' : 'Select JSON File'}
                            <input
                                type="file"
                                hidden
                                accept=".json"
                                onChange={handleFileUpload}
                            />
                        </Button>
                    </Box>
                    {isImporting && <LinearProgress sx={{ mt: 2, borderRadius: 1 }} />}
                </CardContent>
            </Card>

            {result && (
                <Stack spacing={4}>
                    <Card>
                        <CardContent>
                            <Typography variant="h6" gutterBottom sx={{ fontWeight: 'bold' }}>Import Summary</Typography>
                            <Divider sx={{ mb: 3 }} />

                            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 4 }}>
                                <Box>
                                    <Typography variant="caption" color="text.secondary" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                        <CheckCircle size={14} color="green" /> Total Created
                                    </Typography>
                                    <Typography variant="h5" fontWeight="bold" color="success.main">{result.total_created}</Typography>
                                </Box>
                                <Box>
                                    <Typography variant="caption" color="text.secondary" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                        <XCircle size={14} color="red" /> Total Failed
                                    </Typography>
                                    <Typography variant="h5" fontWeight="bold" color="error.main">{result.total_failed}</Typography>
                                </Box>
                                <Box>
                                    <Typography variant="caption" color="text.secondary" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                                        <Clock size={14} /> Processing Time
                                    </Typography>
                                    <Typography variant="h5" fontWeight="bold">{result.processing_time_ms}ms</Typography>
                                </Box>
                            </Box>

                            <Box sx={{ mt: 3 }}>
                                <Typography variant="body2" gutterBottom>Completion Rate</Typography>
                                <LinearProgress
                                    variant="determinate"
                                    value={(result.total_created / result.total_submitted) * 100}
                                    sx={{ height: 10, borderRadius: 5, bgcolor: '#eee', '& .MuiLinearProgress-bar': { bgcolor: result.total_failed > 0 ? 'warning.main' : 'success.main' } }}
                                />
                            </Box>
                        </CardContent>
                    </Card>

                    {result.failed_problems.length > 0 && (
                        <Card sx={{ borderLeft: '4px solid #ef4444' }}>
                            <CardContent>
                                <Typography variant="h6" color="error.main" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <XCircle size={20} /> Failed Imports ({result.failed_problems.length})
                                </Typography>
                                <TableContainer component={Paper} variant="outlined" sx={{ mt: 2 }}>
                                    <Table size="small">
                                        <TableHead sx={{ bgcolor: '#fff5f5' }}>
                                            <TableRow>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Index</TableCell>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Title</TableCell>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Error Message</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {result.failed_problems.map((fail) => (
                                                <TableRow key={fail.index}>
                                                    <TableCell>{fail.index}</TableCell>
                                                    <TableCell sx={{ fontWeight: 'medium' }}>{fail.title || 'N/A'}</TableCell>
                                                    <TableCell color="error.main">
                                                        <Typography variant="body2" color="error">
                                                            {fail.error_message}
                                                            {fail.errors && fail.errors.length > 0 && (
                                                                <Box component="ul" sx={{ m: 0, pl: 2, mt: 0.5 }}>
                                                                    {fail.errors.map((err, i) => (
                                                                        <li key={i}>{err}</li>
                                                                    ))}
                                                                </Box>
                                                            )}
                                                        </Typography>
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            </CardContent>
                        </Card>
                    )}

                    {result.created_problems.length > 0 && (
                        <Card sx={{ borderLeft: '4px solid #10b981' }}>
                            <CardContent>
                                <Typography variant="h6" color="success.main" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                    <CheckCircle size={20} /> Successfully Created ({result.created_problems.length})
                                </Typography>
                                <TableContainer component={Paper} variant="outlined" sx={{ mt: 2 }}>
                                    <Table size="small">
                                        <TableHead sx={{ bgcolor: '#f0fdf4' }}>
                                            <TableRow>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Index</TableCell>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Title</TableCell>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Slug</TableCell>
                                                <TableCell sx={{ fontWeight: 'bold' }}>Status</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {result.created_problems.map((success) => (
                                                <TableRow key={success.problem_id}>
                                                    <TableCell>{success.index}</TableCell>
                                                    <TableCell sx={{ fontWeight: 'medium' }}>{success.title}</TableCell>
                                                    <TableCell sx={{ fontFamily: 'monospace', fontSize: '0.8rem' }}>{success.slug}</TableCell>
                                                    <TableCell>
                                                        <Chip
                                                            label={success.validation_status}
                                                            size="small"
                                                            color={success.validation_status === 'validated' ? 'success' : 'warning'}
                                                            variant="outlined"
                                                        />
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                </TableContainer>
                            </CardContent>
                        </Card>
                    )}
                </Stack>
            )}
        </Box>
    )
}

export default BulkImport
