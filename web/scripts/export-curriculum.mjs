// One-off generator: bundles the TypeScript curriculum (with its `?raw` markdown
// imports) and emits server/seed/curriculum.json, the backend's seed data.
//
//   node scripts/export-curriculum.mjs
//
// Each frontend Module becomes a purchasable Course; Module 1 is free.
import { build } from 'esbuild'
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { pathToFileURL } from 'node:url'

// Resolve Vite's `import x from './file.md?raw'` to the file's text content.
const rawPlugin = {
  name: 'vite-raw',
  setup(b) {
    b.onResolve({ filter: /\?raw$/ }, (args) => ({
      path: resolve(args.resolveDir, args.path.replace(/\?raw$/, '')),
      namespace: 'raw',
    }))
    b.onLoad({ filter: /.*/, namespace: 'raw' }, (args) => ({
      contents: `export default ${JSON.stringify(readFileSync(args.path, 'utf8'))}`,
      loader: 'js',
    }))
  },
}

const bundle = resolve('node_modules/.cache/curriculum.bundle.mjs')
mkdirSync(dirname(bundle), { recursive: true })
await build({
  entryPoints: ['src/content/curriculum.ts'],
  bundle: true,
  format: 'esm',
  platform: 'node',
  outfile: bundle,
  plugins: [rawPlugin],
  logLevel: 'silent',
})

const { curriculum } = await import(pathToFileURL(bundle).href)

const courses = curriculum.map((module, mi) => ({
  slug: module.id,
  title: module.title,
  summary: module.summary ?? '',
  priceCents: mi === 0 ? 0 : 4900, // first module free, others 49 €
  currency: 'eur',
  published: true,
  position: mi,
  chapters: module.chapters.map((chapter, ci) => transformChapter(chapter, ci)),
}))

function transformChapter(chapter, ci) {
  const lessons = []
  const exercises = []
  chapter.items.forEach((item, position) => {
    if (item.kind === 'lesson') {
      lessons.push({ title: item.title, content: item.content, position })
    } else {
      exercises.push({
        title: item.title,
        difficulty: item.difficulty,
        statement: item.statement,
        starter: item.starter,
        solution: item.solution,
        hints: item.hints ?? [],
        position,
        tests: (item.tests ?? []).map((t, ti) => ({ name: t.name, code: t.code, position: ti })),
      })
    }
  })
  return { title: chapter.title, summary: chapter.summary ?? '', position: ci, lessons, exercises }
}

const out = resolve('../server/internal/seed/curriculum.json')
mkdirSync(dirname(out), { recursive: true })
writeFileSync(out, JSON.stringify(courses, null, 2) + '\n', 'utf8')
console.log(`Wrote ${courses.length} courses to ${out}`)
