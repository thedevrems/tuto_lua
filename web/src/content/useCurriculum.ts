import { useEffect, useState } from 'react'
import { api } from '../lib/api'
import type { Module } from '../types'
import { coursesToModules } from './fromApi'

interface CurriculumState {
  modules: Module[] | null
  loading: boolean
  error: string | null
}

/** Loads the full course catalogue from the API and maps it to Module[]. */
export function useCurriculum(): CurriculumState {
  const [modules, setModules] = useState<Module[] | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const list = await api.courses.list()
        const trees = await Promise.all(list.map((c) => api.courses.tree(c.slug)))
        if (!cancelled) setModules(coursesToModules(trees))
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Erreur de chargement des cours')
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

  return { modules, loading: modules === null && error === null, error }
}
