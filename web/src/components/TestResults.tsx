import type { TestOutcome } from '../lib/lua'

interface Props {
  results: TestOutcome[]
  /** A runtime error that prevented the tests from running. */
  runError: string | null
}

export default function TestResults({ results, runError }: Props) {
  if (runError) {
    return (
      <div className="px-3 py-3 text-[12.5px] font-mono text-danger border-l-2 border-danger">
        ✗ Le code n'a pas pu s'exécuter — {runError}
      </div>
    )
  }

  if (results.length === 0) return null

  const passed = results.filter((r) => r.passed).length
  const allPassed = passed === results.length

  return (
    <div className="flex flex-col h-full min-h-0">
      <div className="flex items-center justify-between px-3 py-2 border-b border-gray-200">
        <span className="text-[11px] uppercase tracking-widest text-gray-500">Tests</span>
        <span
          className={
            'text-[11px] font-mono px-2 py-0.5 rounded border ' +
            (allPassed ? 'border-success-border bg-success-bg text-success' : 'border-gray-300 text-gray-600')
          }
        >
          {passed}/{results.length} réussis
        </span>
      </div>

      <div className="flex-1 min-h-0 overflow-auto p-2 space-y-1">
        {results.map((r, i) => (
          <div
            key={i}
            className="flex items-start gap-2 px-2 py-1.5 rounded text-[12.5px] font-mono bg-gray-100"
          >
            <span className={'mt-px ' + (r.passed ? 'text-success' : 'text-danger')}>
              {r.passed ? '✓' : '✗'}
            </span>
            <span className="flex-1">
              <span className={r.passed ? 'text-gray-700' : 'text-black font-medium'}>{r.name}</span>
              {!r.passed && r.message && (
                <span className="block text-gray-500 mt-0.5">{r.message}</span>
              )}
            </span>
          </div>
        ))}

        {allPassed && (
          <div className="px-2 py-2 mt-1 text-[12.5px] text-success font-medium">
            Tous les tests passent. Bien joué !
          </div>
        )}
      </div>
    </div>
  )
}
