import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Switch,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
  Paper,
  FormControlLabel,
  DialogContentText,
  Stack,
} from "@mui/material";
import {
  Add as PlusOutlined,
  Edit as EditOutlined,
  Delete as DeleteOutlined,
  Visibility as EyeOutlined,
  VisibilityOff as EyeInvisibleOutlined,
} from "@mui/icons-material";
import { useState } from "react";
import toast from "react-hot-toast";
import type { CreateTestCaseRequest } from "../../../types/request";
import type { TestCase } from "../../../types";
import { adminTestcaseApi } from "../../../api/adminApi";

interface TestCaseFormValues extends CreateTestCaseRequest {}

export interface TestCaseListProps {
  problemId: number;
}

export default function TestCaseList({ problemId }: TestCaseListProps) {
  const queryClient = useQueryClient();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTestCase, setEditingTestCase] = useState<TestCase | null>(null);
  const [formValues, setFormValues] = useState<TestCaseFormValues>({
    problem_id: problemId,
    input: "",
    expected_output: "",
    is_sample: false
  });
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; id: number | null }>({
    open: false,
    id: null,
  });

  const { data, isFetching } = useQuery({
    queryKey: ["testcases", problemId],
    queryFn: () => adminTestcaseApi.getAll(problemId),
  });

  const testCases = data?.data.data || [];

  const createMutation = useMutation({
    mutationFn: (values: CreateTestCaseRequest) => adminTestcaseApi.create(problemId, values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["testcases", problemId] });
      toast.success("Test case created");
      closeModal();
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ testcaseId, values }: { testcaseId: number; values: CreateTestCaseRequest }) =>
      adminTestcaseApi.update(testcaseId, values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["testcases", problemId] });
      toast.success("Test case updated");
      closeModal();
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (testcaseId: number) => adminTestcaseApi.delete(problemId, testcaseId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["testcases", problemId] });
      toast.success("Test case deleted");
      setDeleteDialog({ open: false, id: null });
    },
  });

  const openModal = (testcase?: TestCase) => {
    if (testcase) {
      setEditingTestCase(testcase);
      setFormValues({
        problem_id: problemId,
        input: testcase.input,
        expected_output: testcase.expected_output,
        is_sample: testcase.is_sample,
        is_hidden: testcase.is_hidden,
      });
    } else {
      setEditingTestCase(null);
      setFormValues({
        problem_id: problemId,
        input: "",
        expected_output: "",
        is_sample: false,
        is_hidden: false,
      });
    }
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setEditingTestCase(null);
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    const values: TestCaseFormValues = {
      ...formValues,
      problem_id: problemId,
    };

    if (!values.input || !values.expected_output) {
      toast.error("Input and Expected Output are required");
      return;
    }

    if (editingTestCase) {
      updateMutation.mutate({ testcaseId: editingTestCase.id, values });
    } else {
      createMutation.mutate(values);
    }
  };

  const handleChangeField = (field: keyof TestCaseFormValues, value: any) => {
    setFormValues((prev) => ({ ...prev, [field]: value }));
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h6" fontWeight="bold">
          Test Cases ({testCases.length})
        </Typography>
        <Button
          variant="contained"
          startIcon={<PlusOutlined />}
          onClick={() => openModal()}
          disabled={createMutation.isPending}
        >
          Add Test Case
        </Button>
      </Box>

      <Paper variant="outlined">
        <TableContainer sx={{ maxHeight: 400 }}>
          <Table size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Order</TableCell>
                <TableCell>Input</TableCell>
                <TableCell>Expected Output</TableCell>
                <TableCell>Sample</TableCell>
                <TableCell>Hidden</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {isFetching ? (
                <TableRow>
                  <TableCell colSpan={7}>Loading...</TableCell>
                </TableRow>
              ) : testCases.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7}>
                    <Typography variant="body2" color="text.secondary">
                      No test cases yet.
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                testCases.map((record: TestCase) => (
                  <TableRow key={record.id} hover>
                    <TableCell>{record.id}</TableCell>
                    <TableCell>{record.order}</TableCell>
                    <TableCell>
                      <Box
                        sx={{
                          maxWidth: 200,
                          bgcolor: "grey.50",
                          p: 1,
                          borderRadius: 1,
                          fontFamily: "monospace",
                          fontSize: "0.75rem",
                          whiteSpace: "nowrap",
                          overflow: "hidden",
                          textOverflow: "ellipsis",
                        }}
                      >
                        {record.input}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box
                        sx={{
                          maxWidth: 200,
                          bgcolor: "success.50",
                          p: 1,
                          borderRadius: 1,
                          fontFamily: "monospace",
                          fontSize: "0.75rem",
                          whiteSpace: "nowrap",
                          overflow: "hidden",
                          textOverflow: "ellipsis",
                        }}
                      >
                        {record.expected_output}
                      </Box>
                    </TableCell>
                    <TableCell>
                      {record.is_sample ? (
                        <Chip label="PUBLIC" size="small" />
                      ) : (
                        <Chip label="PRIVATE" size="small" />
                      )}
                    </TableCell>
                    <TableCell align="right">
                      <Stack direction="row" spacing={1} justifyContent="flex-end">
                        <IconButton
                          size="small"
                          color="primary"
                          onClick={() => openModal(record)}
                        >
                          <EditOutlined fontSize="small" />
                        </IconButton>
                        <IconButton
                          size="small"
                          color="error"
                          onClick={() =>
                            setDeleteDialog({ open: true, id: record.id })
                          }
                        >
                          <DeleteOutlined fontSize="small" />
                        </IconButton>
                      </Stack>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Create / Edit Dialog */}
      <Dialog open={isModalOpen} onClose={closeModal} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingTestCase ? "Edit Test Case" : "Create Test Case"}
        </DialogTitle>
        <Box component="form" onSubmit={handleSubmit}>
          <DialogContent>
            <TextField
              label="Input"
              value={formValues.input}
              onChange={(e) => handleChangeField("input", e.target.value)}
              multiline
              minRows={4}
              fullWidth
              margin="normal"
              required
              placeholder="Enter test case input"
            />
            <TextField
              label="Expected Output"
              value={formValues.expected_output}
              onChange={(e) =>
                handleChangeField("expected_output", e.target.value)
              }
              multiline
              minRows={4}
              fullWidth
              margin="normal"
              required
              placeholder="Enter expected output"
            />
            <FormControlLabel
              control={
                <Switch
                  checked={formValues.is_sample}
                  onChange={(e) =>
                    handleChangeField("is_sample", e.target.checked)
                  }
                />
              }
              label="Is Sample Test Case"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={closeModal}>Cancel</Button>
            <Button
              type="submit"
              variant="contained"
              disabled={createMutation.isPending || updateMutation.isPending}
            >
              {editingTestCase ? "Update" : "Create"}
            </Button>
          </DialogActions>
        </Box>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialog.open}
        onClose={() => setDeleteDialog({ open: false, id: null })}
      >
        <DialogTitle>Delete Test Case</DialogTitle>
        <DialogContent>
          <DialogContentText>Are you sure?</DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => setDeleteDialog({ open: false, id: null })}
            color="inherit"
          >
            Cancel
          </Button>
          <Button
            onClick={() => {
              if (deleteDialog.id != null) {
                deleteMutation.mutate(deleteDialog.id);
              }
            }}
            color="error"
            variant="contained"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
