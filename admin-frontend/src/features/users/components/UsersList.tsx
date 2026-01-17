import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  TableSortLabel,
  Paper,
  Button,
  IconButton,
  Switch,
  Chip,
  Select,
  MenuItem,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Tooltip,
  Box,
  Typography,
  CircularProgress,
  Avatar,
  Card,
  Stack,
  alpha,
  Zoom,
  Skeleton,
} from '@mui/material'
import {
  Delete as DeleteIcon,
  Edit as EditIcon,
  Check as CheckIcon,
  Close as CloseIcon,
  Refresh as RefreshIcon,
  PersonOutline as PersonIcon,
  VerifiedUser as VerifiedIcon,
  Warning as WarningIcon,
} from '@mui/icons-material'
import { adminUsersApi } from '../../../lib/api/admin'
import toast from 'react-hot-toast'
import dayjs from 'dayjs'
import type { User } from '../../../types'

type Order = 'asc' | 'desc'

export const UsersList = () => {
  const queryClient = useQueryClient()
  const [editingRole, setEditingRole] = useState<{ userId: number; role: string } | null>(null)
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; userId: number | null }>({
    open: false,
    userId: null,
  })
  const [page, setPage] = useState(0)
  const [rowsPerPage, setRowsPerPage] = useState(10)
  const [order, setOrder] = useState<Order>('asc')
  const [orderBy, setOrderBy] = useState<keyof User>('id')

  const { data, isFetching, refetch } = useQuery({
    queryKey: ['admin-users'],
    queryFn: async () => {
      const response = await adminUsersApi.getAll()
      return response.data
    },
  })

  const users: User[] = data?.data ?? []

  const deleteMutation = useMutation({
    mutationFn: (userId: number) => adminUsersApi.deleteUser(userId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      toast.success('User deleted successfully')
      setDeleteDialog({ open: false, userId: null })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to delete user')
    },
  })

  const updateRoleMutation = useMutation({
    mutationFn: ({ userId, role }: { userId: number; role: string }) =>
      adminUsersApi.updateRole(userId, role),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      toast.success('Role updated successfully')
      setEditingRole(null)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update role')
    },
  })

  const updateStatusMutation = useMutation({
    mutationFn: ({ userId, isActive }: { userId: number; isActive: boolean }) =>
      adminUsersApi.updateStatus(userId, isActive),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] })
      toast.success('Status updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update status')
    },
  })

  const roleConfig: Record<string, { color: 'error' | 'warning' | 'primary'; bgcolor: string }> = {
    admin: { color: 'error', bgcolor: '#eb3131ff' },
    moderator: { color: 'warning', bgcolor: '#49eb43ff' },
    user: { color: 'primary', bgcolor: '#3186f7ff' },
  }

  const handleRequestSort = (property: keyof User) => {
    const isAsc = orderBy === property && order === 'asc'
    setOrder(isAsc ? 'desc' : 'asc')
    setOrderBy(property)
  }

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage)
  }

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10))
    setPage(0)
  }

  const handleDeleteClick = (userId: number) => {
    setDeleteDialog({ open: true, userId })
  }

  const handleDeleteConfirm = () => {
    if (deleteDialog.userId) {
      deleteMutation.mutate(deleteDialog.userId)
    }
  }

  const descendingComparator = <T,>(a: T, b: T, orderBy: keyof T) => {
    if (b[orderBy] < a[orderBy]) return -1
    if (b[orderBy] > a[orderBy]) return 1
    return 0
  }

  const getComparator = <Key extends keyof any>(
    order: Order,
    orderBy: Key
  ): ((a: { [key in Key]: any }, b: { [key in Key]: any }) => number) => {
    return order === 'desc'
      ? (a, b) => descendingComparator(a, b, orderBy)
      : (a, b) => -descendingComparator(a, b, orderBy)
  }

  const sortedUsers = [...users].sort(getComparator(order, orderBy))
  const paginatedUsers = sortedUsers.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  )

  return (
    <Box sx={{ maxWidth: '100%', px: { xs: 2, md: 4 }, py: 4 }}>
      {/* Header Section */}
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={4}>
        <Box>
          <Typography
            variant="h4"
            sx={{
              fontWeight: 700,
              color: 'text.primary',
              mb: 0.5,
              background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}
          >
            User Management
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Manage user roles, permissions and account status
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={isFetching ? <CircularProgress size={16} color="inherit" /> : <RefreshIcon />}
          onClick={() => refetch()}
          disabled={isFetching}
          sx={{
            borderRadius: 2,
            textTransform: 'none',
            px: 3,
            py: 1.5,
            boxShadow: 2,
            '&:hover': {
              boxShadow: 4,
              transform: 'translateY(-2px)',
            },
            transition: 'all 0.3s ease',
          }}
        >
          Refresh
        </Button>
      </Stack>

      {/* Stats Cards */}
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} mb={3}>
        <Card
          sx={{
            flex: 1,
            p: 2.5,
            borderRadius: 3,
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            color: 'white',
            boxShadow: 3,
          }}
        >
          <Typography variant="body2" sx={{ opacity: 0.9, mb: 0.5 }}>Total Users</Typography>
          <Typography variant="h4" fontWeight={700}>{users.length}</Typography>
        </Card>
        <Card
          sx={{
            flex: 1,
            p: 2.5,
            borderRadius: 3,
            background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
            color: 'white',
            boxShadow: 3,
          }}
        >
          <Typography variant="body2" sx={{ opacity: 0.9, mb: 0.5 }}>Active Users</Typography>
          <Typography variant="h4" fontWeight={700}>
            {users.filter(u => u.is_active).length}
          </Typography>
        </Card>
        <Card
          sx={{
            flex: 1,
            p: 2.5,
            borderRadius: 3,
            background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
            color: 'white',
            boxShadow: 3,
          }}
        >
          <Typography variant="body2" sx={{ opacity: 0.9, mb: 0.5 }}>Verified</Typography>
          <Typography variant="h4" fontWeight={700}>
            {users.filter(u => u.email_verified).length}
          </Typography>
        </Card>
      </Stack>

      {/* Table */}
      <TableContainer
        component={Paper}
        elevation={0}
        sx={{
          borderRadius: 3,
          border: '1px solid',
          borderColor: 'divider',
          overflow: 'hidden',
        }}
      >
        <Table sx={{ minWidth: 900 }}>
          <TableHead>
            <TableRow
              sx={{
                bgcolor: (theme) => alpha(theme.palette.primary.main, 0.05),
              }}
            >
              <TableCell
                sortDirection={orderBy === 'id' ? order : false}
                sx={{ fontWeight: 700, fontSize: '0.875rem' }}
              >
                <TableSortLabel
                  active={orderBy === 'id'}
                  direction={orderBy === 'id' ? order : 'asc'}
                  onClick={() => handleRequestSort('id')}
                >
                  ID
                </TableSortLabel>
              </TableCell>
              <TableCell
                sortDirection={orderBy === 'username' ? order : false}
                sx={{ fontWeight: 700, fontSize: '0.875rem' }}
              >
                <TableSortLabel
                  active={orderBy === 'username'}
                  direction={orderBy === 'username' ? order : 'asc'}
                  onClick={() => handleRequestSort('username')}
                >
                  User
                </TableSortLabel>
              </TableCell>
              <TableCell sx={{ fontWeight: 700, fontSize: '0.875rem' }}>Email</TableCell>
              <TableCell sx={{ fontWeight: 700, fontSize: '0.875rem' }}>Role</TableCell>
              <TableCell sx={{ fontWeight: 700, fontSize: '0.875rem' }}>Status</TableCell>
              <TableCell sx={{ fontWeight: 700, fontSize: '0.875rem' }}>Verification</TableCell>
              <TableCell
                sortDirection={orderBy === 'created_at' ? order : false}
                sx={{ fontWeight: 700, fontSize: '0.875rem' }}
              >
                <TableSortLabel
                  active={orderBy === 'created_at'}
                  direction={orderBy === 'created_at' ? order : 'asc'}
                  onClick={() => handleRequestSort('created_at')}
                >
                  Joined
                </TableSortLabel>
              </TableCell>
              <TableCell align="center" sx={{ fontWeight: 700, fontSize: '0.875rem' }}>
                Actions
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {isFetching && users.length === 0 ? (
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={`skeleton-${index}`}>
                  <TableCell><Skeleton variant="text" width={40} /></TableCell>
                  <TableCell>
                    <Stack direction="row" spacing={1.5} alignItems="center">
                      <Skeleton variant="circular" width={36} height={36} />
                      <Skeleton variant="text" width={100} />
                    </Stack>
                  </TableCell>
                  <TableCell><Skeleton variant="text" width={150} /></TableCell>
                  <TableCell><Skeleton variant="rounded" width={80} height={24} /></TableCell>
                  <TableCell><Skeleton variant="rectangular" width={40} height={20} /></TableCell>
                  <TableCell><Skeleton variant="rounded" width={80} height={24} /></TableCell>
                  <TableCell><Skeleton variant="text" width={100} /></TableCell>
                  <TableCell align="center"><Skeleton variant="circular" width={32} height={32} /></TableCell>
                </TableRow>
              ))
            ) : paginatedUsers.length === 0 ? (
              <TableRow>
                <TableCell colSpan={8} align="center" sx={{ py: 8 }}>
                  <PersonIcon sx={{ fontSize: 60, color: 'text.disabled', mb: 2 }} />
                  <Typography variant="body1" color="text.secondary">
                    No users found
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              paginatedUsers.map((user) => {
                const isEditing = editingRole?.userId === user.id
                return (
                  <TableRow
                    key={user.id}
                    hover
                    sx={{
                      '&:hover': {
                        bgcolor: (theme) => alpha(theme.palette.primary.main, 0.02),
                      },
                      transition: 'background-color 0.2s ease',
                    }}
                  >
                    <TableCell sx={{ fontFamily: 'monospace', color: 'text.secondary', fontWeight: 600 }}>
                      #{user.id}
                    </TableCell>
                    <TableCell>
                      <Stack direction="row" spacing={1.5} alignItems="center">
                        <Avatar
                          sx={{
                            width: 36,
                            height: 36,
                            bgcolor: 'primary.main',
                            fontSize: '0.875rem',
                            fontWeight: 600,
                          }}
                        >
                          {user.username.charAt(0).toUpperCase()}
                        </Avatar>
                        <Typography fontWeight={600} fontSize="0.9rem">
                          {user.username}
                        </Typography>
                      </Stack>
                    </TableCell>
                    <TableCell>
                      <Tooltip title={user.email} arrow placement="top">
                        <Typography
                          sx={{
                            maxWidth: 200,
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                            color: 'text.secondary',
                            fontSize: '0.875rem',
                          }}
                        >
                          {user.email}
                        </Typography>
                      </Tooltip>
                    </TableCell>
                    <TableCell>
                      {isEditing ? (
                        <Stack direction="row" spacing={0.5} alignItems="center">
                          <Select
                            value={editingRole.role}
                            onChange={(e) =>
                              setEditingRole({ userId: user.id, role: e.target.value })
                            }
                            size="small"
                            disabled={updateRoleMutation.isPending}
                            sx={{
                              minWidth: 120,
                              borderRadius: 1.5,
                              '& .MuiOutlinedInput-notchedOutline': {
                                borderColor: 'primary.main',
                              },
                            }}
                          >
                            <MenuItem value="user">User</MenuItem>
                            <MenuItem value="admin">Admin</MenuItem>
                            <MenuItem value="moderator">Moderator</MenuItem>
                          </Select>
                          <IconButton
                            color="success"
                            size="small"
                            onClick={() =>
                              updateRoleMutation.mutate({
                                userId: user.id,
                                role: editingRole.role,
                              })
                            }
                            disabled={updateRoleMutation.isPending}
                            sx={{
                              bgcolor: (theme) => alpha(theme.palette.success.main, 0.1),
                              '&:hover': { bgcolor: (theme) => alpha(theme.palette.success.main, 0.2) },
                            }}
                          >
                            <CheckIcon fontSize="small" />
                          </IconButton>
                          <IconButton
                            size="small"
                            onClick={() => setEditingRole(null)}
                            disabled={updateRoleMutation.isPending}
                            sx={{
                              bgcolor: (theme) => alpha(theme.palette.grey[500], 0.1),
                              '&:hover': { bgcolor: (theme) => alpha(theme.palette.grey[500], 0.2) },
                            }}
                          >
                            <CloseIcon fontSize="small" />
                          </IconButton>
                        </Stack>
                      ) : (
                        <Stack direction="row" spacing={0.5} alignItems="center">
                          <Chip
                            label={user.role.toUpperCase()}
                            color={roleConfig[user.role]?.color || 'default'}
                            size="small"
                            sx={{
                              fontWeight: 700,
                              fontSize: '0.7rem',
                              letterSpacing: 0.5,
                              px: 1,
                              bgcolor: roleConfig[user.role]?.bgcolor,
                            }}
                          />
                          <IconButton
                            size="small"
                            onClick={() => setEditingRole({ userId: user.id, role: user.role })}
                            sx={{
                              opacity: 0.6,
                              '&:hover': { opacity: 1 },
                            }}
                          >
                            <EditIcon fontSize="small" />
                          </IconButton>
                        </Stack>
                      )}
                    </TableCell>
                    <TableCell>
                      <Switch
                        checked={user.is_active}
                        onChange={(e) =>
                          updateStatusMutation.mutate({
                            userId: user.id,
                            isActive: e.target.checked,
                          })
                        }
                        disabled={updateStatusMutation.isPending}
                        sx={{
                          '& .MuiSwitch-switchBase.Mui-checked': {
                            color: 'success.main',
                          },
                          '& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track': {
                            bgcolor: 'success.main',
                          },
                        }}
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        icon={user.email_verified ? <VerifiedIcon /> : <WarningIcon />}
                        label={user.email_verified ? 'VERIFIED' : 'PENDING'}
                        color={user.email_verified ? 'success' : 'warning'}
                        variant={user.email_verified ? 'filled' : 'outlined'}
                        size="small"
                        sx={{
                          fontWeight: 600,
                          fontSize: '0.7rem',
                          letterSpacing: 0.5,
                        }}
                      />
                    </TableCell>
                    <TableCell>
                      <Stack spacing={0.25}>
                        <Typography variant="body2" fontWeight={600}>
                          {dayjs(user.created_at).format('MMM DD, YYYY')}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          {dayjs(user.created_at).format('HH:mm')}
                        </Typography>
                      </Stack>
                    </TableCell>
                    <TableCell align="center">
                      <Tooltip title="Delete user" arrow>
                        <IconButton
                          color="error"
                          size="small"
                          onClick={() => handleDeleteClick(user.id)}
                          disabled={deleteMutation.isPending}
                          sx={{
                            bgcolor: (theme) => alpha(theme.palette.error.main, 0.08),
                            '&:hover': {
                              bgcolor: (theme) => alpha(theme.palette.error.main, 0.15),
                            },
                          }}
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25, 50]}
          component="div"
          count={users.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          sx={{
            borderTop: '1px solid',
            borderColor: 'divider',
            bgcolor: (theme) => alpha(theme.palette.grey[50], 0.5),
          }}
        />
      </TableContainer>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialog.open}
        onClose={() => setDeleteDialog({ open: false, userId: null })}
        TransitionComponent={Zoom}
        PaperProps={{
          sx: {
            borderRadius: 3,
            px: 1,
            py: 1,
          },
        }}
      >
        <DialogTitle sx={{ fontWeight: 700, fontSize: '1.25rem' }}>
          Delete User
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete this user? This action cannot be undone and will permanently remove all user data.
          </DialogContentText>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button
            onClick={() => setDeleteDialog({ open: false, userId: null })}
            disabled={deleteMutation.isPending}
            sx={{ borderRadius: 2, textTransform: 'none', px: 3 }}
          >
            Cancel
          </Button>
          <Button
            onClick={handleDeleteConfirm}
            color="error"
            variant="contained"
            disabled={deleteMutation.isPending}
            startIcon={deleteMutation.isPending ? <CircularProgress size={16} color="inherit" /> : <DeleteIcon />}
            sx={{
              borderRadius: 2,
              textTransform: 'none',
              px: 3,
              boxShadow: 2,
            }}
          >
            {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
