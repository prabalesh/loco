import { Alert } from "@mui/material";

type Props = {
  message: string;
  severity: "error" | "success";
  onClose: () => void;
};

export function AlertMessage({ message, severity, onClose }: Props) {
  if (!message) return null;

  return (
    <Alert severity={severity} sx={{ mb: 2 }} onClose={onClose}>
      {message}
    </Alert>
  );
}
