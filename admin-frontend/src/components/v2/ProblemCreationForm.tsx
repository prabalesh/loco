import React, { useState } from 'react';
import {
  Card,
  CardHeader,
  CardContent,
  TextField,
  MenuItem,
  Button,
  Box,
  IconButton,
  Divider,
  Stack,
  Typography,
  Alert,
  Paper,
  Grid
} from '@mui/material';
import { Add as PlusIcon, Delete as DeleteIcon, Visibility as PreviewIcon } from '@mui/icons-material';
import { adminCodeGenApi } from '../../lib/api/admin';

interface Parameter {
  name: string;
  type: string;
  is_custom: boolean;
}

const PRIMITIVE_TYPES = [
  { value: 'int', label: 'int' },
  { value: 'int[]', label: 'int[]' },
  { value: 'string', label: 'string' },
  { value: 'string[]', label: 'string[]' },
  { value: 'bool', label: 'bool' },
  { value: 'double', label: 'double' },
];

const DIFFICULTIES = [
  { value: 'easy', label: 'Easy' },
  { value: 'medium', label: 'Medium' },
  { value: 'hard', label: 'Hard' },
];

const LANGUAGES = [
  { value: 'python', label: 'Python' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'java', label: 'Java' },
  { value: 'cpp', label: 'C++' },
  { value: 'go', label: 'Go' },
];

export const ProblemCreationForm: React.FC = () => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [difficulty, setDifficulty] = useState('medium');
  const [functionName, setFunctionName] = useState('');
  const [returnType, setReturnType] = useState('int');
  const [parameters, setParameters] = useState<Parameter[]>([{ name: '', type: 'int', is_custom: false }]);
  const [generatedStub, setGeneratedStub] = useState<string>('');
  const [selectedLanguage, setSelectedLanguage] = useState<string>('python');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const addParameter = () => {
    setParameters([...parameters, { name: '', type: 'int', is_custom: false }]);
  };

  const removeParameter = (index: number) => {
    const newParams = [...parameters];
    newParams.splice(index, 1);
    setParameters(newParams);
  };

  const updateParameter = (index: number, field: keyof Parameter, value: any) => {
    const newParams = [...parameters];
    newParams[index] = { ...newParams[index], [field]: value };
    setParameters(newParams);
  };

  const handlePreviewStub = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await adminCodeGenApi.generateStub({
        function_name: functionName,
        return_type: returnType,
        parameters: parameters,
        language_slug: selectedLanguage,
      });

      setGeneratedStub(response.data.data.stub_code);
    } catch (err: any) {
      setError(err.response?.data?.message || err.message || 'Failed to generate stub');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    // Implementation for creating the actual problem in the DB
    console.log('Submit', { title, description, difficulty, functionName, returnType, parameters });
  };

  return (
    <Box sx={{ maxWidth: 1000, margin: '2rem auto', p: 2 }}>
      <form onSubmit={handleSubmit}>
        <Card elevation={3} sx={{ mb: 4 }}>
          <CardHeader title="Create New Problem (V2)" />
          <Divider />
          <CardContent>
            <Stack spacing={3}>
              <TextField
                fullWidth
                label="Problem Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                placeholder="e.g., Two Sum"
              />

              <TextField
                fullWidth
                label="Description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                required
                multiline
                rows={6}
                placeholder="Problem description..."
              />

              <TextField
                select
                label="Difficulty"
                value={difficulty}
                onChange={(e) => setDifficulty(e.target.value)}
                sx={{ width: 200 }}
              >
                {DIFFICULTIES.map((opt) => (
                  <MenuItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </MenuItem>
                ))}
              </TextField>

              <Paper variant="outlined" sx={{ p: 3, bgcolor: '#fafafa' }}>
                <Typography variant="h6" gutterBottom>
                  Function Signature
                </Typography>
                <Grid container spacing={2} sx={{ mb: 3 }}>
                  <Grid size={{ xs: 12, sm: 6 }}>
                    <TextField
                      fullWidth
                      label="Function Name"
                      value={functionName}
                      onChange={(e) => setFunctionName(e.target.value)}
                      required
                      placeholder="e.g., twoSum"
                    />
                  </Grid>
                  <Grid size={{ xs: 12, sm: 6 }}>
                    <TextField
                      select
                      fullWidth
                      label="Return Type"
                      value={returnType}
                      onChange={(e) => setReturnType(e.target.value)}
                    >
                      {PRIMITIVE_TYPES.map((opt) => (
                        <MenuItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </MenuItem>
                      ))}
                    </TextField>
                  </Grid>
                </Grid>

                <Typography variant="subtitle1" gutterBottom>
                  Parameters
                </Typography>
                <Stack spacing={2}>
                  {parameters.map((param, index) => (
                    <Box key={index} sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                      <TextField
                        size="small"
                        label="Name"
                        value={param.name}
                        onChange={(e) => updateParameter(index, 'name', e.target.value)}
                        required
                        sx={{ flex: 1 }}
                      />
                      <TextField
                        select
                        size="small"
                        label="Type"
                        value={param.type}
                        onChange={(e) => updateParameter(index, 'type', e.target.value)}
                        sx={{ width: 150 }}
                      >
                        {PRIMITIVE_TYPES.map((opt) => (
                          <MenuItem key={opt.value} value={opt.value}>
                            {opt.label}
                          </MenuItem>
                        ))}
                      </TextField>
                      <IconButton
                        color="error"
                        onClick={() => removeParameter(index)}
                        disabled={parameters.length === 1}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Box>
                  ))}
                  <Button
                    startIcon={<PlusIcon />}
                    variant="outlined"
                    onClick={addParameter}
                    sx={{ width: 'fit-content' }}
                  >
                    Add Parameter
                  </Button>
                </Stack>

                <Box sx={{ mt: 4, display: 'flex', gap: 2, alignItems: 'center' }}>
                  <TextField
                    select
                    size="small"
                    label="Preview Language"
                    value={selectedLanguage}
                    onChange={(e) => setSelectedLanguage(e.target.value)}
                    sx={{ width: 200 }}
                  >
                    {LANGUAGES.map((opt) => (
                      <MenuItem key={opt.value} value={opt.value}>
                        {opt.label}
                      </MenuItem>
                    ))}
                  </TextField>
                  <Button
                    variant="contained"
                    color="secondary"
                    startIcon={<PreviewIcon />}
                    onClick={handlePreviewStub}
                    disabled={loading || !functionName}
                  >
                    Preview Stub Code
                  </Button>
                </Box>

                {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}

                {generatedStub && (
                  <Box sx={{ mt: 2 }}>
                    <Paper
                      sx={{
                        p: 2,
                        bgcolor: '#1e1e1e',
                        color: '#d4d4d4',
                        fontFamily: 'monospace',
                        overflowX: 'auto'
                      }}
                    >
                      <pre style={{ margin: 0 }}>{generatedStub}</pre>
                    </Paper>
                  </Box>
                )}
              </Paper>

              <Button
                fullWidth
                variant="contained"
                size="large"
                type="submit"
                sx={{ mt: 2 }}
              >
                Create Problem
              </Button>
            </Stack>
          </CardContent>
        </Card>
      </form>
    </Box>
  );
};

export default ProblemCreationForm;
