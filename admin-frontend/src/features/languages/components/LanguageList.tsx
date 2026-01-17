import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { adminLanguagesApi } from "../../../lib/api/admin"
import type { Language } from "../../../types"
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    TextField,
    Switch,
    IconButton,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TablePagination,
    Typography,
    Chip,
    Tooltip,
    Skeleton,
    Stack
} from "@mui/material"
import {
    Add as AddIcon,
    Edit as EditIcon,
    Delete as DeleteIcon,
    Refresh as RefreshIcon
} from "@mui/icons-material"
import dayjs from "dayjs"
import toast from "react-hot-toast"
import { useState } from "react"
import type { CreateOrUpdateLanguageRequest } from "../../../types/request"

interface LanguageFormValues extends CreateOrUpdateLanguageRequest { }

const initialFormValues: LanguageFormValues = {
    language_id: '',
    name: '',
    version: '',
    extension: '',
    default_template: ''
}

export default function LanguageList() {
    const queryClient = useQueryClient()
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [isDeleteOpen, setIsDeleteOpen] = useState(false)
    const [editingLanguage, setEditingLanguage] = useState<Language | null>(null)
    const [deleteId, setDeleteId] = useState<number | null>(null)
    const [formData, setFormData] = useState<LanguageFormValues>(initialFormValues)

    // Pagination state
    const [page, setPage] = useState(0)
    const [rowsPerPage, setRowsPerPage] = useState(10)

    const { data, isFetching, refetch } = useQuery({
        queryKey: ["admin-languages"],
        queryFn: async () => {
            const response = await adminLanguagesApi.getAll()
            return response.data
        }
    })

    const languages = data?.data || []

    const handleChangePage = (_event: unknown, newPage: number) => {
        setPage(newPage)
    }

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRowsPerPage(parseInt(event.target.value, 10))
        setPage(0)
    }

    // Create mutation
    const createMutation = useMutation({
        mutationFn: (values: CreateOrUpdateLanguageRequest) => adminLanguagesApi.create(values),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language created successfully")
            handleCloseModal()
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error?.message || 'Failed to create language')
        },
    })

    // Update mutation
    const updateMutation = useMutation({
        mutationFn: ({ id, values }: { id: number; values: CreateOrUpdateLanguageRequest }) =>
            adminLanguagesApi.update(id, values),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language updated successfully")
            handleCloseModal()
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error?.message || 'Failed to update language')
        },
    })

    // Status update mutation
    const updateStatusMutation = useMutation({
        mutationFn: ({ languageId, activate }: { languageId: number; activate: boolean }) => {
            return activate ? adminLanguagesApi.deactivate(languageId) : adminLanguagesApi.activate(languageId)
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language status updated successfully")
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to update status')
        },
    })

    // DELETE mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => adminLanguagesApi.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-languages"] })
            toast.success("Language deleted successfully")
            setDeleteId(null)
            setIsDeleteOpen(false)
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error?.message || 'Failed to delete language')
        },
    })

    const handleOpenModal = (language?: Language) => {
        if (language) {
            setFormData({
                language_id: language.language_id,
                name: language.name,
                version: language.version,
                extension: language.extension,
                default_template: language.default_template || "",
            })
            setEditingLanguage(language)
        } else {
            setFormData(initialFormValues)
            setEditingLanguage(null)
        }
        setIsModalOpen(true)
    }

    const handleCloseModal = () => {
        setIsModalOpen(false)
        setEditingLanguage(null)
        setFormData(initialFormValues)
    }

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (editingLanguage) {
            updateMutation.mutate({ id: editingLanguage.id, values: formData })
        } else {
            createMutation.mutate(formData)
        }
    }

    const handleDeleteClick = (id: number) => {
        setDeleteId(id)
        setIsDeleteOpen(true)
    }

    const confirmDelete = () => {
        if (deleteId) {
            deleteMutation.mutate(deleteId)
        }
    }

    const handleInputChange = (field: keyof LanguageFormValues) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({ ...formData, [field]: e.target.value })
    }

    // Pagination logic
    const displayedLanguages = languages.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)

    return (
        <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3, p: 2, bgcolor: 'background.paper', borderRadius: 2, boxShadow: 1 }}>
                <Typography variant="h5" component="h1" fontWeight="bold">
                    Language Management
                </Typography>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => handleOpenModal()}
                    disabled={createMutation.isPending || deleteMutation.isPending}
                >
                    Add Language
                </Button>
            </Box>

            <Paper sx={{ width: '100%', mb: 2, boxShadow: 3, borderRadius: 2, overflow: 'hidden' }}>
                <Box sx={{ p: 2, display: 'flex', alignItems: 'center' }}>
                    <Button
                        variant="outlined"
                        startIcon={<RefreshIcon />}
                        onClick={() => refetch()}
                        disabled={isFetching || deleteMutation.isPending}
                        size="small"
                    >
                        Refresh
                    </Button>
                </Box>
                <TableContainer>
                    <Table sx={{ minWidth: 750 }} aria-labelledby="tableTitle" size="medium">
                        <TableHead>
                            <TableRow sx={{ bgcolor: 'grey.50' }}>
                                <TableCell>ID</TableCell>
                                <TableCell>Language ID</TableCell>
                                <TableCell>Name</TableCell>
                                <TableCell>Status</TableCell>
                                <TableCell>Extension</TableCell>
                                <TableCell>Version</TableCell>
                                <TableCell>Template</TableCell>
                                <TableCell>Created At</TableCell>
                                <TableCell>Updated At</TableCell>
                                <TableCell align="right">Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {isFetching && languages.length === 0 ? (
                                Array.from({ length: 5 }).map((_, index) => (
                                    <TableRow key={`skeleton-${index}`}>
                                        <TableCell><Skeleton variant="text" width={40} /></TableCell>
                                        <TableCell><Skeleton variant="rounded" width={60} height={24} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={100} /></TableCell>
                                        <TableCell><Skeleton variant="rectangular" width={40} height={20} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={60} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={60} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={100} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={120} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={120} /></TableCell>
                                        <TableCell align="right">
                                            <Stack direction="row" spacing={1} justifyContent="flex-end">
                                                <Skeleton variant="circular" width={32} height={32} />
                                                <Skeleton variant="circular" width={32} height={32} />
                                            </Stack>
                                        </TableCell>
                                    </TableRow>
                                ))
                            ) : displayedLanguages.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={10} align="center" sx={{ py: 6 }}>
                                        <Typography variant="body1" color="text.secondary">No languages found</Typography>
                                    </TableCell>
                                </TableRow>
                            ) : (
                                displayedLanguages.map((row: Language) => (
                                    <TableRow
                                        hover
                                        key={row.id}
                                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                    >
                                        <TableCell component="th" scope="row">
                                            <Typography variant="body2" fontFamily="monospace" color="text.secondary">{row.id}</Typography>
                                        </TableCell>
                                        <TableCell>
                                            <Chip label={row.language_id} size="small" variant="outlined" sx={{ fontFamily: 'monospace' }} />
                                        </TableCell>
                                        <TableCell>{row.name}</TableCell>
                                        <TableCell>
                                            <Switch
                                                checked={row.is_active}
                                                onChange={() => updateStatusMutation.mutate({ languageId: row.id, activate: row.is_active })}
                                                disabled={updateStatusMutation.isPending}
                                                color="success"
                                                size="small"
                                            />
                                        </TableCell>
                                        <TableCell>
                                            <Typography variant="body2" fontFamily="monospace">{row.extension}</Typography>
                                        </TableCell>
                                        <TableCell>
                                            <Typography variant="body2" fontFamily="monospace">{row.version}</Typography>
                                        </TableCell>
                                        <TableCell>
                                            <Tooltip title={row.default_template || ''}>
                                                <Typography variant="body2" fontFamily="monospace" noWrap sx={{ maxWidth: 150 }}>
                                                    {row.default_template}
                                                </Typography>
                                            </Tooltip>
                                        </TableCell>
                                        <TableCell>
                                            <Typography variant="caption" color="text.secondary">
                                                {dayjs(row.created_at).format('MMM DD, YYYY HH:mm')}
                                            </Typography>
                                        </TableCell>
                                        <TableCell>
                                            <Typography variant="caption" color="text.secondary">
                                                {dayjs(row.updated_at).format('MMM DD, YYYY HH:mm')}
                                            </Typography>
                                        </TableCell>
                                        <TableCell align="right">
                                            <IconButton
                                                size="small"
                                                onClick={() => handleOpenModal(row)}
                                                disabled={updateStatusMutation.isPending || deleteMutation.isPending}
                                                color="primary"
                                            >
                                                <EditIcon fontSize="small" />
                                            </IconButton>
                                            <IconButton
                                                size="small"
                                                onClick={() => handleDeleteClick(row.id)}
                                                disabled={deleteMutation.isPending}
                                                color="error"
                                            >
                                                <DeleteIcon fontSize="small" />
                                            </IconButton>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </TableContainer>
                <TablePagination
                    rowsPerPageOptions={[5, 10, 25]}
                    component="div"
                    count={languages.length}
                    rowsPerPage={rowsPerPage}
                    page={page}
                    onPageChange={handleChangePage}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                />
            </Paper>

            {/* CREATE/UPDATE DIALOG */}
            <Dialog
                open={isModalOpen}
                onClose={handleCloseModal}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>
                    {editingLanguage ? "Edit Language" : "Create Language"}
                </DialogTitle>
                <form onSubmit={handleSubmit}>
                    <DialogContent dividers>
                        <TextField
                            margin="dense"
                            label="Language ID"
                            fullWidth
                            required
                            placeholder="e.g., en, fr, es"
                            value={formData.language_id}
                            onChange={handleInputChange('language_id')}
                            sx={{ mb: 2 }}
                        />
                        <TextField
                            margin="dense"
                            label="Name"
                            fullWidth
                            required
                            placeholder="Enter language name"
                            value={formData.name}
                            onChange={handleInputChange('name')}
                            sx={{ mb: 2 }}
                        />
                        <TextField
                            margin="dense"
                            label="Version"
                            fullWidth
                            required
                            placeholder="e.g., 1.0.0, latest"
                            value={formData.version}
                            onChange={handleInputChange('version')}
                            sx={{ mb: 2 }}
                        />
                        <TextField
                            margin="dense"
                            label="File Extension"
                            fullWidth
                            required
                            placeholder="e.g., .js, .py, .java"
                            value={formData.extension}
                            onChange={handleInputChange('extension')}
                            sx={{ mb: 2 }}
                        />
                        <TextField
                            margin="dense"
                            label="Default Template"
                            fullWidth
                            required
                            placeholder="Enter default code template"
                            multiline
                            rows={4}
                            value={formData.default_template}
                            onChange={handleInputChange('default_template')}
                        />
                    </DialogContent>
                    <DialogActions sx={{ px: 3, py: 2 }}>
                        <Button onClick={handleCloseModal} color="inherit">
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            variant="contained"
                            disabled={createMutation.isPending || updateMutation.isPending}
                        >
                            {editingLanguage ? 'Update' : 'Create'}
                        </Button>
                    </DialogActions>
                </form>
            </Dialog>

            {/* DELETE CONFIRMATION DIALOG */}
            <Dialog
                open={isDeleteOpen}
                onClose={() => setIsDeleteOpen(false)}
            >
                <DialogTitle>Confirm Delete</DialogTitle>
                <DialogContent>
                    <Typography>
                        Are you sure you want to delete this language? This action cannot be undone.
                    </Typography>
                </DialogContent>
                <DialogActions sx={{ px: 3, py: 2 }}>
                    <Button onClick={() => setIsDeleteOpen(false)} color="inherit">
                        Cancel
                    </Button>
                    <Button
                        onClick={confirmDelete}
                        color="error"
                        variant="contained"
                        disabled={deleteMutation.isPending}
                    >
                        Delete
                    </Button>
                </DialogActions>
            </Dialog>
        </Box>
    )
}
