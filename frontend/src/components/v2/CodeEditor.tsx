import React, { useEffect, useState } from 'react';
import {
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Box,
  CircularProgress,
  Typography,
  Paper,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  Alert,
  Stack,
  Chip,
  IconButton,
  type SelectChangeEvent
} from '@mui/material';
import { Close as CloseIcon, Send as SendIcon, CheckCircle as CheckIcon, Error as ErrorIcon } from '@mui/icons-material';
import Editor from '@monaco-editor/react';
import { problemsApi } from '@/features/problems/api/problems';

interface CodeEditorProps {
  problemId: number;
  onCodeChange: (code: string, language: string) => void;
}

const LANGUAGES = [
  { value: 'python', label: 'Python' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'java', label: 'Java' },
  { value: 'cpp', label: 'C++' },
  { value: 'go', label: 'Go' },
];

export const CodeEditor: React.FC<CodeEditorProps> = ({ problemId, onCodeChange }) => {
  const [language, setLanguage] = useState('python');
  const [loading, setLoading] = useState(true);
  const [code, setCode] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [showResult, setShowResult] = useState(false);

  useEffect(() => {
    fetchStubCode();
  }, [problemId, language]);

  const handleSubmit = async () => {
    if (!code.trim()) return;
    setSubmitting(true);
    try {
      const response = await problemsApi.submit(problemId, {
        code: code,
        language_slug: language,
      });
      setResult(response.data.data);
      setShowResult(true);
    } catch (error) {
      console.error('Submission failed:', error);
    } finally {
      setSubmitting(false);
    }
  };

  const fetchStubCode = async () => {
    setLoading(true);
    try {
      const response = await problemsApi.getStub(problemId, language);

      const stubCode = response.data.data?.stub_code || '';
      setCode(stubCode);
      onCodeChange(stubCode, language);
    } catch (error) {
      console.error('Failed to load stub code:', error);
      setCode('');
    } finally {
      setLoading(false);
    }
  };

  const handleCodeChange = (newCode: string | undefined) => {
    const val = newCode || '';
    setCode(val);
    onCodeChange(val, language);
  };

  const handleLanguageChange = (event: SelectChangeEvent) => {
    setLanguage(event.target.value);
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
        <FormControl size="small" sx={{ width: 200 }}>
          <InputLabel>Language</InputLabel>
          <Select
            value={language}
            label="Language"
            onChange={handleLanguageChange}
            disabled={submitting}
          >
            {LANGUAGES.map((lang) => (
              <MenuItem key={lang.value} value={lang.value}>
                {lang.label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        {loading && <CircularProgress size={24} />}
        <Box sx={{ flexGrow: 1 }} />
        <Button
          variant="contained"
          color="primary"
          startIcon={submitting ? <CircularProgress size={20} color="inherit" /> : <SendIcon />}
          onClick={handleSubmit}
          disabled={submitting || !code.trim() || loading}
          sx={{
            px: 4,
            height: 40,
            borderRadius: 2,
            textTransform: 'none',
            fontWeight: 600
          }}
        >
          {submitting ? 'Running...' : 'Submit Code'}
        </Button>
      </Box>

      <Paper
        elevation={0}
        sx={{
          border: '1px solid',
          borderColor: 'divider',
          borderRadius: 2,
          overflow: 'hidden',
          bgcolor: '#1e1e1e'
        }}
      >
        <Editor
          height="500px"
          language={language === 'cpp' ? 'cpp' : language}
          theme="vs-dark"
          value={code}
          onChange={handleCodeChange}
          loading={
            <Box sx={{ height: 500, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white' }}>
              <Typography>Loading Editor...</Typography>
            </Box>
          }
          options={{
            selectOnLineNumbers: true,
            minimap: { enabled: false },
            fontSize: 14,
            padding: { top: 16 },
            readOnly: submitting
          }}
        />
      </Paper>

      <Dialog
        open={showResult}
        onClose={() => setShowResult(false)}
        maxWidth="md"
        fullWidth
        PaperProps={{
          sx: { borderRadius: 3, bgcolor: 'background.paper' }
        }}
      >
        <DialogTitle sx={{ m: 0, p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Typography variant="h6" fontWeight={700}>
            Submission Result
          </Typography>
          <IconButton onClick={() => setShowResult(false)} size="small">
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent dividers sx={{ p: 3 }}>
          {result && (
            <Stack spacing={3}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <Chip
                  label={result.status}
                  color={result.status === 'Accepted' ? 'success' : 'error'}
                  onDelete={() => { }}
                  deleteIcon={result.status === 'Accepted' ? <CheckIcon /> : <ErrorIcon />}
                  sx={{
                    fontWeight: 700,
                    px: 1,
                    '& .MuiChip-deleteIcon': { color: 'inherit' }
                  }}
                />
                <Typography variant="body1" color="text.secondary">
                  Passed {result.passed_tests} / {result.total_tests} test cases
                </Typography>
              </Box>

              {result.error_message && (
                <Alert severity="error" variant="outlined" sx={{ borderRadius: 2 }}>
                  <Typography variant="subtitle2" fontWeight={700} gutterBottom>
                    Error Detail:
                  </Typography>
                  <Box
                    component="pre"
                    sx={{
                      m: 0,
                      p: 1.5,
                      bgcolor: 'rgba(211, 47, 47, 0.05)',
                      borderRadius: 1,
                      overflow: 'auto',
                      fontSize: '0.875rem',
                      fontFamily: 'Monaco, Consolas, "Liberation Mono", monospace'
                    }}
                  >
                    {result.error_message}
                  </Box>
                </Alert>
              )}

              <Typography variant="subtitle1" fontWeight={700}>
                Test Cases
              </Typography>
              <Stack spacing={2}>
                {result.test_results?.map((test: any, idx: number) => (
                  <Paper
                    key={idx}
                    variant="outlined"
                    sx={{
                      p: 2,
                      borderRadius: 2,
                      borderColor: test.status === 'Passed' ? 'success.light' : 'error.light',
                      bgcolor: test.status === 'Passed' ? 'success.lighter' : 'error.lighter'
                    }}
                  >
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                      <Typography variant="subtitle2" fontWeight={700}>
                        Test Case {idx + 1} {test.is_sample && <Chip label="Sample" size="small" sx={{ ml: 1, height: 20 }} />}
                      </Typography>
                      <Typography
                        variant="caption"
                        fontWeight={700}
                        color={test.status === 'Passed' ? 'success.main' : 'error.main'}
                      >
                        {test.status.toUpperCase()}
                      </Typography>
                    </Box>
                    <Stack spacing={1} sx={{ fontSize: '0.875rem' }}>
                      <Box sx={{ display: 'flex', gap: 1 }}>
                        <Typography variant="caption" sx={{ minWidth: 60, fontWeight: 600 }}>Input:</Typography>
                        <Box component="code" sx={{ bgcolor: 'action.hover', px: 0.5, borderRadius: 0.5 }}>{test.input}</Box>
                      </Box>
                      <Box sx={{ display: 'flex', gap: 1 }}>
                        <Typography variant="caption" sx={{ minWidth: 60, fontWeight: 600 }}>Expected:</Typography>
                        <Box component="code" sx={{ bgcolor: 'action.hover', px: 0.5, borderRadius: 0.5 }}>{test.expected_output}</Box>
                      </Box>
                      {test.status !== 'Passed' && (
                        <Box sx={{ display: 'flex', gap: 1 }}>
                          <Typography variant="caption" sx={{ minWidth: 60, fontWeight: 600 }}>Actual:</Typography>
                          <Box component="code" sx={{ bgcolor: 'error.lighter', color: 'error.main', px: 0.5, borderRadius: 0.5 }}>{test.actual_output || 'N/A'}</Box>
                        </Box>
                      )}
                    </Stack>
                  </Paper>
                ))}
              </Stack>
            </Stack>
          )}
        </DialogContent>
      </Dialog>
    </Box>
  );
};

export default CodeEditor;
