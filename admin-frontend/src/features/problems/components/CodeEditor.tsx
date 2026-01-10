import { Box, Typography, Button, Tabs, Tab } from "@mui/material";
import Editor from "@monaco-editor/react";

type EditorTab = "function" | "main" | "solution";

type Props = {
  languageName: string;
  editorLanguage: string;
  currentTab: EditorTab;
  editorContent: Record<EditorTab, string>;
  onTabChange: (tab: EditorTab) => void;
  onCodeChange: (code: string) => void;
  onSave: () => void;
};

export function CodeEditor({
  languageName,
  editorLanguage,
  currentTab,
  editorContent,
  onTabChange,
  onCodeChange,
  onSave,
}: Props) {
  return (
    <>
      <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
        <Typography variant="h6">{languageName}</Typography>
        <Button variant="contained" onClick={onSave}>
          Save
        </Button>
      </Box>

      <Tabs
        value={currentTab}
        onChange={(_, val) => onTabChange(val as EditorTab)}
        sx={{ mb: 2 }}
      >
        <Tab label="Function Code (User sees)" value="function" />
        <Tab label="Main Code (I/O)" value="main" />
        <Tab label="Solution Code (Validator)" value="solution" />
      </Tabs>

      <Box sx={{ border: 1, borderColor: "divider", borderRadius: 1, overflow: "hidden" }}>
        <Editor
          height="500px"
          language={editorLanguage}
          value={editorContent[currentTab]}
          onChange={(value) => onCodeChange(value || "")}
          theme="vs-dark"
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            lineNumbers: "on",
            scrollBeyondLastLine: false,
          }}
        />
      </Box>
    </>
  );
}
