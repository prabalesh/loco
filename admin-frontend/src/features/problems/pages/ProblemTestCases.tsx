import { useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";

import { Spin, Alert } from "antd";
import { adminProblemApi } from "../../../api/adminApi";
import TestCaseList from "../components/TestCasetList";
import ProblemDetails from "../components/ProblemDetails";
import { ProblemStepper } from "../components/ProblemStepper";

export default function ProblemTestCases() {
  const { problemId } = useParams<{ problemId: string }>();

  const { data, isLoading, error } = useQuery({
    queryKey: ["problem", problemId],
    queryFn: () => adminProblemApi.getById(Number(problemId)),
  });

  const problem = data?.data || null

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

      <div className="">
        <div className="bg-white rounded-lg shadow-lg p-6">
          <ProblemDetails problem={problem.data} />
        </div>
        <div className="bg-white rounded-lg shadow-lg p-6">
          <TestCaseList problemId={Number(problemId)} />
        </div>
      </div>
    </div>
  );
}
