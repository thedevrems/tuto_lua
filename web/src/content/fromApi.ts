// Maps the backend course tree (DB-backed) onto the frontend Module/Chapter/Item
// shape, so the learning UI is agnostic of where its content comes from.
import type { Module, Chapter, Item } from '../types'
import type { ApiCourse, ApiChapter } from '../lib/api'

const byPosition = <T extends { position: number }>(a: T, b: T) => a.position - b.position

/** Converts the API course list (each course = a module) into Module[]. */
export function coursesToModules(courses: ApiCourse[]): Module[] {
  return courses.slice().sort(byPosition).map(courseToModule)
}

function courseToModule(c: ApiCourse): Module {
  return {
    id: c.slug,
    title: c.title,
    summary: c.summary || undefined,
    chapters: (c.chapters ?? []).slice().sort(byPosition).map(chapterToChapter),
  }
}

function chapterToChapter(ch: ApiChapter): Chapter {
  const empty = (ch.lessons?.length ?? 0) === 0 && (ch.exercises?.length ?? 0) === 0
  return {
    id: ch.id,
    title: ch.title,
    summary: ch.summary || undefined,
    comingSoon: empty,
    items: mergeItems(ch),
  }
}

/** Interleaves lessons and exercises back into a single ordered item list. */
function mergeItems(ch: ApiChapter): Item[] {
  const positioned: Array<{ position: number; item: Item }> = []
  for (const l of ch.lessons ?? []) {
    positioned.push({ position: l.position, item: { kind: 'lesson', id: l.id, title: l.title, content: l.content } })
  }
  for (const e of ch.exercises ?? []) {
    positioned.push({
      position: e.position,
      item: {
        kind: 'exercise',
        id: e.id,
        title: e.title,
        difficulty: e.difficulty,
        statement: e.statement,
        starter: e.starter,
        solution: e.solution ?? '',
        tests: e.tests?.map((t) => ({ name: t.name, code: t.code })),
        hints: e.hints,
      },
    })
  }
  return positioned.sort(byPosition).map((p) => p.item)
}
