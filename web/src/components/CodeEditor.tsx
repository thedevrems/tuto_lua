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

// Monochrome editor theme (black & white).
const monoTheme = EditorView.theme(
  {
    '&': { color: '#ededf0', backgroundColor: 'transparent', height: '100%' },
    '.cm-content': { caretColor: '#ffffff', padding: '12px 0' },
    '.cm-cursor, .cm-dropCursor': { borderLeftColor: '#ffffff' },
    '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection': {
      backgroundColor: '#33333b',
    },
    '.cm-gutters': { backgroundColor: 'transparent', color: '#5b5b66', border: 'none' },
    '.cm-activeLine': { backgroundColor: 'rgba(255,255,255,0.025)' },
    '.cm-activeLineGutter': { backgroundColor: 'transparent', color: '#b4b4bd' },
    '.cm-lineNumbers .cm-gutterElement': { padding: '0 12px 0 8px' },
    '.cm-selectionMatch': { backgroundColor: '#2a2a31' },
    '.cm-matchingBracket, &.cm-focused .cm-matchingBracket': {
      backgroundColor: '#3a3a42',
      outline: 'none',
    },
  },
  { dark: true },
)

// Grayscale syntax highlighting for Lua tokens.
const monoHighlight = EditorView.theme({
  '.tok-keyword': { color: '#ffffff', fontWeight: '600' },
  '.tok-operator': { color: '#b4b4bd' },
  '.tok-string': { color: '#b4b4bd' },
  '.tok-number': { color: '#d6d6dc' },
  '.tok-comment': { color: '#5b5b66', fontStyle: 'italic' },
  '.tok-variableName': { color: '#ededf0' },
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
