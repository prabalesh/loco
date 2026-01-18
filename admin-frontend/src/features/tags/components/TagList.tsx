import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { adminProblemApi, adminTagApi } from "../../../lib/api/admin"
import type { Tag } from "../../../types"
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    TextField,
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

interface TagFormValues {
    name: string
    slug: string
}

const initialFormValues: TagFormValues = {
    name: '',
    slug: ''
}

export default function TagList() {
    const queryClient = useQueryClient()
    const [isModalOpen, setIsModalOpen] = useState(false)
    const [isDeleteOpen, setIsDeleteOpen] = useState(false)
    const [editingTag, setEditingTag] = useState<Tag | null>(null)
    const [deleteId, setDeleteId] = useState<number | null>(null)
    const [formData, setFormData] = useState<TagFormValues>(initialFormValues)

    // Pagination state
    const [page, setPage] = useState(0)
    const [rowsPerPage, setRowsPerPage] = useState(10)

    const { data: tagResponse, isFetching, refetch } = useQuery({
        queryKey: ["admin-tags"],
        queryFn: async () => {
            const response = await adminProblemApi.getTags()
            return response.data
        }
    })

    const tags = tagResponse?.data || []

    const handleChangePage = (_event: unknown, newPage: number) => {
        setPage(newPage)
    }

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRowsPerPage(parseInt(event.target.value, 10))
        setPage(0)
    }

    // Create mutation
    const createMutation = useMutation({
        mutationFn: (values: TagFormValues) => adminTagApi.create(values),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-tags"] })
            toast.success("Tag created successfully")
            handleCloseModal()
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to create tag')
        },
    })

    // Update mutation
    const updateMutation = useMutation({
        mutationFn: ({ id, values }: { id: number; values: TagFormValues }) =>
            adminTagApi.update(id, values),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-tags"] })
            toast.success("Tag updated successfully")
            handleCloseModal()
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to update tag')
        },
    })

    // DELETE mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => adminTagApi.delete(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["admin-tags"] })
            toast.success("Tag deleted successfully")
            setDeleteId(null)
            setIsDeleteOpen(false)
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to delete tag')
        },
    })

    const handleOpenModal = (tag?: Tag) => {
        if (tag) {
            setFormData({
                name: tag.name,
                slug: tag.slug,
            })
            setEditingTag(tag)
        } else {
            setFormData(initialFormValues)
            setEditingTag(null)
        }
        setIsModalOpen(true)
    }

    const handleCloseModal = () => {
        setIsModalOpen(false)
        setEditingTag(null)
        setFormData(initialFormValues)
    }

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        if (editingTag) {
            updateMutation.mutate({ id: editingTag.id, values: formData })
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

    const handleInputChange = (field: keyof TagFormValues) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setFormData({ ...formData, [field]: e.target.value })
    }

    // Pagination logic
    const displayedTags = tags.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)

    return (
        <Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3, p: 2, bgcolor: 'background.paper', borderRadius: 2, boxShadow: 1 }}>
                <Typography variant="h5" component="h1" fontWeight="bold">
                    Tag Management
                </Typography>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => handleOpenModal()}
                    disabled={createMutation.isPending || deleteMutation.isPending}
                >
                    Add Tag
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
                    <Table sx={{ minWidth: 650 }} aria-labelledby="tableTitle" size="medium">
                        <TableHead>
                            <TableRow sx={{ bgcolor: 'grey.50' }}>
                                <TableCell>ID</TableCell>
                                <TableCell>Name</TableCell>
                                <TableCell>Slug</TableCell>
                                <TableCell>Created At</TableCell>
                                <TableCell>Updated At</TableCell>
                                <TableCell align="right">Actions</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {isFetching && tags.length === 0 ? (
                                Array.from({ length: 5 }).map((_, index) => (
                                    <TableRow key={`skeleton-${index}`}>
                                        <TableCell><Skeleton variant="text" width={40} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={100} /></TableCell>
                                        <TableCell><Skeleton variant="text" width={120} /></TableCell>
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
                            ) : displayedTags.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={6} align="center" sx={{ py: 6 }}>
                                        <Typography variant="body1" color="text.secondary">No tags found</Typography>
                                    </TableCell>
                                </TableRow>
                            ) : (
                                displayedTags.map((row: Tag) => (
                                    <TableRow
                                        hover
                                        key={row.id}
                                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                                    >
                                        <TableCell component="th" scope="row">
                                            <Typography variant="body2" fontFamily="monospace" color="text.secondary">{row.id}</Typography>
                                        </TableCell>
                                        <TableCell>{row.name}</TableCell>
                                        <TableCell>
                                            <Typography variant="body2" fontFamily="monospace">{row.slug}</Typography>
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
                                                disabled={deleteMutation.isPending}
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
                    count={tags.length}
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
                maxWidth="xs"
                fullWidth
            >
                <DialogTitle>
                    {editingTag ? "Edit Tag" : "Create Tag"}
                </DialogTitle>
                <form onSubmit={handleSubmit}>
                    <DialogContent dividers>
                        <TextField
                            margin="dense"
                            label="Name"
                            fullWidth
                            required
                            placeholder="Enter tag name"
                            value={formData.name}
                            onChange={handleInputChange('name')}
                            sx={{ mb: 2 }}
                        />
                        <TextField
                            margin="dense"
                            label="Slug"
                            fullWidth
                            required
                            placeholder="e.g., array, dynamic-programming"
                            value={formData.slug}
                            onChange={handleInputChange('slug')}
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
                            {editingTag ? 'Update' : 'Create'}
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
                        Are you sure you want to delete this tag? This action cannot be undone.
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
