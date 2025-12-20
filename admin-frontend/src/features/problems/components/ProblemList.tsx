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
} from "@mui/material";
import {
  Add as PlusOutlined,
  Edit as EditOutlined,
  Delete as DeleteOutlined,
} from "@mui/icons-material";
import dayjs from "dayjs";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";
import { adminProblemApi } from "../../../api/adminApi";
import { PROBLEM_STEPS, ROUTES } from "../../../config/constant";
import { Play } from "lucide-react";
import { useState } from "react";

export default function ProblemList() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const theme = useTheme();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [confirmDelete, setConfirmDelete] = useState<{ open: boolean; id: number | null; title: string }>({
    open: false,
    id: null,
    title: "",
  });

  const { data } = useQuery({
    queryKey: ["admin-problems"],
    queryFn: async () => {
      const res = await adminProblemApi.getAll();
      return res.data;
    },
  });

  const problems = data?.data || [];

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminProblemApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin-problems"] });
      toast.success("Problem deleted successfully");
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error?.message || "Failed to delete problem");
    },
  });

  const handleChangePage = (_: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleConfirmDelete = (id: number, title: string) => {
    setConfirmDelete({ open: true, id, title });
  };

  const handleCloseConfirm = () => {
    setConfirmDelete({ open: false, id: null, title: "" });
  };

  const handleDelete = () => {
    if (confirmDelete.id !== null) {
      deleteMutation.mutate(confirmDelete.id);
    }
    handleCloseConfirm();
  };

  const columns = [
    { id: "title", label: "Title", minWidth: 150 },
    { id: "difficulty", label: "Difficulty", minWidth: 120 },
    { id: "status", label: "Status", minWidth: 120 },
    { id: "current_step", label: "Current Step", minWidth: 120 },
    { id: "is_active", label: "Active", minWidth: 100 },
    { id: "updated_at", label: "Updated", minWidth: 150 },
    { id: "actions", label: "Actions", minWidth: 150 },
  ];

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
            "&:hover": {
              boxShadow: 2,
            },
          }}
        >
          Create Problem
        </Button>
      </Stack>

      <Paper
        sx={{
          borderRadius: 3,
          boxShadow: 2,
          overflow: "hidden",
          bgcolor: "background.paper",
        }}
      >
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
              {problems.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((record, _) => (
                <TableRow
                  key={record.id}
                  hover
                  sx={{
                    "&:nth-of-type(odd)": {
                      bgcolor: alpha(theme.palette.primary.main, 0.02),
                    },
                    "&:hover": {
                      bgcolor: alpha(theme.palette.primary.main, 0.08),
                    },
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
                      color={
                        record.difficulty === "easy"
                          ? "success"
                          : record.difficulty === "medium"
                          ? "warning"
                          : "error"
                      }
                      size="small"
                      sx={{ borderRadius: 1 }}
                    />
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={record.status.toUpperCase()}
                      color={record.status === "published" ? "info" : "default"}
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
                          "&:hover": {
                            bgcolor: alpha(theme.palette.primary.main, 0.1),
                          },
                        }}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="outlined"
                        size="small"
                        startIcon={<Play style={{ width: 16, height: 16 }} />}
                        onClick={() => {
                          let link = "";
                          switch (record.current_step) {
                            case 1:
                              link = ROUTES.PROBLEMS.TESTCASES(record.id);
                              break;
                            case 2:
                              link = ROUTES.PROBLEMS.LANGUAGES(record.id);
                              break;
                            case 3:
                            case 4:
                              link = ROUTES.PROBLEMS.VALIDATE(record.id);
                              break;
                          }
                          navigate(link);
                        }}
                        sx={{
                          borderRadius: 1,
                          borderColor: "secondary.main",
                          color: "secondary.main",
                          "&:hover": {
                            bgcolor: alpha(theme.palette.secondary.main, 0.1),
                          },
                        }}
                      >
                        Resume
                      </Button>
                      <IconButton
                        color="error"
                        size="small"
                        onClick={() => handleConfirmDelete(record.id, record.title)}
                        disabled={deleteMutation.isPending}
                        sx={{
                          borderRadius: 1,
                          "&:hover": {
                            bgcolor: alpha(theme.palette.error.main, 0.1),
                          },
                        }}
                      >
                        <DeleteOutlined />
                      </IconButton>
                    </Stack>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={data?.total || 0}
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

      <Dialog open={confirmDelete.open} onClose={handleCloseConfirm}>
        <DialogTitle sx={{ fontWeight: "bold" }}>Delete Problem</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete <strong>"{confirmDelete.title}"</strong>?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseConfirm} color="inherit">
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
