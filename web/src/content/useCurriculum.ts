import { useEffect, useState } from 'react'
import { api, type ApiCourse, type User } from '../lib/api'
import type { Module } from '../types'
import { coursesToModules } from './fromApi'

/** A paid course the current user cannot access yet (shown as a paywall). */
export interface LockedCourse {
  id: string
  slug: string
  title: string
  summary: string
  priceCents: number
  currency: string
}

interface CurriculumState {
  modules: Module[] | null
  locked: LockedCourse[]
  loading: boolean
  error: string | null
}

/** Loads the catalogue, splits it into accessible vs locked courses for the
 *  given user, and fetches full content only for the accessible ones. */
export function useCurriculum(user: User | null): CurriculumState {
  const [state, setState] = useState<CurriculumState>({ modules: null, locked: [], loading: true, error: null })

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const result = await loadCurriculum(user)
        if (!cancelled) setState({ ...result, loading: false, error: null })
      } catch (e) {
        if (!cancelled) {
          setState({ modules: null, locked: [], loading: false, error: errorMessage(e) })
        }
      }
    })()
    return () => {
      cancelled = true
    }
  }, [user?.id, user?.role])

  return state
}

async function loadCurriculum(user: User | null): Promise<{ modules: Module[]; locked: LockedCourse[] }> {
  const catalogue = await api.courses.list()
  const accessibleIds = await accessibleCourseIds(user)

  const isOpen = (c: ApiCourse) =>
    c.priceCents === 0 || user?.role === 'admin' || accessibleIds.has(c.id)

  const trees = await Promise.all(catalogue.filter(isOpen).map((c) => api.courses.tree(c.slug)))
  const locked = catalogue.filter((c) => !isOpen(c)).map(toLockedCourse)
  return { modules: coursesToModules(trees), locked }
}

// accessibleCourseIds returns the set of course ids the user is enrolled in.
async function accessibleCourseIds(user: User | null): Promise<Set<string>> {
  if (!user) return new Set()
  try {
    return new Set(await api.enrollments.mine())
  } catch {
    return new Set()
  }
}

function toLockedCourse(c: ApiCourse): LockedCourse {
  return { id: c.id, slug: c.slug, title: c.title, summary: c.summary, priceCents: c.priceCents, currency: c.currency }
}

function errorMessage(e: unknown): string {
  return e instanceof Error ? e.message : 'Erreur de chargement des cours'
}
