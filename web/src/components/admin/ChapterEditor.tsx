import { useState } from 'react'
import { api, type ApiChapter, type ApiExercise } from '../../lib/api'

/** Authoring UI for a single chapter: its lessons, exercises and their tests. */
export default function ChapterEditor({ chapter, onChange }: { chapter: ApiChapter; onChange: () => void }) {
  const lessons = chapter.lessons ?? []
  const exercises = chapter.exercises ?? []
  const nextPos = lessons.length + exercises.length

  return (
    <div className="card">
      <h4 className="font-semibold text-black">{chapter.title}</h4>

      <ul className="mt-3 space-y-1 text-sm text-gray-700">
        {lessons.map((l) => (
          <li key={l.id}>📄 {l.title}</li>
        ))}
        {exercises.map((e) => (
          <ExerciseBlock key={e.id} ex={e} onChange={onChange} />
        ))}
      </ul>

      <div className="mt-4 grid gap-3 sm:grid-cols-2">
        <AddLesson chapterId={chapter.id} position={nextPos} onDone={onChange} />
        <AddExercise chapterId={chapter.id} position={nextPos} onDone={onChange} />
      </div>
    </div>
  )
}

function ExerciseBlock({ ex, onChange }: { ex: ApiExercise; onChange: () => void }) {
  const [open, setOpen] = useState(false)
  return (
    <li className="border-l-2 border-gray-200 pl-3">
      <button onClick={() => setOpen((o) => !o)} className="text-left text-black hover:underline">
        ✏️ {ex.title} <span className="text-xs text-gray-400">({(ex.tests ?? []).length} test(s))</span>
      </button>
      {open && (
        <div className="mt-2">
          {(ex.tests ?? []).map((t) => (
            <div key={t.id} className="text-xs text-gray-500">• {t.name}</div>
          ))}
          <AddTest exerciseId={ex.id} position={(ex.tests ?? []).length} onDone={onChange} />
        </div>
      )}
    </li>
  )
}

function AddLesson({ chapterId, position, onDone }: { chapterId: string; position: number; onDone: () => void }) {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const submit = async () => {
    await api.admin.createLesson(chapterId, { title, content, position })
    setTitle('')
    setContent('')
    onDone()
  }
  return (
    <form onSubmit={runner(submit)} className="space-y-2 rounded-md border border-gray-200 p-3">
      <div className="text-xs font-semibold uppercase tracking-wide text-gray-500">Nouvelle leçon</div>
      <input className="input !py-1.5 text-sm" placeholder="Titre" value={title} onChange={(e) => setTitle(e.target.value)} required />
      <textarea className="input !py-1.5 text-sm" placeholder="Contenu (markdown)" rows={3} value={content} onChange={(e) => setContent(e.target.value)} />
      <button className="btn btn-secondary btn-sm w-full">Ajouter la leçon</button>
    </form>
  )
}

function AddExercise({ chapterId, position, onDone }: { chapterId: string; position: number; onDone: () => void }) {
  const [f, setF] = useState({ title: '', difficulty: 'facile', statement: '', starter: '', solution: '' })
  const set = (k: keyof typeof f) => (e: { target: { value: string } }) => setF((s) => ({ ...s, [k]: e.target.value }))
  const submit = async () => {
    await api.admin.createExercise(chapterId, { ...f, difficulty: f.difficulty as 'facile', hints: [], position })
    setF({ title: '', difficulty: 'facile', statement: '', starter: '', solution: '' })
    onDone()
  }
  return (
    <form onSubmit={runner(submit)} className="space-y-2 rounded-md border border-gray-200 p-3">
      <div className="text-xs font-semibold uppercase tracking-wide text-gray-500">Nouvel exercice</div>
      <input className="input !py-1.5 text-sm" placeholder="Titre" value={f.title} onChange={set('title')} required />
      <select className="input !py-1.5 text-sm" value={f.difficulty} onChange={set('difficulty')}>
        <option value="facile">Facile</option>
        <option value="moyen">Moyen</option>
        <option value="difficile">Difficile</option>
      </select>
      <textarea className="input !py-1.5 text-sm" placeholder="Énoncé (markdown)" rows={2} value={f.statement} onChange={set('statement')} />
      <textarea className="input !py-1.5 font-mono text-sm" placeholder="Code de départ" rows={2} value={f.starter} onChange={set('starter')} />
      <textarea className="input !py-1.5 font-mono text-sm" placeholder="Solution" rows={2} value={f.solution} onChange={set('solution')} />
      <button className="btn btn-secondary btn-sm w-full">Ajouter l'exercice</button>
    </form>
  )
}

function AddTest({ exerciseId, position, onDone }: { exerciseId: string; position: number; onDone: () => void }) {
  const [name, setName] = useState('')
  const [code, setCode] = useState('')
  const submit = async () => {
    await api.admin.createTest(exerciseId, { name, code, position })
    setName('')
    setCode('')
    onDone()
  }
  return (
    <form onSubmit={runner(submit)} className="mt-2 space-y-2 rounded-md border border-gray-200 p-2">
      <input className="input !py-1.5 text-sm" placeholder="Nom du test" value={name} onChange={(e) => setName(e.target.value)} required />
      <textarea className="input !py-1.5 font-mono text-sm" placeholder='Code Lua, ex: assert(...)' rows={2} value={code} onChange={(e) => setCode(e.target.value)} required />
      <button className="btn btn-ghost btn-sm w-full">Ajouter le test</button>
    </form>
  )
}

// runner wraps an async submit handler with preventDefault + error alerting.
function runner(fn: () => Promise<void>) {
  return (e: React.FormEvent) => {
    e.preventDefault()
    fn().catch((err) => alert(err instanceof Error ? err.message : 'Échec de la création'))
  }
}
