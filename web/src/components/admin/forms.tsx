import { useState, type FormEvent, type ReactNode } from 'react'
import type { NewExercise, NewLesson, NewTest } from '../../lib/api'

// Shared submit wrapper: prevents default, alerts on failure, resets create forms.
function useSubmit(isEdit: boolean, run: () => Promise<void>, reset: () => void) {
  return (e: FormEvent) => {
    e.preventDefault()
    run()
      .then(() => {
        if (!isEdit) reset()
      })
      .catch((err) => alert(err instanceof Error ? err.message : 'Échec de l’opération'))
  }
}

function Actions({ isEdit, onCancel }: { isEdit: boolean; onCancel?: () => void }) {
  return (
    <div className="flex gap-2">
      <button className="btn btn-secondary btn-sm flex-1">{isEdit ? 'Enregistrer' : 'Ajouter'}</button>
      {isEdit && onCancel && (
        <button type="button" onClick={onCancel} className="btn btn-ghost btn-sm">
          Annuler
        </button>
      )}
    </div>
  )
}

function FieldsetTitle({ children }: { children: ReactNode }) {
  return <div className="text-xs font-semibold uppercase tracking-wide text-gray-500">{children}</div>
}

export function LessonForm({
  initial,
  position,
  onSubmit,
  onCancel,
}: {
  initial?: NewLesson
  position: number
  onSubmit: (v: NewLesson) => Promise<void>
  onCancel?: () => void
}) {
  const isEdit = !!initial
  const [title, setTitle] = useState(initial?.title ?? '')
  const [content, setContent] = useState(initial?.content ?? '')
  const submit = useSubmit(isEdit, () => onSubmit({ title, content, position }), () => {
    setTitle('')
    setContent('')
  })
  return (
    <form onSubmit={submit} className="space-y-2 rounded-md border border-gray-200 p-3">
      <FieldsetTitle>{isEdit ? 'Modifier la leçon' : 'Nouvelle leçon'}</FieldsetTitle>
      <input className="input !py-1.5 text-sm" placeholder="Titre" value={title} onChange={(e) => setTitle(e.target.value)} required />
      <textarea className="input !py-1.5 text-sm" placeholder="Contenu (markdown)" rows={3} value={content} onChange={(e) => setContent(e.target.value)} />
      <Actions isEdit={isEdit} onCancel={onCancel} />
    </form>
  )
}

export function ExerciseForm({
  initial,
  position,
  onSubmit,
  onCancel,
}: {
  initial?: NewExercise
  position: number
  onSubmit: (v: NewExercise) => Promise<void>
  onCancel?: () => void
}) {
  const isEdit = !!initial
  const [f, setF] = useState({
    title: initial?.title ?? '',
    difficulty: initial?.difficulty ?? 'facile',
    statement: initial?.statement ?? '',
    starter: initial?.starter ?? '',
    solution: initial?.solution ?? '',
    hints: (initial?.hints ?? []).join('\n'),
  })
  const set = (k: keyof typeof f) => (e: { target: { value: string } }) => setF((s) => ({ ...s, [k]: e.target.value }))
  const submit = useSubmit(
    isEdit,
    () =>
      onSubmit({
        title: f.title,
        difficulty: f.difficulty as NewExercise['difficulty'],
        statement: f.statement,
        starter: f.starter,
        solution: f.solution,
        hints: f.hints.split('\n').map((h) => h.trim()).filter(Boolean),
        position,
      }),
    () => setF({ title: '', difficulty: 'facile', statement: '', starter: '', solution: '', hints: '' }),
  )
  return (
    <form onSubmit={submit} className="space-y-2 rounded-md border border-gray-200 p-3">
      <FieldsetTitle>{isEdit ? 'Modifier l’exercice' : 'Nouvel exercice'}</FieldsetTitle>
      <input className="input !py-1.5 text-sm" placeholder="Titre" value={f.title} onChange={set('title')} required />
      <select className="input !py-1.5 text-sm" value={f.difficulty} onChange={set('difficulty')}>
        <option value="facile">Facile</option>
        <option value="moyen">Moyen</option>
        <option value="difficile">Difficile</option>
      </select>
      <textarea className="input !py-1.5 text-sm" placeholder="Énoncé (markdown)" rows={2} value={f.statement} onChange={set('statement')} />
      <textarea className="input !py-1.5 font-mono text-sm" placeholder="Code de départ" rows={2} value={f.starter} onChange={set('starter')} />
      <textarea className="input !py-1.5 font-mono text-sm" placeholder="Solution" rows={2} value={f.solution} onChange={set('solution')} />
      <textarea className="input !py-1.5 text-sm" placeholder="Indices (un par ligne)" rows={2} value={f.hints} onChange={set('hints')} />
      <Actions isEdit={isEdit} onCancel={onCancel} />
    </form>
  )
}

export function TestForm({
  initial,
  position,
  onSubmit,
  onCancel,
}: {
  initial?: NewTest
  position: number
  onSubmit: (v: NewTest) => Promise<void>
  onCancel?: () => void
}) {
  const isEdit = !!initial
  const [name, setName] = useState(initial?.name ?? '')
  const [code, setCode] = useState(initial?.code ?? '')
  const submit = useSubmit(isEdit, () => onSubmit({ name, code, position }), () => {
    setName('')
    setCode('')
  })
  return (
    <form onSubmit={submit} className="mt-2 space-y-2 rounded-md border border-gray-200 p-2">
      <input className="input !py-1.5 text-sm" placeholder="Nom du test" value={name} onChange={(e) => setName(e.target.value)} required />
      <textarea className="input !py-1.5 font-mono text-sm" placeholder="Code Lua, ex: assert(...)" rows={2} value={code} onChange={(e) => setCode(e.target.value)} required />
      <Actions isEdit={isEdit} onCancel={onCancel} />
    </form>
  )
}
