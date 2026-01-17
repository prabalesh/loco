import { useEffect } from 'react'
import Editor, { loader } from '@monaco-editor/react'
import { RotateCcw } from 'lucide-react'
import { Button } from '@/shared/components/ui/Button'
import { LanguageDropdown } from './LanguageDropdown'
import type { ProblemLanguage } from '../types'

interface CodeEditorProps {
    languages: ProblemLanguage[]
    selectedLang: number | null
    currentLang?: ProblemLanguage
    code: string
    onLanguageChange: (langId: number) => void
    onCodeChange: (code: string) => void
    onResetCode: () => void
}

export const CodeEditor = ({
    languages,
    selectedLang,
    currentLang,
    code,
    onLanguageChange,
    onCodeChange,
    onResetCode
}: CodeEditorProps) => {
    // Initialize custom Monaco theme
    useEffect(() => {
        loader.init().then(monaco => {
            monaco.editor.defineTheme('locoCustom', {
                base: 'vs-dark',
                inherit: true,
                rules: [
                    { token: 'comment', foreground: '6A9955', fontStyle: 'italic' },
                    { token: 'keyword', foreground: '569CD6', fontStyle: 'bold' },
                    { token: 'string', foreground: 'CE9178' },
                    { token: 'number', foreground: 'B5CEA8' },
                    { token: 'function', foreground: 'DCDCAA' },
                    { token: 'variable', foreground: '9CDCFE' },
                ],
                colors: {
                    'editor.background': '#0a0a0a',
                    'editor.foreground': '#d4d4d4',
                    'editorLineNumber.foreground': '#4a5568',
                    'editorLineNumber.activeForeground': '#3b82f6',
                    'editor.selectionBackground': '#264f78',
                    'editor.inactiveSelectionBackground': '#3a3d41',
                    'editorCursor.foreground': '#3b82f6',
                    'editor.lineHighlightBackground': '#1a1a1a',
                    'editorIndentGuide.background': '#1a1a1a',
                    'editorIndentGuide.activeBackground': '#2a2a2a',
                }
            })
        })
    }, [])

    return (
        <section className="flex-1 flex flex-col bg-gradient-to-br from-gray-950 to-black font-mono shadow-2xl">
            {/* Code Editor Toolbar */}
            <div className="bg-gray-900/95 backdrop-blur-xl border-b border-white/10 px-4 py-3 flex items-center justify-between z-10 shadow-2xl">
                <div className="flex items-center gap-3">
                    <LanguageDropdown
                        languages={languages}
                        selectedLang={selectedLang}
                        onLanguageChange={onLanguageChange}
                    />
                </div>
                <div className="flex items-center gap-2">
                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={onResetCode}
                        className="text-gray-400 hover:text-white hover:bg-gray-800/80 rounded-xl px-4 py-2 h-auto transition-all duration-200 hover:shadow-md"
                    >
                        <RotateCcw className="h-3.5 w-3.5 mr-2 opacity-70" />
                        <span className="text-[10px] uppercase font-bold tracking-widest">Reset</span>
                    </Button>
                </div>
            </div>

            {/* Editor Container */}
            <div className="flex-1 relative overflow-hidden">
                <Editor
                    height="100%"
                    defaultLanguage={currentLang?.language_name.toLowerCase() || 'javascript'}
                    language={currentLang?.language_name.toLowerCase() || 'javascript'}
                    theme="locoCustom"
                    value={code}
                    onChange={(val) => onCodeChange(val || '')}
                    options={{
                        fontSize: 15,
                        fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace",
                        fontLigatures: true,
                        fontWeight: '500',
                        lineHeight: 22,
                        minimap: { enabled: false },
                        padding: { top: 24, bottom: 24 },
                        smoothScrolling: true,
                        cursorBlinking: 'expand',
                        cursorSmoothCaretAnimation: 'on',
                        lineNumbersMinChars: 3,
                        lineNumbers: 'on',
                        scrollBeyondLastLine: false,
                        automaticLayout: true,
                        letterSpacing: 0.5,
                        foldingStrategy: 'indentation',
                        scrollbar: {
                            vertical: 'auto',
                            horizontal: 'auto',
                            useShadows: false,
                            verticalScrollbarSize: 10,
                            horizontalScrollbarSize: 10,
                        },
                        overviewRulerBorder: false,
                        hideCursorInOverviewRuler: true,
                        renderLineHighlight: 'all',
                        roundedSelection: true,
                        bracketPairColorization: {
                            enabled: true,
                        },
                        guides: {
                            indentation: true,
                            bracketPairs: true,
                        },
                    }}
                />
            </div>
        </section>
    )
}
