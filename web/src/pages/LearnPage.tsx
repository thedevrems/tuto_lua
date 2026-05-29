import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import type { Item } from '../types'
import { useAuth } from '../auth/AuthContext'
import { useCurriculum, type LockedCourse } from '../content/useCurriculum'
import { useProgress } from '../content/useProgress'
import Sidebar from '../components/Sidebar'
import LessonContent from '../components/LessonContent'
import ExercisePanel from '../components/ExercisePanel'
import Paywall from '../components/Paywall'

const LS_LAST = 'lua-academy:last'

export default function LearnPage() {
  const { user, logout } = useAuth()
  const { modules, locked, loading, error } = useCurriculum(user)
  const { codeMap, completed, saveCode, markSolved } = useProgress(user?.id ?? null)

  const [activeId, setActiveId] = useState<string>(() => localStorage.getItem(LS_LAST) ?? '')
  const [lockedView, setLockedView] = useState<LockedCourse | null>(null)
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const allItems = useMemo<Item[]>(
    () => (modules ?? []).flatMap((m) => m.chapters.flatMap((c) => c.items)),
    [modules],
  )

  useEffect(() => {
    if (allItems.length > 0 && !allItems.some((i) => i.id === activeId)) {
      setActiveId(allItems[0].id)
    }
  }, [allItems, activeId])

  const active = useMemo(() => allItems.find((i) => i.id === activeId) ?? null, [allItems, activeId])

  useEffect(() => {
    if (active) localStorage.setItem(LS_LAST, active.id)
  }, [active])

  const index = allItems.findIndex((i) => i.id === active?.id)
  const prev = index > 0 ? allItems[index - 1] : null
  const next = index >= 0 && index < allItems.length - 1 ? allItems[index + 1] : null

  const handleSelect = (item: Item) => {
    setActiveId(item.id)
    setLockedView(null)
    setSidebarOpen(false)
  }

  const handleSelectLocked = (course: LockedCourse) => {
    setLockedView(course)
    setSidebarOpen(false)
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
            'w-[300px] shrink-0 border-r border-gray-200 bg-white z-20 overflow-y-auto ' +
            'lg:static lg:block ' +
            (sidebarOpen ? 'fixed inset-y-0 left-0 top-14 block' : 'hidden')
          }
        >
          <Sidebar modules={modules ?? []} activeId={lockedView ? '' : active?.id ?? ''} completed={completed} onSelect={handleSelect} />
          <LockedList courses={locked} activeSlug={lockedView?.slug ?? ''} onSelect={handleSelectLocked} />
        </aside>

        {sidebarOpen && (
          <div className="lg:hidden fixed inset-0 top-14 bg-black/40 z-10" onClick={() => setSidebarOpen(false)} />
        )}

        {/* Main content */}
        <main className="flex-1 flex flex-col min-h-0">
          <div className="flex-1 min-h-0">
            {lockedView ? (
              <Paywall course={lockedView} />
            ) : loading ? (
              <div className="grid place-items-center h-full text-gray-500">Chargement des cours…</div>
            ) : error ? (
              <BackendError message={error} />
            ) : !active ? (
              <div className="grid place-items-center h-full text-gray-500">Sélectionnez un cours.</div>
            ) : active.kind === 'lesson' ? (
              <LessonContent lesson={active} />
            ) : (
              <ExercisePanel
                exercise={active}
                savedCode={codeMap[active.id]}
                onCodeChange={saveCode}
                onSolved={markSolved}
              />
            )}
          </div>

          {/* Prev / Next */}
          <footer className="flex items-center justify-between gap-2 px-4 py-2.5 border-t border-gray-200 shrink-0">
            <button
              onClick={() => prev && handleSelect(prev)}
              disabled={!prev || !!lockedView}
              className="text-[12px] px-3 py-1.5 rounded-md text-gray-600 hover:bg-gray-100 hover:text-black disabled:opacity-30 disabled:hover:bg-transparent transition-colors max-w-[45%] truncate"
            >
              {prev ? `← ${prev.title}` : ''}
            </button>
            <button
              onClick={() => next && handleSelect(next)}
              disabled={!next || !!lockedView}
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

/** Lists paid courses the user has not unlocked, below the main navigation. */
function LockedList({
  courses,
  activeSlug,
  onSelect,
}: {
  courses: LockedCourse[]
  activeSlug: string
  onSelect: (c: LockedCourse) => void
}) {
  if (courses.length === 0) return null
  return (
    <div className="px-3 pb-6">
      <div className="px-2 pt-4 pb-1 text-[11px] font-medium uppercase tracking-wide text-gray-500">
        Cours à débloquer
      </div>
      {courses.map((c) => (
        <button
          key={c.slug}
          onClick={() => onSelect(c)}
          className={
            'w-full flex items-center gap-2 text-left px-2 py-1.5 rounded-md text-[13px] transition-colors ' +
            (c.slug === activeSlug ? 'bg-black text-white font-medium' : 'text-gray-600 hover:bg-gray-100 hover:text-black')
          }
        >
          <span className="shrink-0">🔒</span>
          <span className="flex-1 leading-tight">{c.title}</span>
        </button>
      ))}
    </div>
  )
}

/** Shown when the course API is unreachable (e.g. backend not started). */
function BackendError({ message }: { message: string }) {
  return (
    <div className="grid h-full place-items-center px-6 text-center">
      <div className="max-w-md">
        <h2 className="text-lg font-semibold text-black">Impossible de charger les cours</h2>
        <p className="mt-2 text-sm text-gray-600">{message}</p>
        <p className="mt-4 text-xs text-gray-500">
          Vérifiez que le backend tourne (<code className="font-mono">cd server &amp;&amp; go run ./cmd/api</code>).
        </p>
      </div>
    </div>
  )
}
