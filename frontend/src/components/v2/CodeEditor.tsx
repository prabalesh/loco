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
  type SelectChangeEvent
} from '@mui/material';
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

  useEffect(() => {
    fetchStubCode();
  }, [problemId, language]);

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
          >
            {LANGUAGES.map((lang) => (
              <MenuItem key={lang.value} value={lang.value}>
                {lang.label}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        {loading && <CircularProgress size={24} />}
      </Box>

      <Paper
        elevation={0}
        sx={{
          border: '1px solid',
          borderColor: 'divider',
          borderRadius: 1,
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
            padding: { top: 16 }
          }}
        />
      </Paper>
    </Box>
  );
};

export default CodeEditor;
