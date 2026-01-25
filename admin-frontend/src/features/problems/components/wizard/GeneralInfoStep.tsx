import React from 'react';
import {
    TextField,
    MenuItem,
    Stack,
    Typography,
    Box,
    Divider,
    FormControlLabel,
    Checkbox,
    Chip
} from '@mui/material';
import TiptapEditor from '../../../../components/editor/TiptapEditor';
import { useQuery } from '@tanstack/react-query';
import { adminProblemApi } from '../../../../lib/api/admin';
import type { Tag, Category } from '../../../../types';

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
    const { data: tags = [] } = useQuery({
        queryKey: ['tags'],
        queryFn: () => adminProblemApi.getTags().then(res => res.data.data || []),
    });

    const { data: categories = [] } = useQuery({
        queryKey: ['categories'],
        queryFn: () => adminProblemApi.getCategories().then(res => res.data.data || []),
    });
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

                <Stack spacing={4}>
                    <Box>
                        <Typography variant="body2" fontWeight="bold" gutterBottom>
                            Problem Description <span style={{ color: '#ef4444' }}>*</span>
                        </Typography>
                        <TiptapEditor
                            content={data.description}
                            onChange={(content) => onChange({ description: content })}
                        />
                        <Typography variant="caption" color="text.secondary">
                            Markdown and HTML are supported
                        </Typography>
                    </Box>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <Box>
                            <Typography variant="body2" fontWeight="bold" gutterBottom>
                                Constraints
                            </Typography>
                            <TiptapEditor
                                content={data.constraints || ''}
                                onChange={(content) => onChange({ constraints: content })}
                            />
                        </Box>
                        <Box>
                            <Typography variant="body2" fontWeight="bold" gutterBottom>
                                Hints (One per line)
                            </Typography>
                            <TextField
                                fullWidth
                                value={data.hints || ''}
                                onChange={(e) => onChange({ hints: e.target.value })}
                                multiline
                                rows={6}
                                placeholder="Check for null\nUse a hashmap"
                                sx={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                            />
                        </Box>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <Box>
                            <Typography variant="body2" fontWeight="bold" gutterBottom>
                                Input Format
                            </Typography>
                            <TiptapEditor
                                content={data.input_format || ''}
                                onChange={(content) => onChange({ input_format: content })}
                            />
                        </Box>
                        <Box>
                            <Typography variant="body2" fontWeight="bold" gutterBottom>
                                Output Format
                            </Typography>
                            <TiptapEditor
                                content={data.output_format || ''}
                                onChange={(content) => onChange({ output_format: content })}
                            />
                        </Box>
                    </div>
                </Stack>
            </Box>

            <Divider />

            <Box>
                <Typography variant="h6" gutterBottom color="primary" fontWeight="bold">Classification</Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                    Assign appropriate tags and categories to help users find this problem.
                </Typography>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div>
                        <Typography variant="subtitle2" fontWeight="bold" gutterBottom>Tags</Typography>
                        <Box sx={{ p: 2, border: '1px solid #e0e0e0', borderRadius: 1, minHeight: 120, maxHeight: 250, overflowY: 'auto', bgcolor: '#fafafa' }}>
                            <div className="flex flex-wrap gap-2">
                                {Array.isArray(tags) && tags.map((tag: Tag) => (
                                    <Chip
                                        key={tag.id}
                                        label={tag.name}
                                        onClick={() => {
                                            const currentIds = data.tag_ids || [];
                                            const nextIds = currentIds.includes(tag.id)
                                                ? currentIds.filter((id: number) => id !== tag.id)
                                                : [...currentIds, tag.id];
                                            onChange({ tag_ids: nextIds });
                                        }}
                                        color={(data.tag_ids || []).includes(tag.id) ? "primary" : "default"}
                                        variant={(data.tag_ids || []).includes(tag.id) ? "filled" : "outlined"}
                                        size="small"
                                        sx={{ cursor: 'pointer' }}
                                    />
                                ))}
                            </div>
                        </Box>
                    </div>

                    <div>
                        <Typography variant="subtitle2" fontWeight="bold" gutterBottom>Categories</Typography>
                        <Box sx={{ p: 2, border: '1px solid #e0e0e0', borderRadius: 1, minHeight: 120, maxHeight: 250, overflowY: 'auto', bgcolor: '#fafafa' }}>
                            <Stack spacing={0.5}>
                                {Array.isArray(categories) && categories.map((cat: Category) => (
                                    <FormControlLabel
                                        key={cat.id}
                                        control={
                                            <Checkbox
                                                size="small"
                                                checked={(data.category_ids || []).includes(cat.id)}
                                                onChange={(e) => {
                                                    const currentIds = data.category_ids || [];
                                                    const nextIds = e.target.checked
                                                        ? [...currentIds, cat.id]
                                                        : currentIds.filter((id: number) => id !== cat.id);
                                                    onChange({ category_ids: nextIds });
                                                }}
                                            />
                                        }
                                        label={<Typography variant="body2">{cat.name}</Typography>}
                                    />
                                ))}
                            </Stack>
                        </Box>
                    </div>
                </div>
            </Box>
        </Stack >
    );
};
