import { useParams } from "react-router-dom";
import { useState, useEffect } from "react";
import {
  Box,
  Container,
  Paper,
  Typography,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Tabs,
  Tab,
  Alert,
  Chip,
  CircularProgress,
} from "@mui/material";
import Editor from "@monaco-editor/react";
import CheckCircleOutlineIcon from "@mui/icons-material/CheckCircleOutline";
import { ProblemStepper } from "../components/ProblemStepper";
import { adminLanguagesApi } from "../../../api/adminApi";
import type { Language } from "../../../types";

type ProblemLanguage = {
  language_id: number;
  language_name: string;
  function_code: string;
  main_code: string;
  solution_code: string;
  is_validated: boolean;
};

type EditorTab = "function" | "main" | "solution";

export default function ProblemLanguage() {
  const { problemId } = useParams<{ problemId: string }>();

  const [availableLanguages, setAvailableLanguages] = useState<Language[]>([]);
  const [problemLanguages, setProblemLanguages] = useState<ProblemLanguage[]>([]);
  const [selectedLangId, setSelectedLangId] = useState<number | null>(null);
  const [currentTab, setCurrentTab] = useState<EditorTab>("function");

  const [functionCode, setFunctionCode] = useState("");
  const [mainCode, setMainCode] = useState("");
  const [solutionCode, setSolutionCode] = useState("");

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Fetch active languages
  useEffect(() => {
    const fetchLanguages = async () => {
      try {
        const response = await adminLanguagesApi.getAllActive();
        const languages = response.data.data.map((lang: any) => ({
          id: lang.id,
          name: lang.name,
          display_name: `${lang.name.charAt(0).toUpperCase() + lang.name.slice(1)} ${lang.version}`,
          version: lang.version,
          extension: lang.extension,
          default_template: lang.default_template,
        }));
        setAvailableLanguages(languages);
      } catch (err) {
        setError("Failed to load languages");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchLanguages();
  }, []);

  // Fetch problem languages
  useEffect(() => {
    const fetchProblemLanguages = async () => {
      if (!problemId) return;

      try {
        // TODO: Replace with actual API call
        // const response = await fetch(`/admin/api/problems/${problemId}/languages`);
        // const data = await response.json();
        // setProblemLanguages(data.languages);

        // Mock data for now
        setProblemLanguages([]);
      } catch (err) {
        console.error("Failed to load problem languages:", err);
        setProblemLanguages([]);
      }
    };

    fetchProblemLanguages();
  }, [problemId]);

  // Load selected language's code into editors
  useEffect(() => {
    if (selectedLangId) {
      const lang = problemLanguages.find((l) => l.language_id === selectedLangId);
      if (lang) {
        setFunctionCode(lang.function_code);
        setMainCode(lang.main_code);
        setSolutionCode(lang.solution_code);
      }
    }
  }, [selectedLangId, problemLanguages]);

  const getEditorLanguage = (langName: string) => {
    const map: Record<string, string> = {
      c: "c",
      "c++": "cpp",
      cpp: "cpp",
      python: "python",
      java: "java",
    };
    return map[langName.toLowerCase()] || "plaintext";
  };

  const handleSave = async () => {
    if (!selectedLangId) {
      setError("Please select a language first");
      return;
    }

    if (!functionCode.trim() || !mainCode.trim() || !solutionCode.trim()) {
      setError("All code fields are required");
      return;
    }

    try {
      const existingLang = problemLanguages.find((l) => l.language_id === selectedLangId);

      if (existingLang) {
        // Update existing
        // await fetch(`/admin/api/problems/${problemId}/languages/${selectedLangId}`, {
        //   method: "PUT",
        //   headers: { "Content-Type": "application/json" },
        //   body: JSON.stringify({
        //     function_code: functionCode,
        //     main_code: mainCode,
        //     solution_code: solutionCode,
        //   }),
        // });

        setProblemLanguages((prev) =>
          prev.map((l) =>
            l.language_id === selectedLangId
              ? {
                  ...l,
                  function_code: functionCode,
                  main_code: mainCode,
                  solution_code: solutionCode,
                  is_validated: false,
                }
              : l
          )
        );
        setSuccess("Language updated successfully");
      } else {
        // Create new
        // await fetch(`/admin/api/problems/${problemId}/languages`, {
        //   method: "POST",
        //   headers: { "Content-Type": "application/json" },
        //   body: JSON.stringify({
        //     language_id: selectedLangId,
        //     function_code: functionCode,
        //     main_code: mainCode,
        //     solution_code: solutionCode,
        //   }),
        // });

        const selectedAvailableLang = availableLanguages.find((l) => l.id === selectedLangId);
        if (selectedAvailableLang) {
          const newProblemLang: ProblemLanguage = {
            language_id: selectedLangId,
            language_name: selectedAvailableLang.name,
            function_code: functionCode,
            main_code: mainCode,
            solution_code: solutionCode,
            is_validated: false,
          };
          setProblemLanguages([...problemLanguages, newProblemLang]);
          setSuccess("Language added successfully");
        }
      }

      setError("");
    } catch (err) {
      setError("Failed to save");
      console.error(err);
    }
  };

  const handleDelete = async (langId: number) => {
    if (!window.confirm("Are you sure you want to remove this language?")) {
      return;
    }

    try {
      // await fetch(`/admin/api/problems/${problemId}/languages/${langId}`, {
      //   method: "DELETE",
      // });

      setProblemLanguages((prev) => prev.filter((l) => l.language_id !== langId));
      if (selectedLangId === langId) {
        setSelectedLangId(null);
        setFunctionCode("");
        setMainCode("");
        setSolutionCode("");
      }
      setSuccess("Language removed");
    } catch (err) {
      setError("Failed to remove language");
      console.error(err);
    }
  };

  const handleSelectLanguageToAdd = (langId: number) => {
    const lang = availableLanguages.find((l) => l.id === langId);
    if (lang) {
      setSelectedLangId(langId);
      setFunctionCode(lang.default_template || "// Write starter code here");
      setMainCode("// Write I/O handling code here");
      setSolutionCode("// Write your solution here");
      setCurrentTab("function");
      setSuccess(`${lang.name} selected. Fill in the code and click Save.`);
      setError("");
    }
  };

  const handleSelectExistingLanguage = (langId: number) => {
    setSelectedLangId(langId);
    setCurrentTab("function");
    setError("");
  };

  const getCurrentCode = () => {
    if (currentTab === "function") return functionCode;
    if (currentTab === "main") return mainCode;
    return solutionCode;
  };

  const handleEditorChange = (value: string | undefined) => {
    const code = value || "";
    if (currentTab === "function") setFunctionCode(code);
    else if (currentTab === "main") setMainCode(code);
    else setSolutionCode(code);
  };

  const selectedLang =
    problemLanguages.find((l) => l.language_id === selectedLangId) ||
    (selectedLangId ? availableLanguages.find((l) => l.id === selectedLangId) : null);

  const unusedLanguages = availableLanguages.filter(
    (al) => !problemLanguages.some((pl) => pl.language_id === al.id)
  );

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ mt: 4, display: "flex", justifyContent: "center" }}>
        <CircularProgress />
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <ProblemStepper currentStep={3} model="validate" />

      <Paper sx={{ p: 3, mb: 2 }}>
        <Typography variant="h5" gutterBottom>
          Configure Languages
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Add boilerplate code and solution for each language
        </Typography>
      </Paper>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError("")}>
          {error}
        </Alert>
      )}
      {success && (
        <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess("")}>
          {success}
        </Alert>
      )}

      <Paper sx={{ p: 3 }}>
        {/* Language selector */}
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: "flex", gap: 2, alignItems: "center", mb: 2 }}>
            <FormControl sx={{ minWidth: 250 }}>
              <InputLabel>Add a language</InputLabel>
              <Select
                value=""
                onChange={(e) => {
                    const langId = Number(e.target.value);
                    if (langId) handleSelectLanguageToAdd(langId);
                }}
                label="Add a language"
                disabled={unusedLanguages.length === 0}
              >
                {unusedLanguages.map((lang) => (
                  <MenuItem key={lang.id} value={lang.id}>
                    {lang.name + "(" + lang.version + ")"}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            {unusedLanguages.length === 0 && problemLanguages.length > 0 && (
              <Typography variant="body2" color="text.secondary">
                All languages added
              </Typography>
            )}
          </Box>

          {/* Added languages as chips */}
          {problemLanguages.length > 0 && (
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Added Languages:
              </Typography>
              <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
                {problemLanguages.map((lang) => (
                  <Chip
                    key={lang.language_id}
                    label={lang.language_name}
                    onClick={() => handleSelectExistingLanguage(lang.language_id)}
                    onDelete={() => handleDelete(lang.language_id)}
                    color={selectedLangId === lang.language_id ? "primary" : "default"}
                    icon={lang.is_validated ? <CheckCircleOutlineIcon /> : undefined}
                  />
                ))}
              </Box>
            </Box>
          )}
        </Box>

        {/* Code editor */}
        {selectedLangId ? (
          <>
            <Box sx={{ display: "flex", justifyContent: "space-between", alignItems: "center", mb: 2 }}>
              <Typography variant="h6">
                {(selectedLang as any)?.display_name || "Selected Language"}
              </Typography>
              <Button variant="contained" onClick={handleSave}>
                Save
              </Button>
            </Box>

            <Tabs value={currentTab} onChange={(_, val) => setCurrentTab(val)} sx={{ mb: 2 }}>
              <Tab label="Function Code (User sees)" value="function" />
              <Tab label="Main Code (I/O)" value="main" />
              <Tab label="Solution Code (Validator)" value="solution" />
            </Tabs>

            <Box sx={{ border: 1, borderColor: "divider", borderRadius: 1, overflow: "hidden" }}>
              <Editor
                height="500px"
                language={getEditorLanguage(
                  (selectedLang as any)?.language_name || (selectedLang as any)?.name || ""
                )}
                value={getCurrentCode()}
                onChange={handleEditorChange}
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
        ) : (
          <Box sx={{ textAlign: "center", py: 8, color: "text.secondary" }}>
            <Typography>Select a language from the dropdown above to start</Typography>
          </Box>
        )}
      </Paper>

      {/* Navigation */}
      <Box sx={{ mt: 3, display: "flex", justifyContent: "space-between" }}>
        <Button variant="outlined" onClick={() => window.history.back()}>
          Back
        </Button>
        <Button
          variant="contained"
          disabled={problemLanguages.length === 0}
          onClick={() => (window.location.href = `/admin/problems/${problemId}/validate`)}
        >
          Continue to Validate
        </Button>
      </Box>
    </Container>
  );
}
