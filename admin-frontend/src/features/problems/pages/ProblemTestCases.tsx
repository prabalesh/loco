import { useNavigate, useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { CircularProgress, Alert, Button, Box, Typography, Paper, Stack } from "@mui/material";
import { adminProblemApi } from "../../../lib/api/admin";
import { ProblemStepper } from "../components/ProblemStepper";
import toast from "react-hot-toast";
import TestCaseList from "../components/TestCasetList";

export default function ProblemTestCases() {
  const navigate = useNavigate();
  const { problemId } = useParams<{ problemId: string }>();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["problem", problemId],
    queryFn: () => adminProblemApi.getById(String(problemId)),
  });

  const problem = data?.data || null;

  // Mutation to validate test cases
  const validateMutation = useMutation({
    mutationFn: () => {
      return adminProblemApi.validateTestCases(Number(problemId));
    },
    onSuccess: () => {
      toast.success("Test cases validated. Moving to next step.");
      queryClient.invalidateQueries({ queryKey: ["problem", problemId] });
      navigate(`/problems/${problemId}/languages`)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Validation failed");
    },
  });

  if (isLoading) return (
    <Box display="flex" justifyContent="center" alignItems="center" height="250px">
      <CircularProgress size={40} />
    </Box>
  );
  if (error || !problem?.data) return <Alert severity="error">Problem not found</Alert>;

  return (
    <Box sx={{ p: 4 }}>
      <Box textAlign="center" mb={4}>
        <Typography variant="h4" fontWeight="bold" gutterBottom>
          Test Cases - {problem.data.title}
        </Typography>
        <Box mt={2}>
          <ProblemStepper currentStep={2} model="validate" problemId={problemId || "create"} />
        </Box>
      </Box>

      <Stack spacing={4}>
        <Paper elevation={3} sx={{ p: 3 }}>
          <TestCaseList problemId={Number(problemId)} />
        </Paper>

        <Box textAlign="center">
          <Button
            variant="contained"
            size="large"
            disabled={validateMutation.isPending}
            onClick={() => validateMutation.mutate()}
          >
            {validateMutation.isPending ? 'Validating...' : 'Validate Test Cases'}
          </Button>
        </Box>
      </Stack>
    </Box>
  );
}
