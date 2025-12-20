import { useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Spin, Alert, Button } from "antd";
import { adminProblemApi } from "../../../api/adminApi";
import ProblemDetails from "../components/ProblemDetails";
import { ProblemStepper } from "../components/ProblemStepper";
import toast from "react-hot-toast";
import TestCaseList from "../components/TestCasetList";

export default function ProblemTestCases() {
  const { problemId } = useParams<{ problemId: string }>();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["problem", problemId],
    queryFn: () => adminProblemApi.getById(Number(problemId)),
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
      // Optionally, redirect or update UI
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || "Validation failed");
    },
  });

  if (isLoading) return <Spin size="large" className="flex h-64 items-center justify-center" />;
  if (error || !problem?.data) return <Alert message="Problem not found" type="error" />;

  return (
    <div className="p-6 space-y-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-4">Test Cases - {problem.data.title}</h1>
        <div>
          <ProblemStepper currentStep={2} model="validate" />
        </div>
      </div>

      <div className="space-y-6">
        <div className="bg-white rounded-lg shadow-lg p-6">
          <ProblemDetails problem={problem.data} />
        </div>
        <div className="bg-white rounded-lg shadow-lg p-6">
          <TestCaseList problemId={Number(problemId)} />
        </div>
        <div className="text-center">
          <Button
            type="primary"
            loading={validateMutation.isPending}
            onClick={() => validateMutation.mutate()}
          >
            Validate Test Cases
          </Button>
        </div>
      </div>
    </div>
  );
}
