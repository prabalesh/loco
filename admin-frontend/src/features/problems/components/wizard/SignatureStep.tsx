import React, { useState, useEffect } from 'react';
import {
    TextField,
    MenuItem,
    Stack,
    Typography,
    Box,
    Divider,
    Button,
    IconButton,
    Tooltip,
    Paper
} from '@mui/material';
import {
    Add as PlusIcon,
    Delete as DeleteIcon,
    Info as InfoIcon
} from '@mui/icons-material';
import { adminProblemApi } from '../../../../lib/api/admin';

interface SignatureStepProps {
    data: any;
    onChange: (newData: Partial<any>) => void;
}

const PRIMITIVE_TYPES = [
    { value: 'integer', label: 'int', is_custom: false },
    { value: 'integer_array', label: 'int[]', is_custom: false },
    { value: 'string', label: 'string', is_custom: false },
    { value: 'string_array', label: 'string[]', is_custom: false },
    { value: 'boolean', label: 'bool', is_custom: false },
    { value: 'double', label: 'double', is_custom: false },
];

export const SignatureStep: React.FC<SignatureStepProps> = ({ data, onChange }) => {
    const [customTypes, setCustomTypes] = useState<any[]>([]);

    useEffect(() => {
        fetchCustomTypes();
    }, []);

    const fetchCustomTypes = async () => {
        try {
            const response = await adminProblemApi.getCustomTypes();
            setCustomTypes(response.data.data.map((ct: any) => ({
                value: ct.name,
                label: ct.name,
                is_custom: true
            })));
        } catch (err) {
            console.error('Failed to fetch custom types', err);
        }
    };

    const allTypes = [...PRIMITIVE_TYPES, ...customTypes];

    const handleAddParam = () => {
        const newParams = [...data.parameters, { name: '', type: 'integer', is_custom: false }];
        onChange({ parameters: newParams });
    };

    const handleRemoveParam = (index: number) => {
        const newParams = [...data.parameters];
        newParams.splice(index, 1);
        onChange({ parameters: newParams });
    };

    const handleParamChange = (index: number, field: string, value: any) => {
        const newParams = [...data.parameters];
        if (field === 'type') {
            const typeObj = allTypes.find(t => t.value === value);
            newParams[index] = { ...newParams[index], [field]: value, is_custom: typeObj?.is_custom || false };
        } else {
            newParams[index] = { ...newParams[index], [field]: value };
        }
        onChange({ parameters: newParams });
    };

    return (
        <Stack spacing={4}>
            <Box>
                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">Function Signature</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                    Define the interface of the solution function. This structure will be used across all supported languages.
                </Typography>

                <Paper variant="outlined" sx={{ p: 4, bgcolor: '#fafafa', borderRadius: 2 }}>
                    <Stack spacing={4}>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                            <TextField
                                fullWidth
                                label="Function Name"
                                value={data.function_name}
                                onChange={(e) => onChange({ function_name: e.target.value })}
                                required
                                placeholder="e.g., solve, twoSum"
                            />
                            <TextField
                                select
                                fullWidth
                                label="Return Type"
                                value={data.return_type}
                                onChange={(e) => onChange({ return_type: e.target.value })}
                            >
                                {allTypes.map((opt) => (
                                    <MenuItem key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </MenuItem>
                                ))}
                            </TextField>
                        </div>

                        <Divider />

                        <Box>
                            <Typography variant="subtitle2" gutterBottom sx={{ display: 'flex', alignItems: 'center', gap: 1, fontWeight: 'bold' }}>
                                Input Parameters
                                <Tooltip title="Specify the input arguments for the function.">
                                    <InfoIcon fontSize="small" color="action" />
                                </Tooltip>
                            </Typography>

                            <Stack spacing={2} sx={{ mt: 2 }}>
                                {data.parameters.map((param: any, index: number) => (
                                    <Box key={index} sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                                        <TextField
                                            sx={{ flex: 1 }}
                                            label="Parameter Name"
                                            size="small"
                                            value={param.name}
                                            onChange={(e) => handleParamChange(index, 'name', e.target.value)}
                                            required
                                            placeholder="e.g., nums"
                                        />
                                        <TextField
                                            select
                                            sx={{ width: 180 }}
                                            label="Type"
                                            size="small"
                                            value={param.type}
                                            onChange={(e) => handleParamChange(index, 'type', e.target.value)}
                                        >
                                            {allTypes.map((opt) => (
                                                <MenuItem key={opt.value} value={opt.value}>
                                                    {opt.label}
                                                </MenuItem>
                                            ))}
                                        </TextField>
                                        <IconButton
                                            color="error"
                                            onClick={() => handleRemoveParam(index)}
                                            disabled={data.parameters.length === 1}
                                            size="small"
                                        >
                                            <DeleteIcon fontSize="small" />
                                        </IconButton>
                                    </Box>
                                ))}

                                <Button
                                    startIcon={<PlusIcon />}
                                    onClick={handleAddParam}
                                    sx={{ alignSelf: 'flex-start', mt: 1 }}
                                    size="small"
                                >
                                    Add Parameter
                                </Button>
                            </Stack>
                        </Box>
                    </Stack>
                </Paper>
            </Box>

            <Box>
                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">Validation Strategy</Typography>
                <Paper variant="outlined" sx={{ p: 3, borderRadius: 2 }}>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                        <TextField
                            select
                            fullWidth
                            label="Comparison Method"
                            value={data.validation_type}
                            onChange={(e) => onChange({ validation_type: e.target.value })}
                            helperText="How should user output be compared to expected?"
                        >
                            <MenuItem value="EXACT">Exact Match</MenuItem>
                            <MenuItem value="UNORDERED">Unordered Array (Sets)</MenuItem>
                            <MenuItem value="SUBSET">Subset Match</MenuItem>
                        </TextField>

                        <Box sx={{ display: 'flex', gap: 2 }}>
                            <TextField
                                sx={{ flex: 1 }}
                                label="Time Limit (ms)"
                                type="number"
                                defaultValue={2000}
                            />
                            <TextField
                                sx={{ flex: 1 }}
                                label="Memory Limit (MB)"
                                type="number"
                                defaultValue={256}
                            />
                        </Box>
                    </div>
                </Paper>
            </Box>
        </Stack>
    );
};
