import { Chip, Typography, Box, Grid } from "@mui/material";
import type { Problem } from "../../../types";

interface ProblemDetailsProps {
  problem: Problem;
}

export default function ProblemDetails({ problem }: ProblemDetailsProps) {
  const getDifficultyColor = (difficulty: Problem["difficulty"]) => {
    return difficulty === "easy" ? "success" : difficulty === "medium" ? "warning" : "error";
  };

  const getStatusColor = (status: string) => {
    return status === "published" ? "primary" : "default";
  }

  const DetailItem = ({ label, children }: { label: string, children: React.ReactNode }) => (
    <Grid container spacing={1} sx={{ py: 1, borderBottom: '1px solid', borderColor: 'divider' }}>
      <Grid size={{ xs: 4, sm: 3 }}>
        <Typography variant="body2" color="text.secondary" fontWeight="medium">
          {label}
        </Typography>
      </Grid>
      <Grid size={{ xs: 8, sm: 9 }}>
        <Box>{children}</Box>
      </Grid>
    </Grid>
  );

  return (
    <Box>
      <Typography variant="h5" fontWeight="bold" gutterBottom>
        {problem.title}
      </Typography>

      <Box sx={{ mb: 3 }}>
        <DetailItem label="Slug">
          <Typography variant="body2">{problem.slug}</Typography>
        </DetailItem>
        <DetailItem label="Difficulty">
          <Chip
            label={problem.difficulty.toUpperCase()}
            color={getDifficultyColor(problem.difficulty)}
            size="small"
            variant="outlined"
          />
        </DetailItem>
        <DetailItem label="Status">
          <Chip
            label={problem.status.toUpperCase()}
            color={getStatusColor(problem.status) as any}
            size="small"
            variant="outlined"
          />
        </DetailItem>
        <DetailItem label="Time Limit">
          <Typography variant="body2">{problem.time_limit}ms</Typography>
        </DetailItem>
        <DetailItem label="Memory Limit">
          <Typography variant="body2">{problem.memory_limit}MB</Typography>
        </DetailItem>
        <DetailItem label="Acceptance">
          <Typography variant="body2">{problem.acceptance_rate.toFixed(1)}%</Typography>
        </DetailItem>
      </Box>
    </Box>
  );
}
