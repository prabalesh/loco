import React, { useState } from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Typography,
    Box,
    Alert,
    CircularProgress,
    Stack,
} from '@mui/material';
import type { Parameter } from '../../../types';

interface BulkTestCaseDialogProps {
    open: boolean;
    onClose: () => void;
    onImport: (testCases: any[]) => Promise<void>;
    parameters?: Parameter[];
    returnType?: string;
    isImporting: boolean;
}

const BulkTestCaseDialog: React.FC<BulkTestCaseDialogProps> = ({
    open,
    onClose,
    onImport,
    parameters,
    returnType,
    isImporting,
}) => {
    const [jsonInput, setJsonInput] = useState('');
    const [error, setError] = useState<string | null>(null);

    const getDefaultValueForType = (type: string) => {
        switch (type.toLowerCase()) {
            case 'integer':
            case 'double':
                return 0;
            case 'integer_array':
                return [0, 1, 2];
            case 'string':
                return "";
            case 'string_array':
                return ["a", "b"];
            case 'boolean':
                return false;
            default:
                return null;
        }
    };

    const generateTemplate = () => {
        let template: any;
        const inputObj: Record<string, any> = {};
        if (parameters && parameters.length > 0) {
            parameters.forEach((p) => {
                inputObj[p.name] = getDefaultValueForType(p.type);
            });
        }

        template = [
            {
                input: parameters && parameters.length > 0 ? inputObj : "input_string",
                expected_output: returnType ? getDefaultValueForType(returnType) : "expected_value",
                is_sample: false
            }
        ];
        setJsonInput(JSON.stringify(template, null, 2));
        setError(null);
    };

    const handleImport = async () => {
        try {
            const parsed = JSON.parse(jsonInput);
            let testCasesToCreate: any[] = [];

            if (Array.isArray(parsed)) {
                // Format 1 and 2: Array of objects
                testCasesToCreate = parsed.map((item: any) => {
                    let input = item.input;
                    if (parameters && parameters.length > 0 && typeof input === 'object' && input !== null && !Array.isArray(input)) {
                        // Map object keys to parameter array
                        const inputValues = parameters.map(p => input[p.name]);
                        input = JSON.stringify(inputValues);
                    } else if (typeof input === 'object' && input !== null) {
                        input = JSON.stringify(input);
                    } else if (typeof input === 'string') {
                        // Already a string, assume it's correct or needs [ ] wrapping if not
                        try {
                            const p = JSON.parse(input);
                            if (!Array.isArray(p)) {
                                input = JSON.stringify([p]);
                            }
                        } catch {
                            input = JSON.stringify([input]);
                        }
                    }

                    return {
                        input: String(input),
                        expected_output: JSON.stringify(item.expected_output),
                        is_sample: !!item.is_sample,
                        is_hidden: !!item.is_hidden,
                    };
                });
            } else if (parsed && typeof parsed === 'object') {
                if (parsed.parameters && Array.isArray(parsed.rows)) {
                    // Format 3: Parameter-based rows
                    const params = parsed.parameters as string[];
                    const expectedOutputIdx = params.indexOf('expected_output');
                    const isSampleIdx = params.indexOf('is_sample');
                    const isHiddenIdx = params.indexOf('is_hidden');

                    testCasesToCreate = parsed.rows.map((row: any[]) => {
                        const inputValues: any[] = [];
                        params.forEach((param, idx) => {
                            if (param !== 'expected_output' && param !== 'is_sample' && param !== 'is_hidden') {
                                inputValues.push(row[idx]);
                            }
                        });

                        const rawExpectedOutput = expectedOutputIdx !== -1 ? row[expectedOutputIdx] : "";
                        const expected_output = JSON.stringify(rawExpectedOutput);

                        const is_sample = isSampleIdx !== -1 ? !!row[isSampleIdx] : false;
                        const is_hidden = isHiddenIdx !== -1 ? !!row[isHiddenIdx] : false;

                        let finalInput: string;
                        if (parameters && parameters.length > 0) {
                            finalInput = JSON.stringify(inputValues);
                        } else {
                            // Legacy fallback or non-parameter problems
                            const rawInput = row[0];
                            finalInput = Array.isArray(rawInput) ? JSON.stringify(rawInput) : JSON.stringify([rawInput]);
                        }

                        return {
                            input: finalInput,
                            expected_output,
                            is_sample,
                            is_hidden,
                        };
                    });
                } else {
                    throw new Error('Invalid JSON structure. Expected an array of objects or an object with "parameters" and "rows".');
                }
            } else {
                throw new Error('Invalid JSON structure.');
            }

            // Basic validation
            if (testCasesToCreate.length === 0) {
                throw new Error('No test cases found in JSON.');
            }

            for (const tc of testCasesToCreate) {
                if (!tc.expected_output) {
                    throw new Error('Each test case must have an expected_output.');
                }
            }

            await onImport(testCasesToCreate);
            setJsonInput('');
            onClose();
        } catch (e: any) {
            setError(e.message);
        }
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
            <DialogTitle>Bulk Add Test Cases</DialogTitle>
            <DialogContent>
                <Stack spacing={2} sx={{ mt: 1 }}>
                    <Typography variant="body2" color="text.secondary">
                        Paste a JSON array of test cases or use the template generator.
                    </Typography>

                    <Box display="flex" justifyContent="flex-end">
                        <Button size="small" variant="outlined" onClick={generateTemplate}>
                            Generate Template
                        </Button>
                    </Box>

                    <TextField
                        multiline
                        minRows={10}
                        maxRows={20}
                        fullWidth
                        placeholder='[ { "input": "...", "expected_output": "...", "is_sample": false } ]'
                        value={jsonInput}
                        onChange={(e) => {
                            setJsonInput(e.target.value);
                            setError(null);
                        }}
                        error={!!error}
                        helperText={error}
                        sx={{ fontFamily: 'monospace' }}
                    />

                    {error && (
                        <Alert severity="error">{error}</Alert>
                    )}

                    {isImporting && (
                        <Box display="flex" alignItems="center" gap={2}>
                            <CircularProgress size={20} />
                            <Typography variant="body2">Importing test cases...</Typography>
                        </Box>
                    )}
                </Stack>
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} disabled={isImporting}>
                    Cancel
                </Button>
                <Button
                    onClick={handleImport}
                    variant="contained"
                    disabled={!jsonInput || isImporting}
                >
                    Import {isImporting ? <CircularProgress size={20} sx={{ ml: 1, color: 'white' }} /> : ''}
                </Button>
            </DialogActions>
        </Dialog>
    );
};

export default BulkTestCaseDialog;
