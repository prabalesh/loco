import React from 'react';
import {
    TextField,
    Stack,
    Typography,
    Box,
    Button,
    IconButton,
    Paper,
    Checkbox,
    FormControlLabel,
    Alert,
    CircularProgress
} from '@mui/material';
import {
    Add as PlusIcon,
    Delete as DeleteIcon
} from '@mui/icons-material';
import BulkTestCaseDialog from '../BulkTestCaseDialog';

interface TestCasesStepProps {
    data: any;
    onChange: (newData: Partial<any>) => void;
    onSave?: () => void;
    saving?: boolean;
}

const getDefaultValueForType = (type: string) => {
    switch (type) {
        case 'integer':
        case 'double':
            return 0;
        case 'integer_array':
        case 'string_array':
            return [];
        case 'string':
            return "string";
        case 'boolean':
            return false;
        default:
            return null;
    }
};

export const TestCasesStep: React.FC<TestCasesStepProps> = ({ data, onChange, onSave, saving }) => {
    const [isBulkModalOpen, setIsBulkModalOpen] = React.useState(false);

    const handleAddTestCase = () => {
        const defaultInputs = data.parameters.map((p: any) => getDefaultValueForType(p.type));
        const defaultOutput = getDefaultValueForType(data.return_type);

        const newCases = [
            ...data.test_cases,
            {
                input: JSON.stringify(defaultInputs),
                expected_output: JSON.stringify(defaultOutput),
                is_sample: false
            }
        ];
        onChange({ test_cases: newCases });
    };

    const handleRemoveTestCase = (index: number) => {
        const newCases = [...data.test_cases];
        newCases.splice(index, 1);
        onChange({ test_cases: newCases });
    };

    const handleTestCaseChange = (index: number, field: string, value: any) => {
        const newCases = [...data.test_cases];
        newCases[index] = { ...newCases[index], [field]: value };
        onChange({ test_cases: newCases });
    };

    const handleBulkImport = async (importedCases: any[]) => {
        // importedCases are already formatted and stringified by the dialog
        const newCases = [...data.test_cases, ...importedCases];
        onChange({ test_cases: newCases });
    };

    const paramNames = data.parameters.map((p: any) => p.name || `param${data.parameters.indexOf(p) + 1}`);

    return (
        <Stack spacing={4}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Box>
                    <Typography variant="h6" color="primary">Test Cases</Typography>
                    <Typography variant="body2" color="text.secondary">
                        Provide inputs as a JSON array corresponding to function parameters.
                    </Typography>
                </Box>
                <Stack direction="row" spacing={2}>
                    {onSave && (
                        <Button
                            variant="contained"
                            color="primary"
                            onClick={onSave}
                            disabled={saving}
                            startIcon={saving ? <CircularProgress size={20} color="inherit" /> : null}
                        >
                            {saving ? 'Saving...' : 'Save Test Cases'}
                        </Button>
                    )}
                    <Button startIcon={<PlusIcon />} variant="outlined" onClick={handleAddTestCase}>
                        Add Test Case
                    </Button>
                    <Button variant="outlined" color="secondary" onClick={() => setIsBulkModalOpen(true)}>
                        Bulk Add
                    </Button>
                </Stack>
            </Box>

            <Alert severity="info">
                <Box>
                    <Typography variant="subtitle2">Signature: <code>{data.function_name}({paramNames.join(', ')})</code> returns <code>{data.return_type}</code></Typography>
                    <Typography variant="body2" sx={{ mt: 1 }}>
                        <strong>JSON Format Example:</strong>
                        <Box sx={{ mt: 0.5, fontFamily: 'monospace', bgcolor: 'rgba(0,0,0,0.05)', p: 1, borderRadius: 1 }}>
                            Input: <code>{JSON.stringify(data.parameters.map((p: any) => getDefaultValueForType(p.type)))}</code>
                            <br />
                            Output: <code>{JSON.stringify(getDefaultValueForType(data.return_type))}</code>
                        </Box>
                    </Typography>
                </Box>
            </Alert>

            <Stack spacing={3}>
                {data.test_cases.map((tc: any, index: number) => (
                    <Paper key={index} variant="outlined" sx={{ p: 3, position: 'relative' }}>
                        <Box sx={{ position: 'absolute', top: 12, right: 12 }}>
                            <IconButton
                                color="error"
                                size="small"
                                onClick={() => handleRemoveTestCase(index)}
                                disabled={data.test_cases.length === 1}
                            >
                                <DeleteIcon />
                            </IconButton>
                        </Box>

                        <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
                            <div className="md:col-span-12">
                                <Typography variant="subtitle2" gutterBottom>
                                    Test Case #{index + 1} {tc.is_sample && <span style={{ color: '#2e7d32' }}>(Sample)</span>}
                                </Typography>
                            </div>

                            <div className="md:col-span-7">
                                <TextField
                                    fullWidth
                                    multiline
                                    rows={3}
                                    label={`Inputs [${paramNames.join(', ')}]`}
                                    value={tc.input}
                                    onChange={(e) => handleTestCaseChange(index, 'input', e.target.value)}
                                    placeholder="e.g., [[2,7,11,15], 9]"
                                    helperText="Must be a valid JSON array of arguments"
                                    sx={{ fontFamily: 'monospace' }}
                                />
                            </div>

                            <div className="md:col-span-5">
                                <Stack spacing={2}>
                                    <TextField
                                        fullWidth
                                        multiline
                                        rows={3}
                                        label="Expected Output"
                                        value={tc.expected_output}
                                        onChange={(e) => handleTestCaseChange(index, 'expected_output', e.target.value)}
                                        placeholder="e.g., [0, 1]"
                                        helperText="Valid JSON representing return type"
                                        sx={{ fontFamily: 'monospace' }}
                                    />
                                    <FormControlLabel
                                        control={
                                            <Checkbox
                                                checked={tc.is_sample}
                                                onChange={(e) => handleTestCaseChange(index, 'is_sample', e.target.checked)}
                                            />
                                        }
                                        label="Mark as Sample Case"
                                    />
                                </Stack>
                            </div>
                        </div>
                    </Paper>
                ))}
            </Stack>

            <BulkTestCaseDialog
                open={isBulkModalOpen}
                onClose={() => setIsBulkModalOpen(false)}
                onImport={handleBulkImport}
                parameters={data.parameters}
                returnType={data.return_type}
                isImporting={false} // State is managed locally in wizard
            />
        </Stack>
    );
};
