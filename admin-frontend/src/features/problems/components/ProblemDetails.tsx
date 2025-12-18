import { Tag, Descriptions } from "antd";
import type { Problem } from "../../../types";

interface ProblemDetailsProps {
  problem: Problem;
}

export default function ProblemDetails({ problem }: ProblemDetailsProps) {
  const getDifficultyColor = (difficulty: Problem["difficulty"]) => {
    return difficulty === "easy" ? "green" : difficulty === "medium" ? "orange" : "red";
  };

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">{problem.title}</h2>
      
      <Descriptions column={1} className="mb-6">
        <Descriptions.Item label="Slug">{problem.slug}</Descriptions.Item>
        <Descriptions.Item label="Difficulty">
          <Tag color={getDifficultyColor(problem.difficulty)}>{problem.difficulty.toUpperCase()}</Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Status">
          <Tag color={problem.status === "published" ? "blue" : "default"}>
            {problem.status.toUpperCase()}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Time Limit">{problem.time_limit}ms</Descriptions.Item>
        <Descriptions.Item label="Memory Limit">{problem.memory_limit}MB</Descriptions.Item>
        <Descriptions.Item label="Acceptance">
          {problem.acceptance_rate.toFixed(1)}%
        </Descriptions.Item>
      </Descriptions>
    </div>
  );
}
