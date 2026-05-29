interface Props {
  output: string
  error: string | null
  running: boolean
  /** Has the user run anything yet? */
  ran: boolean
}

export default function Console({ output, error, running, ran }: Props) {
  return (
    <div className="flex flex-col h-full min-h-0">
      <div className="flex items-center gap-2 px-3 py-2 border-b border-ink-800 text-[11px] uppercase tracking-widest text-ink-400">
        <span className="flex gap-1.5">
          <span className="w-2.5 h-2.5 rounded-full bg-ink-700" />
          <span className="w-2.5 h-2.5 rounded-full bg-ink-700" />
          <span className="w-2.5 h-2.5 rounded-full bg-ink-700" />
        </span>
        Console
      </div>
      <div className="flex-1 min-h-0 overflow-auto p-3 font-mono text-[12.5px] leading-relaxed">
        {running && <span className="text-ink-400">Exécution…</span>}

        {!running && !ran && (
          <span className="text-ink-500">
            La sortie de <code className="text-ink-300">print()</code> s'affichera ici. Cliquez sur «&nbsp;Lancer&nbsp;».
          </span>
        )}

        {!running && ran && (
          <>
            {output && <pre className="whitespace-pre-wrap text-ink-100 m-0">{output}</pre>}
            {error && (
              <pre className="whitespace-pre-wrap text-white m-0 mt-2 border-l-2 border-white pl-3">
                ✗ Erreur — {error}
              </pre>
            )}
            {!output && !error && <span className="text-ink-500">(aucune sortie)</span>}
          </>
        )}
      </div>
    </div>
  )
}
