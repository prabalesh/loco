import React, { useState, useEffect } from 'react';
import {
    Stack,
    Typography,
    Box,
    Button,
    IconButton,
    Paper,
    Alert,
    Chip,
    CircularProgress,
    Collapse
} from '@mui/material';
import {
    Refresh as RefreshIcon,
    CheckCircle as SuccessIcon,
    ExpandMore,
    ExpandLess,
    Language as LanguageIcon,
    Delete as DeleteIcon
} from '@mui/icons-material';
import { adminLanguagesApi, adminCodeGenApi, adminProblemApi } from '../../../../lib/api/admin';

interface BoilerplateStepProps {
    data: any;
    onChange: (newData: Partial<any>) => void;
    onRefresh: () => void;
}

export const BoilerplateStep: React.FC<BoilerplateStepProps> = ({ data, onChange, onRefresh }) => {
    const [languages, setLanguages] = useState<any[]>([]);
    const [generating, setGenerating] = useState<string | null>(null);
    const [batchGenerating, setBatchGenerating] = useState(false);
    const [previews, setPreviews] = useState<Record<string, string>>({});
    const [expanded, setExpanded] = useState<string | null>(null);

    const selectedLanguages = data?.selected_languages || [];

    useEffect(() => {
        if (data.boilerplates && data.boilerplates.length > 0) {
            const initialPreviews: Record<string, string> = {};
            data.boilerplates.forEach((b: any) => {
                const lang = b.language?.language_id || languages.find(l => l.id === b.language_id)?.language_id;
                if (lang) {
                    initialPreviews[lang] = b.stub_code;
                }
            });
            setPreviews(initialPreviews);
        }
    }, [data.boilerplates, languages]);

    useEffect(() => {
        fetchLanguages();
    }, []);

    const fetchLanguages = async () => {
        try {
            const response = await adminLanguagesApi.getAllActive();
            setLanguages(response.data.data);
        } catch (err) {
            console.error('Failed to fetch languages', err);
        }
    };

    const toggleLanguage = (langId: string) => {
        const next = selectedLanguages.includes(langId)
            ? selectedLanguages.filter((s: string) => s !== langId)
            : [...selectedLanguages, langId];
        onChange({ selected_languages: next });
    };

    const handleGeneratePreview = async (langSlug: string) => {
        setGenerating(langSlug);
        try {
            const response = await adminCodeGenApi.generateStub({
                function_name: data.function_name,
                return_type: data.return_type,
                parameters: data.parameters,
                language_slug: langSlug,
            });
            setPreviews(prev => ({ ...prev, [langSlug]: response.data.data.stub_code }));
            setExpanded(langSlug);
        } catch (err) {
            console.error(`Failed to generate stub for ${langSlug}`, err);
        } finally {
            setGenerating(null);
        }
    };

    const handleRegenerateAll = async () => {
        if (!data.id) return;
        setBatchGenerating(true);
        try {
            await adminProblemApi.v2RegenerateBoilerplates(data.id);
            onRefresh();
        } catch (err) {
            console.error('Failed to regenerate all boilerplates', err);
        } finally {
            setBatchGenerating(false);
        }
    };

    return (
        <Stack spacing={4}>
            <Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 3 }}>
                    <Box>
                        <Typography variant="h6" color="primary" fontWeight="bold">Language Support & Boilerplates</Typography>
                        <Typography variant="body2" color="text.secondary">
                            Select the languages you want to support for this problem and preview the generated code stubs.
                        </Typography>
                    </Box>
                    {data.id && (
                        <Button
                            variant="contained"
                            size="small"
                            startIcon={batchGenerating ? <CircularProgress size={16} color="inherit" /> : <RefreshIcon />}
                            onClick={handleRegenerateAll}
                            disabled={batchGenerating}
                        >
                            {batchGenerating ? 'Regenerating...' : 'Regenerate All'}
                        </Button>
                    )}
                </Box>

                <Paper variant="outlined" sx={{ p: 3, mb: 4, borderRadius: 2 }}>
                    <Typography variant="subtitle2" gutterBottom fontWeight="bold">Select Supported Languages</Typography>
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 1 }}>
                        {languages.map(lang => (
                            <Chip
                                key={lang.id}
                                label={lang.name}
                                onClick={() => toggleLanguage(lang.language_id)}
                                color={selectedLanguages.includes(lang.language_id) ? "primary" : "default"}
                                variant={selectedLanguages.includes(lang.language_id) ? "filled" : "outlined"}
                                icon={<LanguageIcon />}
                                sx={{ px: 1, py: 2, borderRadius: '8px' }}
                            />
                        ))}
                    </Box>
                </Paper>

                <Stack spacing={2}>
                    {selectedLanguages.map((slug: string) => {
                        const lang = languages.find(l => l.language_id === slug);
                        if (!lang) return null;

                        console.log("Language", slug);

                        return (
                            <Paper key={slug} variant="outlined" sx={{ overflow: 'hidden', borderRadius: 2 }}>
                                <Box sx={{ p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center', bgcolor: expanded === slug ? '#f0f7ff' : '#fcfcfc' }}>
                                    <Stack direction="row" spacing={2} alignItems="center">
                                        <Typography fontWeight="bold" color={expanded === slug ? 'primary' : 'textPrimary'}>{lang.name}</Typography>
                                        {previews[slug] && <SuccessIcon color="success" fontSize="small" />}
                                    </Stack>
                                    <Stack direction="row" spacing={1}>
                                        <Button
                                            size="small"
                                            startIcon={generating === slug ? <CircularProgress size={16} /> : <RefreshIcon />}
                                            onClick={() => handleGeneratePreview(slug)}
                                            disabled={!!generating}
                                            variant="text"
                                        >
                                            {previews[slug] ? 'Regenerate' : 'Generate Preview'}
                                        </Button>
                                        <Button
                                            size="small"
                                            startIcon={<DeleteIcon />}
                                            onClick={() => toggleLanguage(slug)}
                                            color="error"
                                            variant="text"
                                        >
                                            Remove
                                        </Button>
                                        <IconButton
                                            size="small"
                                            onClick={() => setExpanded(expanded === slug ? null : slug)}
                                            disabled={!previews[slug]}
                                        >
                                            {expanded === slug ? <ExpandLess /> : <ExpandMore />}
                                        </IconButton>
                                    </Stack>
                                </Box>

                                <Collapse in={expanded === slug}>
                                    <Box sx={{ p: 0, bgcolor: '#1e1e1e', color: '#d4d4d4', fontFamily: 'monospace', fontSize: '0.8rem' }}>
                                        <pre style={{ margin: 0, padding: '1.5rem', overflowX: 'auto', maxHeight: '300px' }}>
                                            {previews[slug] || '// No preview generated yet'}
                                        </pre>
                                    </Box>
                                </Collapse>
                            </Paper>
                        );
                    })}
                </Stack>
            </Box>

            {selectedLanguages.length === 0 && (
                <Alert severity="warning" sx={{ borderRadius: 2 }}>Please select at least one language to support.</Alert>
            )}
        </Stack>
    );
};
