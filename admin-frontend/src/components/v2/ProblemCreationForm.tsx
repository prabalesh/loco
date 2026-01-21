import React, { useState, useEffect } from 'react';
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
  Checkbox,
  FormControlLabel,
  Tooltip
} from '@mui/material';
import { Add as PlusIcon, Delete as DeleteIcon, Visibility as PreviewIcon, Info as InfoIcon } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { adminCodeGenApi, adminProblemApi } from '../../lib/api/admin';
import toast from 'react-hot-toast';

interface Parameter {
  name: string;
  type: string;
  is_custom: boolean;
}

const PRIMITIVE_TYPES = [
  { value: 'int', label: 'int', is_custom: false },
  { value: 'int[]', label: 'int[]', is_custom: false },
  { value: 'string', label: 'string', is_custom: false },
  { value: 'string[]', label: 'string[]', is_custom: false },
  { value: 'bool', label: 'bool', is_custom: false },
  { value: 'double', label: 'double', is_custom: false },
];

const VALIDATION_TYPES = [
  { value: 'EXACT', label: 'Exact Match' },
  { value: 'UNORDERED', label: 'Unordered List' },
  { value: 'SUBSET', label: 'Subset Match' },
  { value: 'ANY_MATCH', label: 'Any Match (Custom)' },
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
  const [customTypes, setCustomTypes] = useState<{ value: string; label: string; is_custom: boolean }[]>([]);
  const [allTypes, setAllTypes] = useState<any[]>([...PRIMITIVE_TYPES]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [difficulty, setDifficulty] = useState('medium');
  const [functionName, setFunctionName] = useState('');
  const [returnType, setReturnType] = useState('int');
  const [parameters, setParameters] = useState<Parameter[]>([{ name: '', type: 'int', is_custom: false }]);
  const [generatedStub, setGeneratedStub] = useState<string>('');
  const [selectedLanguage, setSelectedLanguage] = useState<string>('python');

  const [validationType, setValidationType] = useState('EXACT');
  const [expectedTimeComplexity, setExpectedTimeComplexity] = useState('O(n)');
  const [expectedSpaceComplexity, setExpectedSpaceComplexity] = useState('O(1)');
  const [testCases, setTestCases] = useState<any[]>([{ input: '', expected_output: '', is_sample: true }]);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [creationStatus, setCreationStatus] = useState<'idle' | 'creating' | 'generating' | 'done'>('idle');
  const [boilerplateStats, setBoilerplateStats] = useState<{ total_languages: number; languages: string[] } | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    fetchCustomTypes();
  }, []);

  const fetchCustomTypes = async () => {
    try {
      const response = await adminProblemApi.getCustomTypes();
      const types = response.data.data.map((ct: any) => ({
        value: ct.name,
        label: ct.name,
        is_custom: true,
      }));
      setCustomTypes(types);
      setAllTypes([...PRIMITIVE_TYPES, ...types]);
    } catch (err) {
      console.error('Failed to fetch custom types', err);
    }
  };

  useEffect(() => {
    if (customTypes.length > 0) {
      console.log('Loaded custom types:', customTypes);
    }
  }, [customTypes]);

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
    if (field === 'type') {
      // Check if custom
      const typeObj = allTypes.find(t => t.value === value);
      newParams[index] = { ...newParams[index], type: value, is_custom: typeObj?.is_custom || false };
    } else {
      newParams[index] = { ...newParams[index], [field]: value };
    }
    setParameters(newParams);
  };

  const addTestCase = () => {
    setTestCases([...testCases, { input: '', expected_output: '', is_sample: false }]);
  };

  const removeTestCase = (index: number) => {
    const newCases = [...testCases];
    newCases.splice(index, 1);
    setTestCases(newCases);
  };

  const updateTestCase = (index: number, field: string, value: any) => {
    const newCases = [...testCases];
    newCases[index] = { ...newCases[index], [field]: value };
    setTestCases(newCases);
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

  const pollBoilerplates = async (problemId: number) => {
    setCreationStatus('generating');
    let attempts = 0;
    const maxAttempts = 15;

    const interval = setInterval(async () => {
      try {
        const response = await adminCodeGenApi.getBoilerplateStats(problemId);
        const stats = response.data.data;
        if (stats) {
          setBoilerplateStats(stats);
          if (stats.total_languages >= 5 || attempts >= maxAttempts) {
            clearInterval(interval);
            setCreationStatus('done');
            toast.success(`Problem created and boilerplates generated!`);
            setTimeout(() => {
              navigate('/problems');
            }, 3000);
          }
        }
        attempts++;
      } catch (err) {
        console.error('Failed to poll boilerplates', err);
        attempts++;
        if (attempts >= maxAttempts) clearInterval(interval);
      }
    }, 2000);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setCreationStatus('creating');
    setError(null);

    try {
      // Parse test cases JSON
      const formattedTestCases = testCases.map((tc, idx) => {
        try {
          return {
            ...tc,
            input: JSON.parse(tc.input || '[]'),
            expected_output: JSON.parse(tc.expected_output || 'null')
          };
        } catch (err) {
          throw new Error(`Invalid JSON in test case #${idx + 1}`);
        }
      });

      const problemData = {
        title,
        description,
        difficulty,
        function_name: functionName,
        return_type: returnType,
        parameters,
        test_cases: formattedTestCases,
        validation_type: validationType,
        expected_time_complexity: expectedTimeComplexity,
        expected_space_complexity: expectedSpaceComplexity,
        status: 'draft',
        visibility: 'public',
        is_active: true,
        tag_ids: [],
        category_ids: []
      };

      // Using bracket notation to bypass potential TS issues with dynamic method addition
      const api = adminProblemApi as any;
      const response = await api.v2Create(problemData);
      const newProblem = response.data.data;

      if (newProblem && newProblem.id) {
        pollBoilerplates(newProblem.id);
      } else {
        setCreationStatus('done');
        toast.success('Problem created successfully!');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to create problem');
      setCreationStatus('idle');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: 1000, margin: '2rem auto', p: 2 }}>
      <form onSubmit={handleSubmit}>
        <Card elevation={3} sx={{ mb: 4 }}>
          <CardHeader
            title="Create New Problem (Automated V2)"
            subheader="Define the signature and test cases, we'll handle the boilerplate."
          />
          <Divider />
          <CardContent>
            <Stack spacing={4}>
              {/* Basic Info */}
              <Box>
                <Typography variant="h6" gutterBottom color="primary">Basic Information</Typography>
                <div className="grid grid-cols-1 sm:grid-cols-12 gap-6">
                  <div className="sm:col-span-8">
                    <TextField
                      fullWidth
                      label="Problem Title"
                      value={title}
                      onChange={(e) => setTitle(e.target.value)}
                      required
                      placeholder="e.g., Two Sum"
                    />
                  </div>
                  <div className="sm:col-span-4">
                    <TextField
                      select
                      fullWidth
                      label="Difficulty"
                      value={difficulty}
                      onChange={(e) => setDifficulty(e.target.value)}
                    >
                      {DIFFICULTIES.map((opt) => (
                        <MenuItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </MenuItem>
                      ))}
                    </TextField>
                  </div>
                  <div className="sm:col-span-12">
                    <TextField
                      fullWidth
                      label="Description"
                      value={description}
                      onChange={(e) => setDescription(e.target.value)}
                      required
                      multiline
                      rows={6}
                      placeholder="Problem description in detail..."
                    />
                  </div>
                </div>
              </Box>

              <Divider />

              {/* Specification */}
              <Box>
                <Typography variant="h6" gutterBottom color="primary">Technical Specification</Typography>
                <Paper variant="outlined" sx={{ p: 3, bgcolor: '#fcfcfc' }}>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 mb-6">
                    <div>
                      <TextField
                        fullWidth
                        label="Function Name"
                        value={functionName}
                        onChange={(e) => setFunctionName(e.target.value)}
                        required
                        placeholder="e.g., twoSum"
                        helperText="Used for auto-generating function signature"
                      />
                    </div>
                    <div>
                      <TextField
                        select
                        fullWidth
                        label="Return Type"
                        value={returnType}
                        onChange={(e) => setReturnType(e.target.value)}
                      >
                        {allTypes.map((opt) => (
                          <MenuItem key={opt.value} value={opt.value}>
                            {opt.label}
                          </MenuItem>
                        ))}
                      </TextField>
                    </div>
                  </div>

                  <Typography variant="subtitle1" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    Parameters
                    <Tooltip title="Define input parameters for your function.">
                      <InfoIcon fontSize="small" color="action" />
                    </Tooltip>
                  </Typography>
                  <Stack spacing={2} sx={{ mb: 3 }}>
                    {parameters.map((param, index) => (
                      <Box key={index} sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                        <TextField
                          size="small"
                          label="Name"
                          value={param.name}
                          onChange={(e) => updateParameter(index, 'name', e.target.value)}
                          required
                          placeholder="e.g., nums"
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
                          {allTypes.map((opt) => (
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
                      variant="text"
                      onClick={addParameter}
                      sx={{ width: 'fit-content' }}
                    >
                      Add Parameter
                    </Button>
                  </Stack>

                  <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
                    <div>
                      <TextField
                        select
                        fullWidth
                        label="Validation Type"
                        value={validationType}
                        onChange={(e) => setValidationType(e.target.value)}
                      >
                        {VALIDATION_TYPES.map((opt) => (
                          <MenuItem key={opt.value} value={opt.value}>
                            {opt.label}
                          </MenuItem>
                        ))}
                      </TextField>
                    </div>
                    <div>
                      <TextField
                        fullWidth
                        label="Expected Time Complexity"
                        value={expectedTimeComplexity}
                        onChange={(e) => setExpectedTimeComplexity(e.target.value)}
                        placeholder="e.g., O(n log n)"
                      />
                    </div>
                    <div>
                      <TextField
                        fullWidth
                        label="Expected Space Complexity"
                        value={expectedSpaceComplexity}
                        onChange={(e) => setExpectedSpaceComplexity(e.target.value)}
                        placeholder="e.g., O(n)"
                      />
                    </div>
                  </div>

                  <Box sx={{ mt: 3, p: 2, borderRadius: 1, border: '1px dashed #ccc', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                      <TextField
                        select
                        size="small"
                        label="Preview Language"
                        value={selectedLanguage}
                        onChange={(e) => setSelectedLanguage(e.target.value)}
                        sx={{ width: 150 }}
                      >
                        {LANGUAGES.map((opt) => (
                          <MenuItem key={opt.value} value={opt.value}>
                            {opt.label}
                          </MenuItem>
                        ))}
                      </TextField>
                      <Button
                        variant="outlined"
                        color="secondary"
                        startIcon={<PreviewIcon />}
                        onClick={handlePreviewStub}
                        disabled={loading || !functionName || !parameters[0].name}
                      >
                        Preview Stub
                      </Button>
                    </Box>
                    <Typography variant="caption" color="text.secondary">
                      Stub code is auto-generated based on signature.
                    </Typography>
                  </Box>

                  {generatedStub && (
                    <Box sx={{ mt: 2 }}>
                      <Paper
                        sx={{
                          p: 2,
                          bgcolor: '#1e1e1e',
                          color: '#d4d4d4',
                          fontFamily: 'monospace',
                          overflowX: 'auto',
                          fontSize: '0.875rem'
                        }}
                      >
                        <pre style={{ margin: 0 }}>{generatedStub}</pre>
                      </Paper>
                    </Box>
                  )}
                </Paper>
              </Box>

              <Divider />

              {/* Test Cases */}
              <Box>
                <Typography variant="h6" gutterBottom color="primary" sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  Test Cases
                  <Button startIcon={<PlusIcon />} onClick={addTestCase}>Add Test Case</Button>
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  Provide inputs and expected outputs in JSON format.
                </Typography>

                <Stack spacing={3}>
                  {testCases.map((tc, index) => (
                    <Paper key={index} variant="outlined" sx={{ p: 2, position: 'relative' }}>
                      <Box sx={{ position: 'absolute', top: 8, right: 8 }}>
                        <IconButton
                          color="error"
                          size="small"
                          onClick={() => removeTestCase(index)}
                          disabled={testCases.length === 1}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </Box>
                      <div className="grid grid-cols-1 sm:grid-cols-12 gap-4">
                        <div className="sm:col-span-5">
                          <TextField
                            fullWidth
                            multiline
                            rows={2}
                            label={`Input (JSON)`}
                            value={tc.input}
                            onChange={(e) => updateTestCase(index, 'input', e.target.value)}
                            placeholder='e.g., [2, 7, 11, 15], 9'
                            size="small"
                            helperText="Values for each parameter in order"
                          />
                        </div>
                        <div className="sm:col-span-4">
                          <TextField
                            fullWidth
                            multiline
                            rows={2}
                            label="Expected Output (JSON)"
                            value={tc.expected_output}
                            onChange={(e) => updateTestCase(index, 'expected_output', e.target.value)}
                            placeholder='e.g., [0, 1]'
                            size="small"
                          />
                        </div>
                        <div className="sm:col-span-3 flex items-center">
                          <FormControlLabel
                            control={
                              <Checkbox
                                checked={tc.is_sample}
                                onChange={(e) => updateTestCase(index, 'is_sample', e.target.checked)}
                              />
                            }
                            label="Is Sample"
                          />
                        </div>
                      </div>
                    </Paper>
                  ))}
                </Stack>
              </Box>

              {error && <Alert severity="error">{error}</Alert>}

              {creationStatus === 'generating' && (
                <Alert severity="info" variant="filled">
                  Creating problem and generating boilerplates for all active languages...
                  {boilerplateStats && (
                    <Box sx={{ mt: 1 }}>
                      <strong>Ready:</strong> {boilerplateStats.languages.join(', ')} ({boilerplateStats.total_languages}/5)
                    </Box>
                  )}
                </Alert>
              )}

              <Box sx={{ pt: 2 }}>
                <Button
                  fullWidth
                  variant="contained"
                  size="large"
                  type="submit"
                  disabled={loading || creationStatus === 'creating' || creationStatus === 'generating'}
                  sx={{ height: 56, fontSize: '1.1rem' }}
                >
                  {creationStatus === 'creating' ? 'Creating Problem...' :
                    creationStatus === 'generating' ? 'Generating Boilerplates...' :
                      'Create Problem & Generate Code Stubs'}
                </Button>
              </Box>
            </Stack>
          </CardContent>
        </Card>
      </form>
    </Box>
  );
};

export default ProblemCreationForm;
