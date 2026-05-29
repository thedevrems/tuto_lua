import { useEffect, useRef, useState } from 'react'
import { api } from '../lib/api'

const LS_CODE = 'lua-academy:code'
const LS_DONE = 'lua-academy:completed'
const PUSH_DELAY = 1200 // ms of inactivity before syncing code to the server

function loadJSON<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key)
    return raw ? (JSON.parse(raw) as T) : fallback
  } catch {
    return fallback
  }
}

interface ProgressApi {
  codeMap: Record<string, string>
  completed: Set<string>
  saveCode: (exerciseId: string, code: string) => void
  markSolved: (exerciseId: string) => void
}

/** Tracks per-exercise code and completion. localStorage is the offline cache;
 *  when a user is logged in, state is loaded from and synced to the backend. */
export function useProgress(userId: string | null): ProgressApi {
  const [codeMap, setCodeMap] = useState<Record<string, string>>(() => loadJSON(LS_CODE, {}))
  const [completed, setCompleted] = useState<Set<string>>(() => new Set(loadJSON<string[]>(LS_DONE, [])))

  const codeRef = useRef(codeMap)
  const doneRef = useRef(completed)
  codeRef.current = codeMap
  doneRef.current = completed
  const timers = useRef<Record<string, ReturnType<typeof setTimeout>>>({})

  // On login, merge the server's saved progress into the local state.
  useEffect(() => {
    if (!userId) return
    let cancelled = false
    api.progress
      .list()
      .then((rows) => {
        if (cancelled) return
        setCodeMap((prev) => {
          const next = { ...prev }
          rows.forEach((r) => (next[r.exerciseId] = r.code))
          localStorage.setItem(LS_CODE, JSON.stringify(next))
          return next
        })
        setCompleted((prev) => {
          const next = new Set(prev)
          rows.forEach((r) => r.completed && next.add(r.exerciseId))
          localStorage.setItem(LS_DONE, JSON.stringify([...next]))
          return next
        })
      })
      .catch(() => {})
    return () => {
      cancelled = true
    }
  }, [userId])

  const saveCode = (exerciseId: string, code: string) => {
    setCodeMap((prev) => {
      const next = { ...prev, [exerciseId]: code }
      localStorage.setItem(LS_CODE, JSON.stringify(next))
      return next
    })
    if (userId) schedulePush(timers.current, exerciseId, () => api.progress.save(exerciseId, code, doneRef.current.has(exerciseId)))
  }

  const markSolved = (exerciseId: string) => {
    setCompleted((prev) => {
      if (prev.has(exerciseId)) return prev
      const next = new Set(prev)
      next.add(exerciseId)
      localStorage.setItem(LS_DONE, JSON.stringify([...next]))
      return next
    })
    if (userId) {
      api.progress.save(exerciseId, codeRef.current[exerciseId] ?? '', true).catch(() => {})
    }
  }

  return { codeMap, completed, saveCode, markSolved }
}

// schedulePush debounces a server sync per exercise so we don't POST on every keystroke.
function schedulePush(
  timers: Record<string, ReturnType<typeof setTimeout>>,
  key: string,
  push: () => Promise<unknown>,
) {
  if (timers[key]) clearTimeout(timers[key])
  timers[key] = setTimeout(() => {
    push().catch(() => {})
  }, PUSH_DELAY)
}
