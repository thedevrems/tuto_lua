import { useMemo } from 'react'
import CodeMirror from '@uiw/react-codemirror'
import { EditorView } from '@codemirror/view'
import { StreamLanguage } from '@codemirror/language'
import { lua } from '@codemirror/legacy-modes/mode/lua'

interface Props {
  value: string
  onChange: (value: string) => void
  readOnly?: boolean
  minHeight?: string
}

// Light monochrome editor theme (charte graphique: black & white).
const monoTheme = EditorView.theme(
  {
    '&': { color: '#171717', backgroundColor: 'transparent', height: '100%' },
    '.cm-content': { caretColor: '#0a0a0a', padding: '12px 0' },
    '.cm-cursor, .cm-dropCursor': { borderLeftColor: '#0a0a0a' },
    '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': {
      backgroundColor: '#e5e5e5',
    },
    '.cm-gutters': { backgroundColor: 'transparent', color: '#a3a3a3', border: 'none' },
    '.cm-activeLine': { backgroundColor: 'rgba(0,0,0,0.025)' },
    '.cm-activeLineGutter': { backgroundColor: 'transparent', color: '#525252' },
    '.cm-lineNumbers .cm-gutterElement': { padding: '0 12px 0 8px' },
    '.cm-selectionMatch': { backgroundColor: '#f5f5f5' },
    '.cm-matchingBracket, &.cm-focused .cm-matchingBracket': {
      backgroundColor: '#e5e5e5',
      outline: 'none',
    },
  },
  { dark: false },
)

// Grayscale syntax highlighting for Lua tokens (charte: monochrome).
const monoHighlight = EditorView.theme({
  '.tok-keyword': { color: '#0a0a0a', fontWeight: '600' },
  '.tok-operator': { color: '#525252' },
  '.tok-string': { color: '#525252' },
  '.tok-number': { color: '#404040' },
  '.tok-comment': { color: '#a3a3a3', fontStyle: 'italic' },
  '.tok-variableName': { color: '#171717' },
})

export default function CodeEditor({ value, onChange, readOnly = false, minHeight = '180px' }: Props) {
  const extensions = useMemo(
    () => [StreamLanguage.define(lua), monoTheme, monoHighlight, EditorView.lineWrapping],
    [],
  )

  return (
    <CodeMirror
      value={value}
      onChange={onChange}
      extensions={extensions}
      readOnly={readOnly}
      basicSetup={{
        lineNumbers: true,
        foldGutter: false,
        highlightActiveLine: !readOnly,
        highlightActiveLineGutter: !readOnly,
        autocompletion: false,
        bracketMatching: true,
        closeBrackets: true,
        indentOnInput: true,
      }}
      style={{ minHeight, fontSize: '13px' }}
      theme="none"
    />
  )
}
