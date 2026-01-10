import { Box, Typography } from "@mui/material";

export function EmptyEditorState() {
  return (
    <Box sx={{ textAlign: "center", py: 8, color: "text.secondary" }}>
      <Typography>Select a language from the dropdown above to start</Typography>
    </Box>
  );
}
