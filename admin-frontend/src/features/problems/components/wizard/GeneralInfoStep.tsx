import React from 'react';
import {
    TextField,
    MenuItem,
    Stack,
    Typography,
    Box,
    Divider,
    Paper
} from '@mui/material';

interface GeneralInfoStepProps {
    data: any;
    onChange: (newData: Partial<any>) => void;
    isConsolidated?: boolean;
}

const DIFFICULTIES = [
    { value: 'easy', label: 'Easy' },
    { value: 'medium', label: 'Medium' },
    { value: 'hard', label: 'Hard' },
];

export const GeneralInfoStep: React.FC<GeneralInfoStepProps> = ({ data, onChange }) => {
    return (
        <Stack spacing={4}>
            <Box>
                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">Basic Information</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                    Define the title, difficulty, and unique identifier for the problem.
                </Typography>

                <Stack spacing={3}>
                    <div className="grid grid-cols-1 sm:grid-cols-12 gap-6">
                        <div className="sm:col-span-8">
                            <TextField
                                fullWidth
                                label="Problem Title"
                                value={data.title}
                                onChange={(e) => onChange({ title: e.target.value })}
                                required
                                placeholder="e.g., Two Sum"
                            />
                        </div>
                        <div className="sm:col-span-4">
                            <TextField
                                select
                                fullWidth
                                label="Difficulty"
                                value={data.difficulty}
                                onChange={(e) => onChange({ difficulty: e.target.value })}
                            >
                                {DIFFICULTIES.map((opt) => (
                                    <MenuItem key={opt.value} value={opt.value}>
                                        {opt.label}
                                    </MenuItem>
                                ))}
                            </TextField>
                        </div>
                    </div>

                    <TextField
                        fullWidth
                        label="URL Slug (Optional)"
                        value={data.slug || ''}
                        onChange={(e) => onChange({ slug: e.target.value })}
                        placeholder="e.g., two-sum"
                        helperText="Leave empty to auto-generate from title"
                    />
                </Stack>
            </Box>

            <Divider />

            <Box>
                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">Problem Content</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                    Describe the problem statement, constraints, and provide hints for users.
                </Typography>

                <Stack spacing={3}>
                    <TextField
                        fullWidth
                        label="Problem Description"
                        value={data.description}
                        onChange={(e) => onChange({ description: e.target.value })}
                        required
                        multiline
                        rows={10}
                        placeholder="### Problem Statement\n\nGiven an array of integers `nums`, return..."
                        sx={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                        helperText="Markdown is supported"
                    />

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <TextField
                            fullWidth
                            label="Constraints"
                            value={data.constraints || ''}
                            onChange={(e) => onChange({ constraints: e.target.value })}
                            multiline
                            rows={4}
                            placeholder="- 1 <= nums.length <= 10^4\n- -10^9 <= nums[i] <= 10^9"
                            sx={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                        />
                        <TextField
                            fullWidth
                            label="Hints (One per line)"
                            value={data.hints || ''}
                            onChange={(e) => onChange({ hints: e.target.value })}
                            multiline
                            rows={4}
                            placeholder="Check for null\nUse a hashmap"
                            sx={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                        />
                    </div>
                </Stack>
            </Box>
        </Stack>
    );
};
