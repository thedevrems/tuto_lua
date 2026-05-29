export type Difficulty = 'facile' | 'moyen' | 'difficile'

/** A single automated test run against the student's code. */
export interface LuaTest {
  /** Short label shown in the results list. */
  name: string
  /**
   * Lua chunk executed in the SAME state, right after the student's code.
   * It should `assert(condition, "message")`. Throwing/asserting => fail.
   *
   * Helpers injected into the state and usable here:
   *   - `printed(substr)`  -> true if any printed line contains substr
   *   - `output_text()`    -> all printed lines joined with "\n"
   *   - `output_lines()`   -> table of printed lines
   * Globals defined by the student (PascalCase functions, globals) are visible.
   */
  code: string
}

/** A lesson: pure reading material (markdown). */
export interface Lesson {
  kind: 'lesson'
  id: string
  title: string
  /** Markdown body. */
  content: string
}

/** An interactive exercise with editor, console, tests and a solution. */
export interface Exercise {
  kind: 'exercise'
  id: string
  title: string
  difficulty: Difficulty
  /** Markdown statement / subject shown above the editor. */
  statement: string
  /** Code pre-filled in the editor. */
  starter: string
  /** Reference solution revealed on demand. */
  solution: string
  /** Optional automated tests. */
  tests?: LuaTest[]
  /** Optional hints revealed one by one. */
  hints?: string[]
}

export type Item = Lesson | Exercise

export interface Chapter {
  id: string
  title: string
  /** Short blurb shown in the chapter header. */
  summary?: string
  items: Item[]
  /** When true the chapter is a placeholder (curriculum announced, content WIP). */
  comingSoon?: boolean
}

export interface Module {
  id: string
  title: string
  summary?: string
  chapters: Chapter[]
}
