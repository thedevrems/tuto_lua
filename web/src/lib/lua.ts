import { LuaFactory, type LuaEngine } from 'wasmoon'
// Let Vite resolve & serve the wasm glue so it loads reliably in the browser
// (avoids wasmoon's Node-style require/createRequire resolution).
import glueWasmUrl from 'wasmoon/dist/glue.wasm?url'
import type { LuaTest } from '../types'

export interface RunResult {
  /** Captured print() output, lines joined with "\n". */
  output: string
  /** Runtime/compile error message, or null when the code ran cleanly. */
  error: string | null
}

export interface TestOutcome {
  name: string
  passed: boolean
  /** Failure message (assert message or error), null on success. */
  message: string | null
}

export interface TestRunResult extends RunResult {
  results: TestOutcome[]
}

const MAX_OUTPUT_LINES = 5000
// Abort the VM after this many instructions to avoid freezing on infinite loops.
const INSTRUCTION_LIMIT = 8_000_000

let factoryPromise: Promise<LuaFactory> | null = null

function getFactory(): Promise<LuaFactory> {
  if (!factoryPromise) {
    factoryPromise = Promise.resolve(new LuaFactory(glueWasmUrl))
  }
  return factoryPromise
}

/**
 * Build a fresh Lua engine with a captured `print`, output helpers and an
 * instruction-count guard. Returns the engine plus the live output buffer.
 */
async function createSandbox(): Promise<{ lua: LuaEngine; lines: string[] }> {
  const factory = await getFactory()
  const lua = await factory.createEngine({ openStandardLibs: true })
  const lines: string[] = []

  lua.global.set('__capture', (s: string) => {
    if (lines.length < MAX_OUTPUT_LINES) {
      lines.push(s)
    } else if (lines.length === MAX_OUTPUT_LINES) {
      lines.push('… (sortie tronquée)')
    }
  })

  await lua.doString(`
    _printed = {}

    local function __tostr(...)
      local parts = {}
      local n = select('#', ...)
      for i = 1, n do parts[i] = tostring(select(i, ...)) end
      return table.concat(parts, "\\t")
    end

    function print(...)
      local s = __tostr(...)
      _printed[#_printed + 1] = s
      __capture(s)
    end

    -- Test helpers
    function output_lines() return _printed end
    function output_text() return table.concat(_printed, "\\n") end
    function printed(sub)
      for _, line in ipairs(_printed) do
        if string.find(line, sub, 1, true) then return true end
      end
      return false
    end

    -- Instruction guard against infinite loops
    debug.sethook(function()
      error("Temps d'execution depasse : boucle infinie ou code trop long ?", 2)
    end, "", ${INSTRUCTION_LIMIT})
  `)

  return { lua, lines }
}

function cleanError(message: unknown): string {
  let raw = typeof message === 'string' ? message : String((message as Error)?.message ?? message)
  // Drop the Lua stack traceback — not useful for beginners.
  raw = raw.split('\nstack traceback:')[0]
  // Replace the "[string \"...\"]:N:" chunk prefix (content may contain quotes)
  // with a friendlier "ligne N :".
  raw = raw.replace(/\[string "[\s\S]*?"\]:(\d+):/g, 'ligne $1 :')
  // wasmoon prefixes engine errors with "Lua Error(...)"; strip it.
  raw = raw.replace(/^Lua Error\([^)]*\):?\s*/, '')
  return raw.trim()
}

/** Test assertion messages live on a single synthetic line; drop that prefix. */
function cleanTestMessage(message: unknown): string {
  return cleanError(message).replace(/^ligne \d+ :\s*/, '')
}

/** Run the student's code and capture its output. */
export async function runLua(code: string): Promise<RunResult> {
  let sandbox: { lua: LuaEngine; lines: string[] } | null = null
  try {
    sandbox = await createSandbox()
  } catch (e) {
    return { output: '', error: `Impossible de démarrer Lua : ${cleanError(e)}` }
  }
  const { lua, lines } = sandbox
  let error: string | null = null
  try {
    await lua.doString(code)
  } catch (e) {
    error = cleanError(e)
  } finally {
    lua.global.close()
  }
  return { output: lines.join('\n'), error }
}

/** Run the student's code, then each test in the same state. */
export async function runLuaTests(code: string, tests: LuaTest[]): Promise<TestRunResult> {
  let sandbox: { lua: LuaEngine; lines: string[] } | null = null
  try {
    sandbox = await createSandbox()
  } catch (e) {
    return { output: '', error: `Impossible de démarrer Lua : ${cleanError(e)}`, results: [] }
  }
  const { lua, lines } = sandbox
  const results: TestOutcome[] = []
  let error: string | null = null

  try {
    await lua.doString(code)
  } catch (e) {
    error = cleanError(e)
  }

  if (!error) {
    for (const test of tests) {
      try {
        await lua.doString(test.code)
        results.push({ name: test.name, passed: true, message: null })
      } catch (e) {
        results.push({ name: test.name, passed: false, message: cleanTestMessage(e) })
      }
    }
  }

  lua.global.close()
  return { output: lines.join('\n'), error, results }
}
