import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { curriculum } from '../content/curriculum'
import type { Item } from '../types'
import { useAuth } from '../auth/AuthContext'
import Sidebar from '../components/Sidebar'
import LessonContent from '../components/LessonContent'
import ExercisePanel from '../components/ExercisePanel'

const LS_CODE = 'lua-academy:code'
const LS_DONE = 'lua-academy:completed'
const LS_LAST = 'lua-academy:last'

// Flattened, ordered list of all reachable items (for lookup + prev/next).
const allItems: Item[] = curriculum.flatMap((m) => m.chapters.flatMap((c) => c.items))

function loadJSON<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key)
    return raw ? (JSON.parse(raw) as T) : fallback
  } catch {
    return fallback
  }
}

export default function LearnPage() {
  const { user, logout } = useAuth()
  const [codeMap, setCodeMap] = useState<Record<string, string>>(() => loadJSON(LS_CODE, {}))
  const [completed, setCompleted] = useState<Set<string>>(() => new Set(loadJSON<string[]>(LS_DONE, [])))
  const [activeId, setActiveId] = useState<string>(() => {
    const last = localStorage.getItem(LS_LAST)
    if (last && allItems.some((i) => i.id === last)) return last
    return allItems[0]?.id ?? ''
  })
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const active = useMemo(() => allItems.find((i) => i.id === activeId) ?? allItems[0], [activeId])

  const index = allItems.findIndex((i) => i.id === active?.id)
  const prev = index > 0 ? allItems[index - 1] : null
  const next = index >= 0 && index < allItems.length - 1 ? allItems[index + 1] : null

  useEffect(() => {
    if (active) localStorage.setItem(LS_LAST, active.id)
  }, [active])

  const handleSelect = (item: Item) => {
    setActiveId(item.id)
    setSidebarOpen(false)
  }

  const handleCodeChange = (id: string, code: string) => {
    setCodeMap((prev) => {
      const next = { ...prev, [id]: code }
      localStorage.setItem(LS_CODE, JSON.stringify(next))
      return next
    })
  }

  const handleSolved = (id: string) => {
    setCompleted((prev) => {
      if (prev.has(id)) return prev
      const next = new Set(prev)
      next.add(id)
      localStorage.setItem(LS_DONE, JSON.stringify([...next]))
      return next
    })
  }

  const exerciseCount = allItems.filter((i) => i.kind === 'exercise').length
  const doneCount = [...completed].filter((id) => allItems.some((i) => i.id === id && i.kind === 'exercise')).length

  return (
    <div className="h-screen flex flex-col bg-white">
      {/* Top bar */}
      <header className="flex items-center gap-4 px-4 h-14 border-b border-gray-200 shrink-0">
        <button
          onClick={() => setSidebarOpen((o) => !o)}
          className="lg:hidden text-gray-600 hover:text-black px-2 py-1"
          aria-label="Menu"
        >
          ☰
        </button>
        <Link to="/" className="flex items-center gap-2.5">
          <span className="w-7 h-7 rounded-md bg-black text-white font-bold grid place-items-center font-mono text-sm">
            L
          </span>
          <div className="leading-tight">
            <div className="text-[15px] font-bold text-black tracking-tight">Lua Academy</div>
            <div className="text-[10px] text-gray-500 uppercase tracking-widest">Lua pour FiveM</div>
          </div>
        </Link>

        <div className="ml-auto flex items-center gap-3">
          <span className="text-[12px] text-gray-600 font-mono px-2.5 py-1 rounded-md border border-gray-200">
            {doneCount}/{exerciseCount} exercices
          </span>
          {user ? (
            <div className="flex items-center gap-2">
              <span className="hidden sm:inline text-[12px] text-gray-600">{user.username}</span>
              <button onClick={logout} className="btn btn-ghost btn-sm">
                Déconnexion
              </button>
            </div>
          ) : (
            <Link to="/login" className="btn btn-secondary btn-sm">
              Connexion
            </Link>
          )}
        </div>
      </header>

      <div className="flex-1 flex min-h-0">
        {/* Sidebar (overlay on mobile, fixed on desktop) */}
        <aside
          className={
            'w-[300px] shrink-0 border-r border-gray-200 bg-white z-20 ' +
            'lg:static lg:block ' +
            (sidebarOpen ? 'fixed inset-y-0 left-0 top-14 block' : 'hidden')
          }
        >
          <Sidebar modules={curriculum} activeId={active?.id ?? ''} completed={completed} onSelect={handleSelect} />
        </aside>

        {sidebarOpen && (
          <div className="lg:hidden fixed inset-0 top-14 bg-black/40 z-10" onClick={() => setSidebarOpen(false)} />
        )}

        {/* Main content */}
        <main className="flex-1 flex flex-col min-h-0">
          <div className="flex-1 min-h-0">
            {!active ? (
              <div className="grid place-items-center h-full text-gray-500">Sélectionnez un cours.</div>
            ) : active.kind === 'lesson' ? (
              <LessonContent lesson={active} />
            ) : (
              <ExercisePanel
                exercise={active}
                savedCode={codeMap[active.id]}
                onCodeChange={handleCodeChange}
                onSolved={handleSolved}
              />
            )}
          </div>

          {/* Prev / Next */}
          <footer className="flex items-center justify-between gap-2 px-4 py-2.5 border-t border-gray-200 shrink-0">
            <button
              onClick={() => prev && handleSelect(prev)}
              disabled={!prev}
              className="text-[12px] px-3 py-1.5 rounded-md text-gray-600 hover:bg-gray-100 hover:text-black disabled:opacity-30 disabled:hover:bg-transparent transition-colors max-w-[45%] truncate"
            >
              {prev ? `← ${prev.title}` : ''}
            </button>
            <button
              onClick={() => next && handleSelect(next)}
              disabled={!next}
              className="text-[12px] px-3 py-1.5 rounded-md text-gray-600 hover:bg-gray-100 hover:text-black disabled:opacity-30 disabled:hover:bg-transparent transition-colors max-w-[45%] truncate"
            >
              {next ? `${next.title} →` : ''}
            </button>
          </footer>
        </main>
      </div>
    </div>
  )
}
