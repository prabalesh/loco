import type { Editor } from '@tiptap/react'
import { EditorContent, useEditor } from '@tiptap/react'
import { useEffect, useState } from 'react'
import StarterKit from '@tiptap/starter-kit'
import TextAlign from '@tiptap/extension-text-align'
import Highlight from '@tiptap/extension-highlight'
import Code from '@tiptap/extension-code'
import { CodeBlockLowlight } from '@tiptap/extension-code-block-lowlight'
import { Markdown } from 'tiptap-markdown'
import { common, createLowlight } from 'lowlight'
import styles from './TiptapEditor.module.css'

import { 
  AlignCenter, 
  AlignLeft, 
  AlignRight, 
  Bold, 
  Heading1, 
  Heading2, 
  Heading3, 
  Highlighter, 
  Italic, 
  ListOrdered, 
  Strikethrough, 
  List,
  Code2,
  FileCode,
  FileText,
  Eye
} from 'lucide-react'

interface TiptapEditorProps {
  content: string
  onChange?: (content: string) => void
}

const lowlight = createLowlight(common)

function MenuBar({ editor, isRawMode, onToggleMode }: { 
  editor: Editor | null
  isRawMode: boolean
  onToggleMode: () => void 
}) {
  if (!editor) {
    return null
  }

  const Options = [
    {
      icon: <Heading1 className="size-4" />,
      onClick: () => editor.chain().focus().toggleHeading({ level: 1 }).run(),
      pressed: editor.isActive("heading", { level: 1 }),
      disabled: isRawMode,
    },
    {
      icon: <Heading2 className="size-4" />,
      onClick: () => editor.chain().focus().toggleHeading({ level: 2 }).run(),
      pressed: editor.isActive("heading", { level: 2 }),
      disabled: isRawMode,
    },
    {
      icon: <Heading3 className="size-4" />,
      onClick: () => editor.chain().focus().toggleHeading({ level: 3 }).run(),
      pressed: editor.isActive("heading", { level: 3 }),
      disabled: isRawMode,
    },
    {
      icon: <Bold className="size-4" />,
      onClick: () => editor.chain().focus().toggleBold().run(),
      pressed: editor.isActive("bold"),
      disabled: isRawMode,
    },
    {
      icon: <Italic className="size-4" />,
      onClick: () => editor.chain().focus().toggleItalic().run(),
      pressed: editor.isActive("italic"),
      disabled: isRawMode,
    },
    {
      icon: <Strikethrough className="size-4" />,
      onClick: () => editor.chain().focus().toggleStrike().run(),
      pressed: editor.isActive("strike"),
      disabled: isRawMode,
    },
    {
      icon: <Code2 className="size-4" />,
      onClick: () => editor.chain().focus().toggleCode().run(),
      pressed: editor.isActive("code"),
      disabled: isRawMode,
    },
    {
      icon: <FileCode className="size-4" />,
      onClick: () => editor.chain().focus().toggleCodeBlock().run(),
      pressed: editor.isActive("codeBlock"),
      disabled: isRawMode,
    },
    {
      icon: <AlignLeft className="size-4" />,
      onClick: () => editor.chain().focus().setTextAlign("left").run(),
      pressed: editor.isActive({ textAlign: "left" }),
      disabled: isRawMode,
    },
    {
      icon: <AlignCenter className="size-4" />,
      onClick: () => editor.chain().focus().setTextAlign("center").run(),
      pressed: editor.isActive({ textAlign: "center" }),
      disabled: isRawMode,
    },
    {
      icon: <AlignRight className="size-4" />,
      onClick: () => editor.chain().focus().setTextAlign("right").run(),
      pressed: editor.isActive({ textAlign: "right" }),
      disabled: isRawMode,
    },
    {
      icon: <List className="size-4" />,
      onClick: () => editor.chain().focus().toggleBulletList().run(),
      pressed: editor.isActive("bulletList"),
      disabled: isRawMode,
    },
    {
      icon: <ListOrdered className="size-4" />,
      onClick: () => editor.chain().focus().toggleOrderedList().run(),
      pressed: editor.isActive("orderedList"),
      disabled: isRawMode,
    },
    {
      icon: <Highlighter className="size-4" />,
      onClick: () => editor.chain().focus().toggleHighlight().run(),
      pressed: editor.isActive("highlight"),
      disabled: isRawMode,
    },
  ]

  return (
    <div className="border-b border-gray-200 bg-gray-50">
      <div className="flex flex-wrap items-center justify-between gap-1 p-2">
        <div className="flex flex-wrap items-center gap-1">
          {Options.map((option, index) => (
            <button
              key={index}
              onClick={option.onClick}
              type="button"
              disabled={option.disabled}
              className={`
                rounded-md p-2 transition-colors
                ${option.disabled ? 'opacity-40 cursor-not-allowed' : 'hover:bg-gray-200'}
                ${option.pressed 
                  ? 'bg-blue-100 text-blue-700 hover:bg-blue-200' 
                  : 'text-gray-700'
                }
              `}
            >
              {option.icon}
            </button>
          ))}
        </div>
        
        {/* Toggle Raw/Visual Mode */}
        <button
          type="button"
          onClick={onToggleMode}
          className={`
            flex items-center gap-2 px-3 py-1.5 rounded-md text-sm font-medium transition-colors
            ${isRawMode 
              ? 'bg-purple-100 text-purple-700 hover:bg-purple-200' 
              : 'bg-blue-100 text-blue-700 hover:bg-blue-200'
            }
          `}
        >
          {isRawMode ? (
            <>
              <Eye className="size-4" />
              Visual
            </>
          ) : (
            <>
              <FileText className="size-4" />
              HTML
            </>
          )}
        </button>
      </div>
    </div>
  )
}

export default function TiptapEditor({ content, onChange }: TiptapEditorProps) {
  const [isRawMode, setIsRawMode] = useState(false)
  const [rawContent, setRawContent] = useState('')

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        codeBlock: false,
        bulletList: {
          keepMarks: true,
          keepAttributes: false,
        },
        orderedList: {
          keepMarks: true,
          keepAttributes: false,
        },
      }),
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      }),
      Highlight,
      Code.configure({
        HTMLAttributes: {
          class: 'inline-code',
        },
      }),
      CodeBlockLowlight.configure({
        lowlight,
        HTMLAttributes: {
          class: 'code-block',
        },
      }),
      Markdown.configure({
        html: true,
        transformPastedText: true,
        transformCopiedText: true,
      }),
    ],
    content: content || '',
    editorProps: {
      attributes: {
        class: `${styles.editor} focus:outline-none p-4 min-h-[300px]`
      }
    },
    onUpdate: ({ editor }) => {
      if (!isRawMode) {
        onChange?.(editor.getHTML())
      }
    },
    immediatelyRender: false,
  })

  // Update editor content when prop changes externally
  useEffect(() => {
    if (editor && content !== editor.getHTML() && !isRawMode) {
      editor.commands.setContent(content)
    }
  }, [content, editor, isRawMode])

  // Initialize raw content when switching to raw mode
  useEffect(() => {
    if (isRawMode && editor) {
      setRawContent(editor.getHTML())
    }
  }, [isRawMode, editor])

  const handleToggleMode = () => {
    if (isRawMode && editor) {
      // Switching from raw to visual - apply raw content to editor
      try {
        editor.commands.setContent(rawContent)
        onChange?.(rawContent)
      } catch (error) {
        console.error('Error parsing HTML:', error)
        alert('Invalid HTML content. Please check your markup.')
        return
      }
    } else if (editor) {
      // Switching from visual to raw - get current HTML
      setRawContent(editor.getHTML())
    }
    setIsRawMode(!isRawMode)
  }

  const handleRawContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setRawContent(e.target.value)
  }

  return (
    <div className="rounded-lg border border-gray-300 bg-white shadow-sm overflow-hidden">
      <MenuBar editor={editor} isRawMode={isRawMode} onToggleMode={handleToggleMode} />
      
      {isRawMode ? (
        <div className="p-4">
          <textarea
            value={rawContent}
            onChange={handleRawContentChange}
            className="w-full min-h-[300px] p-4 font-mono text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Paste your HTML or Markdown here..."
            spellCheck={false}
          />
          <div className="mt-2 text-xs text-gray-500">
            Tip: You can paste HTML or Markdown. Click "Visual" to see the rendered result.
          </div>
        </div>
      ) : (
        <div className="bg-white">
          <EditorContent editor={editor} />
        </div>
      )}
    </div>
  )
}
