import { useParams, useNavigate } from 'react-router-dom'
import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
  Card,
  Button,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Tabs,
  Tab,
  Alert,
  CircularProgress,
  Chip,
  Stack,
  Box,
  Typography,
  CardHeader,
  CardContent,
  Divider,
  AlertTitle
} from '@mui/material'
import {
  Save as SaveIcon,
  ArrowBack as ArrowBackIcon,
  ArrowForward as ArrowForwardIcon
} from '@mui/icons-material'
import Editor from '@monaco-editor/react'
import toast from 'react-hot-toast'
import { adminLanguagesApi, adminProblemLanguagesApi } from '../../../lib/api/admin'
import { ProblemStepper } from '../components/ProblemStepper'

interface Language {
  id: number
  name: string
  language_id: string
  extension: string
  default_template?: string
}

interface ProblemLanguage {
  problem_id: number
  language_id: number
  function_code: string
  main_code: string
  solution_code: string
  language?: Language
}

export default function ProblemLanguage() {
  const { problemId } = useParams<{ problemId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [selectedLanguageId, setSelectedLanguageId] = useState<number | null>(null)
  const [functionCode, setFunctionCode] = useState('')
  const [mainCode, setMainCode] = useState('')
  const [solutionCode, setSolutionCode] = useState('')
  const [activeTab, setActiveTab] = useState('function')

  // Fetch all available languages
  const { data: languagesData, isLoading: languagesLoading } = useQuery({
    queryKey: ['languages', 'active'],
    queryFn: () => adminLanguagesApi.getAllActive(),
  })

  // Fetch problem languages
  const { data: problemLanguagesData, isLoading: problemLanguagesLoading } = useQuery({
    queryKey: ['problem-languages', problemId],
    queryFn: () => adminProblemLanguagesApi.getAll(String(problemId)),
  })

  const languages = languagesData?.data?.data || []
  const problemLanguages = problemLanguagesData?.data?.data || []

  // Get languages not yet added to this problem
  const availableLanguages = languages.filter(
    (lang) => !problemLanguages.some((pl) => pl.language_id === lang.id)
  )

  // Save/Update mutation
  const saveMutation = useMutation({
    mutationFn: async (data: { language_id: number; function_code: string; main_code: string; solution_code: string }) => {
      const existing = problemLanguages.find((pl) => pl.language_id === data.language_id)
      if (existing) {
        return adminProblemLanguagesApi.update(String(problemId), data.language_id, data)
      } else {
        return adminProblemLanguagesApi.create(String(problemId), data)
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['problem-languages', problemId] })
      toast.success('Language configuration saved successfully!')
    },
    onError: () => {
      toast.error('Failed to save language configuration')
    },
  })

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: (languageId: number) => adminProblemLanguagesApi.delete(String(problemId), languageId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['problem-languages', problemId] })
      toast.success('Language removed successfully!')
      if (selectedLanguageId === deleteMutation.variables) {
        setSelectedLanguageId(null)
        resetEditor()
      }
    },
    onError: () => {
      toast.error('Failed to remove language')
    },
  })

  const resetEditor = () => {
    setFunctionCode('')
    setMainCode('')
    setSolutionCode('')
    setActiveTab('function')
  }

  // Load language data when selected
  useEffect(() => {
    if (selectedLanguageId) {
      const problemLang = problemLanguages.find((pl) => pl.language_id === selectedLanguageId)
      const lang = languages.find((l) => l.id === selectedLanguageId)

      if (problemLang) {
        // Load existing configuration
        setFunctionCode(problemLang.function_code || '')
        setMainCode(problemLang.main_code || '')
        setSolutionCode(problemLang.solution_code || '')
      } else if (lang) {
        // Initialize with template
        setFunctionCode(lang.default_template || '// Write starter code here')
        setMainCode('// Write I/O handling code here')
        setSolutionCode('// Write your solution here')
      }
    }
  }, [selectedLanguageId, problemLanguages, languages])

  const handleSave = () => {
    if (!selectedLanguageId) {
      toast.error('Please select a language')
      return
    }

    if (!functionCode.trim() || !mainCode.trim() || !solutionCode.trim()) {
      toast.error('All code sections are required')
      return
    }

    saveMutation.mutate({
      language_id: selectedLanguageId,
      function_code: functionCode,
      main_code: mainCode,
      solution_code: solutionCode,
    })
  }

  const handleDelete = (languageId: number) => {
    if (window.confirm('Are you sure you want to remove this language?')) {
      deleteMutation.mutate(languageId)
    }
  }

  const handleLanguageSelect = (langId: number) => {
    setSelectedLanguageId(langId)
  }

  const getLanguageName = (langId: number) => {
    return languages.find((l) => l.id === langId)?.name || 'Unknown'
  }

  const getEditorLanguage = (langId: number) => {
    let lang = languages.find((l) => l.id === langId)?.language_id;

    switch (lang) {
      case "c++":
        lang = "cpp"
        break;
    }

    return lang || 'plaintext'
  }

  if (languagesLoading || problemLanguagesLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
        <CircularProgress size={60} />
      </Box>
    )
  }

  return (
    <Stack spacing={4} sx={{ p: 4 }}>
      {/* Header */}
      <Box textAlign="center">
        <Typography variant="h4" fontWeight="bold" gutterBottom>
          Configure Languages
        </Typography>
        <ProblemStepper currentStep={3} model="validate" problemId={problemId || "create"} />
      </Box>

      {/* Language Selection */}
      <Card variant="outlined">
        <CardHeader title="Select Language" />
        <Divider />
        <CardContent>
          <Stack spacing={3}>
            <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
              <FormControl fullWidth sx={{ maxWidth: 300 }}>
                <InputLabel id="language-add-label">Add a language</InputLabel>
                <Select
                  labelId="language-add-label"
                  label="Add a language"
                  value=""
                  onChange={(e) => handleLanguageSelect(Number(e.target.value))}
                  disabled={availableLanguages.length === 0}
                >
                  {availableLanguages.map((lang) => (
                    <MenuItem key={lang.id} value={lang.id}>
                      {lang.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Typography variant="body2" color="text.secondary">
                {availableLanguages.length === 0 ? 'All languages added' : `${availableLanguages.length} available`}
              </Typography>
            </Box>

            {/* Added Languages */}
            {problemLanguages.length > 0 && (
              <Box>
                <Typography variant="subtitle2" gutterBottom color="text.secondary">
                  Configured Languages:
                </Typography>
                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                  {problemLanguages.map((pl) => (
                    <Chip
                      key={pl.language_id}
                      label={getLanguageName(pl.language_id)}
                      color={selectedLanguageId === pl.language_id ? 'primary' : 'default'}
                      variant={selectedLanguageId === pl.language_id ? 'filled' : 'outlined'}
                      onDelete={() => handleDelete(pl.language_id)}
                      onClick={() => handleLanguageSelect(pl.language_id)}
                      sx={{ cursor: 'pointer' }}
                    />
                  ))}
                </Stack>
              </Box>
            )}
          </Stack>
        </CardContent>
      </Card>

      {/* Code Editor */}
      {selectedLanguageId ? (
        <Card variant="outlined">
          <CardHeader
            title={`${getLanguageName(selectedLanguageId)} Configuration`}
            action={
              <Button
                variant="contained"
                startIcon={<SaveIcon />}
                onClick={handleSave}
                disabled={saveMutation.isPending}
              >
                Save
              </Button>
            }
          />
          <Divider />
          <CardContent>
            <Alert severity="info" sx={{ mb: 3 }}>
              <AlertTitle>Code Sections</AlertTitle>
              Function: Starter code shown to users | Main: I/O handling | Solution: Reference solution
            </Alert>

            <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
              <Tabs value={activeTab} onChange={(_e, v) => setActiveTab(v)}>
                <Tab label="Function Code" value="function" />
                <Tab label="Main Code" value="main" />
                <Tab label="Solution Code" value="solution" />
              </Tabs>
            </Box>

            <Box sx={{ border: 1, borderColor: 'divider', borderRadius: 1, overflow: 'hidden' }}>
              <Editor
                height="400px"
                language={getEditorLanguage(selectedLanguageId)}
                value={activeTab === 'function' ? functionCode : activeTab === 'main' ? mainCode : solutionCode}
                onChange={(value) => {
                  if (activeTab === 'function') setFunctionCode(value || '')
                  else if (activeTab === 'main') setMainCode(value || '')
                  else setSolutionCode(value || '')
                }}
                theme="vs-dark"
                options={{
                  minimap: { enabled: false },
                  fontSize: 14,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                }}
              />
            </Box>
          </CardContent>
        </Card>
      ) : (
        <Card variant="outlined">
          <CardContent>
            <Typography color="text.secondary" align="center" py={4} variant="h6">
              Select a language to configure boilerplate code
            </Typography>
          </CardContent>
        </Card>
      )}

      {/* Navigation */}
      <Box display="flex" justifyContent="space-between" mt={4}>
        <Button
          variant="outlined"
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate(`/problems/${problemId}/testcases`)}
        >
          Back to Test Cases
        </Button>
        <Button
          variant="contained"
          endIcon={<ArrowForwardIcon />}
          onClick={() => navigate(`/problems/${problemId}/validate`)}
          disabled={problemLanguages.length === 0}
        >
          Continue
        </Button>
      </Box>
    </Stack>
  )
}
