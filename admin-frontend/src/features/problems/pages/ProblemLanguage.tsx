import { useParams, useNavigate } from "react-router-dom";
import { useState, useCallback, useEffect } from "react";
import {
  Container,
  Paper,
  Typography,
  Button,
  Box,
  CircularProgress,
} from "@mui/material";
import { ProblemStepper } from "../components/ProblemStepper";
import { useAvailableLanguages } from "../hooks/useAvailableLanguages";
import { useProblemLanguages } from "../hooks/useProblemLanguages";
import { useEditorState } from "../hooks/useEditorState";
import { useLanguageManager } from "../hooks/useLanguageManager";
import { LanguageSelector } from "../components/LanguageSelector";
import { LanguageChips } from "../components/LanguageChips";
import { CodeEditor } from "../components/CodeEditor";
import { EmptyEditorState } from "../components/EmptyEditorState";
import { AlertMessage } from "../components/AlertMessage";

export default function ProblemLanguage() {
  const { problemId } = useParams<{ problemId: string }>();
  const navigate = useNavigate();

  // Fetch available languages
  const { availableLanguages, loading: langLoading } = useAvailableLanguages();

  // Fetch and manage problem languages
  const {
    problemLanguages,
    loading: probLoading,
    saving,
    error: hookError,
    saveLanguage,
    deleteLanguage,
  } = useProblemLanguages(problemId);

  // Local UI state
  const [selectedLangId, setSelectedLangId] = useState<number | null>(null);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Editor state management
  const {
    currentTab,
    setCurrentTab,
    editorContent,
    updateEditorContent,
    resetEditorContent,
  } = useEditorState(selectedLangId, problemLanguages);

  // Language utilities
  const { unusedLanguages, getLanguageById, getEditorLanguage } =
    useLanguageManager(availableLanguages, problemLanguages);

  // Sync hook error with local error state
  useEffect(() => {
    if (hookError) {
      setError(hookError);
    }
  }, [hookError]);

  // Handle selecting a new language to add
  const handleSelectLanguageToAdd = useCallback(
    (langId: number) => {
      const lang = availableLanguages.find((l) => l.id === langId);
      if (lang) {
        setSelectedLangId(langId);
        resetEditorContent({
          function: lang.default_template || "// Write starter code here",
          main: "// Write I/O handling code here",
          solution: "// Write your solution here",
        });
        setCurrentTab("function");
        setSuccess(`${lang.name} selected. Fill in the code and click Save.`);
        setError("");
      }
    },
    [availableLanguages, resetEditorContent, setCurrentTab]
  );

  // Handle selecting an existing language
  const handleSelectExistingLanguage = useCallback(
    (langId: number) => {
      setSelectedLangId(langId);
      setCurrentTab("function");
      setError("");
      setSuccess("");
    },
    [setCurrentTab]
  );

  // Validate input fields
  const validateFields = useCallback(() => {
    const { function: functionCode, main: mainCode, solution: solutionCode } = editorContent;

    if (!functionCode.trim()) {
      setError("Function code is required");
      return false;
    }

    if (!mainCode.trim()) {
      setError("Main code (I/O handling) is required");
      return false;
    }

    if (!solutionCode.trim()) {
      setError("Solution code is required");
      return false;
    }

    return true;
  }, [editorContent]);

  // Handle save
  const handleSave = useCallback(async () => {
    if (!selectedLangId) {
      setError("Please select a language first");
      return;
    }

    if (!validateFields()) {
      return;
    }

    const { function: functionCode, main: mainCode, solution: solutionCode } = editorContent;

    const result = await saveLanguage(selectedLangId, {
      function_code: functionCode,
      main_code: mainCode,
      solution_code: solutionCode,
    });

    if (result.success) {
      setSuccess(result.message);
      setError("");
    } else {
      setError(result.message);
      setSuccess("");
    }
  }, [selectedLangId, editorContent, saveLanguage, validateFields]);

  // Handle delete
  const handleDelete = useCallback(
    async (langId: number) => {
      if (!window.confirm("Are you sure you want to remove this language?")) {
        return;
      }

      const result = await deleteLanguage(langId);

      if (result.success) {
        if (selectedLangId === langId) {
          setSelectedLangId(null);
          resetEditorContent({ function: "", main: "", solution: "" });
        }
        setSuccess(result.message);
        setError("");
      } else {
        setError(result.message);
        setSuccess("");
      }
    },
    [selectedLangId, deleteLanguage, resetEditorContent]
  );

  // Handle navigation
  const handleBack = useCallback(() => {
    navigate(`/admin/problems/${problemId}/test-cases`);
  }, [navigate, problemId]);

  const handleContinue = useCallback(() => {
    navigate(`/admin/problems/${problemId}/validate`);
  }, [navigate, problemId]);

  // Loading state
  if (langLoading || probLoading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, display: "flex", justifyContent: "center" }}>
        <CircularProgress />
      </Container>
    );
  }

  // Get current language details
  const selectedLang = selectedLangId ? getLanguageById(selectedLangId) : null;
  const languageName =
    (selectedLang as any)?.display_name ||
    (selectedLang as any)?.name ||
    "Selected Language";
  const editorLanguage = selectedLang
    ? getEditorLanguage(
        (selectedLang as any)?.language_name || (selectedLang as any)?.name || ""
      )
    : "plaintext";

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      {/* Stepper */}
      <ProblemStepper currentStep={3} model="validate" />

      {/* Header */}
      <Paper sx={{ p: 3, mb: 2 }}>
        <Typography variant="h5" gutterBottom>
          Configure Languages
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Add boilerplate code and solution for each language
        </Typography>
      </Paper>

      {/* Alerts */}
      <AlertMessage message={error} severity="error" onClose={() => setError("")} />
      <AlertMessage message={success} severity="success" onClose={() => setSuccess("")} />

      {/* Main Content */}
      <Paper sx={{ p: 3 }}>
        {/* Language Selection Section */}
        <Box sx={{ mb: 3 }}>
          <LanguageSelector
            unusedLanguages={unusedLanguages}
            onSelectLanguage={handleSelectLanguageToAdd}
            hasAddedLanguages={problemLanguages.length > 0}
          />

          <LanguageChips
            languages={problemLanguages}
            selectedLangId={selectedLangId}
            onSelectLanguage={handleSelectExistingLanguage}
            onDeleteLanguage={handleDelete}
          />
        </Box>

        {/* Code Editor Section */}
        {selectedLangId ? (
          <CodeEditor
            languageName={languageName}
            editorLanguage={editorLanguage}
            currentTab={currentTab}
            editorContent={editorContent}
            onTabChange={setCurrentTab}
            onCodeChange={(code) => updateEditorContent(currentTab, code)}
            onSave={handleSave}
            saving={saving}
          />
        ) : (
          <EmptyEditorState />
        )}
      </Paper>

      {/* Navigation */}
      <Box sx={{ mt: 3, display: "flex", justifyContent: "space-between" }}>
        <Button variant="outlined" onClick={handleBack}>
          Back
        </Button>
        <Button
          variant="contained"
          disabled={problemLanguages.length === 0}
          onClick={handleContinue}
        >
          Continue to Validate
        </Button>
      </Box>
    </Container>
  );
}
