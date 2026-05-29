import { useEffect, useMemo, useState } from 'react'
import type { Exercise } from '../types'
import { runLua, runLuaTests, type TestOutcome } from '../lib/lua'
import CodeEditor from './CodeEditor'
import Console from './Console'
import TestResults from './TestResults'
import Markdown from './Markdown'

interface Props {
  exercise: Exercise
  savedCode: string | undefined
  onCodeChange: (id: string, code: string) => void
  onSolved: (id: string) => void
}

const difficultyLabel: Record<Exercise['difficulty'], string> = {
  facile: 'Facile',
  moyen: 'Moyen',
  difficile: 'Difficile',
}

type Tab = 'console' | 'tests'

export default function ExercisePanel({ exercise, savedCode, onCodeChange, onSolved }: Props) {
  const [code, setCode] = useState(savedCode ?? exercise.starter)
  const [tab, setTab] = useState<Tab>('console')
  const [running, setRunning] = useState(false)
  const [ran, setRan] = useState(false)

  const [output, setOutput] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [testResults, setTestResults] = useState<TestOutcome[]>([])
  const [testError, setTestError] = useState<string | null>(null)

  const [showSolution, setShowSolution] = useState(false)
  const [hintsShown, setHintsShown] = useState(0)

  const hasTests = !!exercise.tests && exercise.tests.length > 0

  // Reset everything when switching exercise.
  useEffect(() => {
    setCode(savedCode ?? exercise.starter)
    setTab('console')
    setRunning(false)
    setRan(false)
    setOutput('')
    setError(null)
    setTestResults([])
    setTestError(null)
    setShowSolution(false)
    setHintsShown(0)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [exercise.id])

  const handleChange = (next: string) => {
    setCode(next)
    onCodeChange(exercise.id, next)
  }

  const handleRun = async () => {
    setTab('console')
    setRunning(true)
    const result = await runLua(code)
    setOutput(result.output)
    setError(result.error)
    setRunning(false)
    setRan(true)
  }

  const handleTest = async () => {
    if (!exercise.tests) return
    setTab('tests')
    setRunning(true)
    const result = await runLuaTests(code, exercise.tests)
    setOutput(result.output)
    setError(result.error)
    setTestResults(result.results)
    setTestError(result.error)
    setRunning(false)
    setRan(true)
    if (!result.error && result.results.length > 0 && result.results.every((r) => r.passed)) {
      onSolved(exercise.id)
    }
  }

  const handleReset = () => {
    handleChange(exercise.starter)
    setRan(false)
    setOutput('')
    setError(null)
    setTestResults([])
    setTestError(null)
  }

  const passedCount = useMemo(() => testResults.filter((r) => r.passed).length, [testResults])

  return (
    <div className="flex flex-col lg:flex-row h-full min-h-0">
      {/* ---- Subject / statement ---- */}
      <section className="lg:w-[42%] lg:max-w-2xl shrink-0 overflow-y-auto border-b lg:border-b-0 lg:border-r border-gray-200 px-6 py-6">
        <div className="flex items-center gap-3 mb-4">
          <span className="badge uppercase tracking-widest">{difficultyLabel[exercise.difficulty]}</span>
        </div>
        <h1 className="text-2xl font-bold text-black tracking-tight mb-4">{exercise.title}</h1>

        <Markdown>{exercise.statement}</Markdown>

        {/* Hints */}
        {exercise.hints && exercise.hints.length > 0 && (
          <div className="mt-6">
            {exercise.hints.slice(0, hintsShown).map((hint, i) => (
              <div key={i} className="mb-2 px-3 py-2 rounded-md bg-gray-100 border border-gray-200 text-[13px] text-gray-700">
                <span className="text-gray-500">Indice {i + 1} — </span>
                {hint}
              </div>
            ))}
            {hintsShown < exercise.hints.length && (
              <button
                onClick={() => setHintsShown((n) => n + 1)}
                className="text-[12px] text-gray-500 hover:text-black underline underline-offset-2 decoration-gray-300"
              >
                Afficher un indice ({hintsShown}/{exercise.hints.length})
              </button>
            )}
          </div>
        )}

        {/* Solution */}
        <div className="mt-6 pt-5 border-t border-gray-200">
          <button
            onClick={() => setShowSolution((s) => !s)}
            className="text-[13px] text-gray-600 hover:text-black inline-flex items-center gap-2"
          >
            <span className={'text-xs transition-transform ' + (showSolution ? 'rotate-90' : '')}>▶</span>
            {showSolution ? 'Masquer la réponse' : 'Afficher la réponse'}
          </button>
          {showSolution && (
            <div className="mt-3 rounded-lg border border-gray-200 overflow-hidden">
              <div className="px-3 py-1.5 text-[11px] uppercase tracking-widest text-gray-500 border-b border-gray-200 bg-gray-100 flex items-center justify-between">
                Solution
                <button
                  onClick={() => handleChange(exercise.solution)}
                  className="text-[11px] text-gray-600 hover:text-black normal-case tracking-normal"
                >
                  Charger dans l'éditeur
                </button>
              </div>
              <div className="bg-white">
                <CodeEditor value={exercise.solution} onChange={() => {}} readOnly minHeight="0" />
              </div>
            </div>
          )}
        </div>
      </section>

      {/* ---- Workspace : editor + output ---- */}
      <section className="flex-1 flex flex-col min-h-0">
        {/* Toolbar */}
        <div className="flex items-center gap-2 px-4 py-2.5 border-b border-gray-200">
          <span className="text-[11px] uppercase tracking-widest text-gray-500 mr-auto font-mono">solution.lua</span>
          <button
            onClick={handleReset}
            className="text-[12px] px-3 py-1.5 rounded-md text-gray-600 hover:bg-gray-100 hover:text-black transition-colors"
          >
            Réinitialiser
          </button>
          {hasTests && (
            <button
              onClick={handleTest}
              disabled={running}
              className="text-[12px] px-3 py-1.5 rounded-md border-2 border-black text-black hover:bg-black hover:text-white disabled:opacity-50 transition-all"
            >
              Tester
            </button>
          )}
          <button
            onClick={handleRun}
            disabled={running}
            className="text-[12px] px-4 py-1.5 rounded-md bg-black text-white font-semibold hover:bg-black-muted disabled:opacity-50 transition-colors"
          >
            {running ? '…' : '▶ Lancer'}
          </button>
        </div>

        {/* Editor */}
        <div className="flex-1 min-h-0 overflow-auto bg-white-soft">
          <CodeEditor value={code} onChange={handleChange} minHeight="100%" />
        </div>

        {/* Output : console / tests */}
        <div className="h-[40%] min-h-[180px] flex flex-col border-t border-gray-200 bg-white">
          <div className="flex items-center border-b border-gray-200">
            <button
              onClick={() => setTab('console')}
              className={
                'px-4 py-2 text-[12px] border-b-2 -mb-px transition-colors ' +
                (tab === 'console' ? 'border-black text-black' : 'border-transparent text-gray-500 hover:text-black')
              }
            >
              Console
            </button>
            {hasTests && (
              <button
                onClick={() => setTab('tests')}
                className={
                  'px-4 py-2 text-[12px] border-b-2 -mb-px transition-colors flex items-center gap-2 ' +
                  (tab === 'tests' ? 'border-black text-black' : 'border-transparent text-gray-500 hover:text-black')
                }
              >
                Tests
                {ran && testResults.length > 0 && (
                  <span className="text-[10px] font-mono text-gray-500">
                    {passedCount}/{testResults.length}
                  </span>
                )}
              </button>
            )}
          </div>

          <div className="flex-1 min-h-0">
            {tab === 'console' ? (
              <Console output={output} error={error} running={running} ran={ran} />
            ) : (
              <TestResults results={testResults} runError={testError} />
            )}
          </div>
        </div>
      </section>
    </div>
  )
}
