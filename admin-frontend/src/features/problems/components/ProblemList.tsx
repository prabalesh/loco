import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Typography,
  TablePagination,
  IconButton,
  Stack,
  Box,
  useTheme,
  alpha,
  Skeleton,
} from "@mui/material";
import {
  Add as PlusOutlined,
  Edit as EditOutlined,
  Delete as DeleteOutlined,
} from "@mui/icons-material";
import dayjs from "dayjs";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";
import { adminProblemApi } from "../../../lib/api/admin";
import { PROBLEM_STEPS, ROUTES } from "../../../config/constant";
import { Play } from "lucide-react";
import { useState, useMemo } from "react";

const DIFFICULTY_COLORS = {
  easy: "success",
  medium: "warning",
  hard: "error",
} as const;

const STATUS_COLORS = {
  published: "info",
  draft: "default",
} as const;

const getResumeLink = (step: number, problemId: number) => {
  const stepRoutes = {
    1: ROUTES.PROBLEMS.TESTCASES,
    2: ROUTES.PROBLEMS.LANGUAGES,
    3: ROUTES.PROBLEMS.VALIDATE,
    4: ROUTES.PROBLEMS.VALIDATE,
  };
  return stepRoutes[step as keyof typeof stepRoutes]?.(problemId) || "";
};

const TableSkeleton = () => (
  <>
    {Array.from({ length: 5 }).map((_, index) => (
      <TableRow key={`skeleton-${index}`}>
        <TableCell>
          <Skeleton variant="text" width="60%" height={24} />
          <Skeleton variant="text" width="40%" height={16} />
        </TableCell>
        <TableCell><Skeleton variant="rounded" width={80} height={24} /></TableCell>
        <TableCell><Skeleton variant="rounded" width={80} height={24} /></TableCell>
        <TableCell><Skeleton variant="rounded" width={80} height={24} /></TableCell>
        <TableCell><Skeleton variant="rounded" width={40} height={24} /></TableCell>
        <TableCell><Skeleton variant="text" width={100} /></TableCell>
        <TableCell>
          <Stack direction="row" spacing={1}>
            <Skeleton variant="rectangular" width={60} height={32} />
            <Skeleton variant="rectangular" width={60} height={32} />
            <Skeleton variant="circular" width={32} height={32} />
          </Stack>
        </TableCell>
      </TableRow>
    ))}
  </>
);

const columns = [
  { id: "title", label: "Title", minWidth: 150 },
  { id: "difficulty", label: "Difficulty", minWidth: 120 },
  { id: "status", label: "Status", minWidth: 120 },
  { id: "current_step", label: "Current Step", minWidth: 120 },
  { id: "is_active", label: "Active", minWidth: 100 },
  { id: "updated_at", label: "Updated", minWidth: 150 },
  { id: "actions", label: "Actions", minWidth: 150 },
];

export default function ProblemList() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const theme = useTheme();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [deleteId, setDeleteId] = useState<number | null>(null);

  const { data, isFetching } = useQuery({
    queryKey: ["admin-problems"],
    queryFn: async () => {
      const res = await adminProblemApi.getAll();
      return res.data;
    },
  });

  const problems = data?.data ?? [];

  const paginatedProblems = useMemo(
    () => problems.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage),
    [problems, page, rowsPerPage]
  );

  const problemToDelete = useMemo(
    () => problems.find((p) => p.id === deleteId),
    [problems, deleteId]
  );

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminProblemApi.delete(String(id)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-problems"] });
      toast.success("Problem deleted successfully");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Failed to delete problem");
    },
  });

  const handleChangePage = (_: unknown, newPage: number) => setPage(newPage);

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleDelete = () => {
    if (deleteId !== null) {
      deleteMutation.mutate(deleteId);
    }
    setDeleteId(null);
  };

  return (
    <Box sx={{ p: 3, bgcolor: "background.default" }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center" mb={3}>
        <div>
          <Typography variant="h4" component="h1" gutterBottom fontWeight="bold" color="text.primary">
            Problem Management
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Manage coding problems for your platform
          </Typography>
        </div>
        <Button
          variant="contained"
          color="primary"
          startIcon={<PlusOutlined />}
          onClick={() => navigate("/problems/create")}
          sx={{
            borderRadius: 2,
            boxShadow: 1,
            "&:hover": { boxShadow: 2 },
          }}
        >
          Create Problem
        </Button>
      </Stack>

      <Paper sx={{ borderRadius: 3, boxShadow: 2, overflow: "hidden", bgcolor: "background.paper" }}>
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                {columns.map((column) => (
                  <TableCell
                    key={column.id}
                    sx={{
                      fontWeight: "bold",
                      color: "text.primary",
                      bgcolor: "background.default",
                      borderBottom: `2px solid ${alpha(theme.palette.primary.main, 0.2)}`,
                    }}
                  >
                    {column.label}
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {isFetching && problems.length === 0 ? (
                <TableSkeleton />
              ) : (
                paginatedProblems.map((record) => (
                  <TableRow
                    key={record.id}
                    hover
                    sx={{
                      "&:nth-of-type(odd)": { bgcolor: alpha(theme.palette.primary.main, 0.02) },
                      "&:hover": { bgcolor: alpha(theme.palette.primary.main, 0.08) },
                    }}
                  >
                    <TableCell>
                      <Typography variant="subtitle1" fontWeight="bold">
                        {record.title}
                      </Typography>
                      <Typography variant="caption" color="textSecondary">
                        {record.slug}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={record.difficulty.toUpperCase()}
                        color={DIFFICULTY_COLORS[record.difficulty as keyof typeof DIFFICULTY_COLORS]}
                        size="small"
                        sx={{ borderRadius: 1 }}
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={record.status.toUpperCase()}
                        color={STATUS_COLORS[record.status as keyof typeof STATUS_COLORS] || "default"}
                        size="small"
                        sx={{ borderRadius: 1 }}
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={PROBLEM_STEPS[record.current_step - 1]?.label || ""}
                        size="small"
                        sx={{ borderRadius: 1 }}
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={record.is_active ? "YES" : "NO"}
                        color={record.is_active ? "success" : "error"}
                        size="small"
                        sx={{ borderRadius: 1 }}
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="caption" color="textSecondary">
                        {dayjs(record.updated_at).format("MMM DD, YYYY")}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Stack direction="row" spacing={1}>
                        <Button
                          variant="outlined"
                          size="small"
                          startIcon={<EditOutlined />}
                          onClick={() => navigate(`/problems/edit/${record.id}`)}
                          sx={{
                            borderRadius: 1,
                            borderColor: "primary.main",
                            color: "primary.main",
                            "&:hover": { bgcolor: alpha(theme.palette.primary.main, 0.1) },
                          }}
                        >
                          Edit
                        </Button>
                        <Button
                          variant="outlined"
                          size="small"
                          startIcon={<Play style={{ width: 16, height: 16 }} />}
                          onClick={() => navigate(getResumeLink(record.current_step, record.id))}
                          sx={{
                            borderRadius: 1,
                            borderColor: "secondary.main",
                            color: "secondary.main",
                            "&:hover": { bgcolor: alpha(theme.palette.secondary.main, 0.1) },
                          }}
                        >
                          Resume
                        </Button>
                        <IconButton
                          color="error"
                          size="small"
                          onClick={() => setDeleteId(record.id)}
                          disabled={deleteMutation.isPending}
                          sx={{
                            borderRadius: 1,
                            "&:hover": { bgcolor: alpha(theme.palette.error.main, 0.1) },
                          }}
                        >
                          <DeleteOutlined />
                        </IconButton>
                      </Stack>
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
          count={problems?.length ?? 0}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          labelRowsPerPage="Rows per page"
          labelDisplayedRows={({ from, to, count }) => `${from}-${to} of ${count}`}
          sx={{
            borderTop: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
            bgcolor: "background.default",
          }}
        />
      </Paper>

      <Dialog open={deleteId !== null} onClose={() => setDeleteId(null)}>
        <DialogTitle sx={{ fontWeight: "bold" }}>Delete Problem</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete <strong>"{problemToDelete?.title}"</strong>?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteId(null)} color="inherit">
            Cancel
          </Button>
          <Button onClick={handleDelete} color="error" variant="contained">
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
