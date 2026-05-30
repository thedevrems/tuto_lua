import { useState } from 'react'
import { api, type ApiChapter, type ApiExercise, type ApiLesson, type ApiTest } from '../../lib/api'
import { ExerciseForm, LessonForm, TestForm } from './forms'

const alertErr = (e: unknown) => alert(e instanceof Error ? e.message : 'Échec de l’opération')
const confirmDelete = (label: string) => window.confirm(`Supprimer ${label} ? Cette action est irréversible.`)

/** Authoring UI for a single chapter: edit/delete its lessons, exercises and tests. */
export default function ChapterEditor({ chapter, onChange }: { chapter: ApiChapter; onChange: () => void }) {
  const lessons = chapter.lessons ?? []
  const exercises = chapter.exercises ?? []
  const nextPos = lessons.length + exercises.length
  const [renaming, setRenaming] = useState(false)

  const deleteChapter = () => {
    if (!confirmDelete(`le chapitre « ${chapter.title} »`)) return
    api.admin.deleteChapter(chapter.id).then(onChange).catch(alertErr)
  }

  return (
    <div className="card">
      <div className="flex items-center justify-between gap-2">
        <h4 className="font-semibold text-black">{chapter.title}</h4>
        <div className="flex gap-1">
          <button onClick={() => setRenaming((r) => !r)} className="btn btn-ghost btn-sm">Renommer</button>
          <button onClick={deleteChapter} className="btn btn-ghost btn-sm text-danger">Supprimer</button>
        </div>
      </div>
      {renaming && <RenameChapter chapter={chapter} onDone={() => { setRenaming(false); onChange() }} />}

      <div className="mt-3 space-y-2">
        {lessons.map((l) => (
          <LessonRow key={l.id} lesson={l} onChange={onChange} />
        ))}
        {exercises.map((e) => (
          <ExerciseBlock key={e.id} ex={e} onChange={onChange} />
        ))}
      </div>

      <div className="mt-4 grid gap-3 sm:grid-cols-2">
        <LessonForm position={nextPos} onSubmit={(v) => api.admin.createLesson(chapter.id, v).then(onChange)} />
        <ExerciseForm position={nextPos} onSubmit={(v) => api.admin.createExercise(chapter.id, v).then(onChange)} />
      </div>
    </div>
  )
}

function RenameChapter({ chapter, onDone }: { chapter: ApiChapter; onDone: () => void }) {
  const [title, setTitle] = useState(chapter.title)
  const [summary, setSummary] = useState(chapter.summary)
  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    api.admin.updateChapter(chapter.id, { title, summary, position: chapter.position }).then(onDone).catch(alertErr)
  }
  return (
    <form onSubmit={submit} className="mt-2 flex flex-wrap gap-2">
      <input className="input !py-1.5 flex-1 text-sm" value={title} onChange={(e) => setTitle(e.target.value)} required />
      <input className="input !py-1.5 flex-1 text-sm" placeholder="Résumé" value={summary} onChange={(e) => setSummary(e.target.value)} />
      <button className="btn btn-secondary btn-sm">Enregistrer</button>
    </form>
  )
}

function LessonRow({ lesson, onChange }: { lesson: ApiLesson; onChange: () => void }) {
  const [edit, setEdit] = useState(false)
  const del = () => {
    if (!confirmDelete(`la leçon « ${lesson.title} »`)) return
    api.admin.deleteLesson(lesson.id).then(onChange).catch(alertErr)
  }
  if (edit) {
    return (
      <LessonForm
        initial={{ title: lesson.title, content: lesson.content, position: lesson.position }}
        position={lesson.position}
        onSubmit={(v) => api.admin.updateLesson(lesson.id, v).then(() => { setEdit(false); onChange() })}
        onCancel={() => setEdit(false)}
      />
    )
  }
  return (
    <Row icon="📄" label={lesson.title} onEdit={() => setEdit(true)} onDelete={del} />
  )
}

function ExerciseBlock({ ex, onChange }: { ex: ApiExercise; onChange: () => void }) {
  const [open, setOpen] = useState(false)
  const [edit, setEdit] = useState(false)
  const tests = ex.tests ?? []
  const del = () => {
    if (!confirmDelete(`l’exercice « ${ex.title} »`)) return
    api.admin.deleteExercise(ex.id).then(onChange).catch(alertErr)
  }
  return (
    <div className="rounded-md border border-gray-200 p-2">
      <div className="flex items-center gap-2 text-sm">
        <button onClick={() => setOpen((o) => !o)} className="flex-1 text-left text-black hover:underline">
          ✏️ {ex.title} <span className="text-xs text-gray-400">({tests.length} test(s))</span>
        </button>
        <button onClick={() => setEdit((e) => !e)} className="text-gray-500 hover:text-black">Modifier</button>
        <button onClick={del} className="text-gray-500 hover:text-danger">Supprimer</button>
      </div>

      {edit && (
        <div className="mt-2">
          <ExerciseForm
            initial={{ title: ex.title, difficulty: ex.difficulty, statement: ex.statement, starter: ex.starter, solution: ex.solution ?? '', hints: ex.hints ?? [], position: ex.position }}
            position={ex.position}
            onSubmit={(v) => api.admin.updateExercise(ex.id, v).then(() => { setEdit(false); onChange() })}
            onCancel={() => setEdit(false)}
          />
        </div>
      )}

      {open && (
        <div className="mt-2 space-y-1 border-t border-gray-100 pt-2">
          {tests.map((t) => (
            <TestRow key={t.id} test={t} onChange={onChange} />
          ))}
          <TestForm position={tests.length} onSubmit={(v) => api.admin.createTest(ex.id, v).then(onChange)} />
        </div>
      )}
    </div>
  )
}

function TestRow({ test, onChange }: { test: ApiTest; onChange: () => void }) {
  const [edit, setEdit] = useState(false)
  const del = () => {
    if (!confirmDelete(`le test « ${test.name} »`)) return
    api.admin.deleteTest(test.id).then(onChange).catch(alertErr)
  }
  if (edit) {
    return (
      <TestForm
        initial={{ name: test.name, code: test.code, position: test.position }}
        position={test.position}
        onSubmit={(v) => api.admin.updateTest(test.id, v).then(() => { setEdit(false); onChange() })}
        onCancel={() => setEdit(false)}
      />
    )
  }
  return <Row icon="•" label={test.name} small onEdit={() => setEdit(true)} onDelete={del} />
}

/** A read-only content row with Modifier / Supprimer actions. */
function Row({
  icon,
  label,
  small,
  onEdit,
  onDelete,
}: {
  icon: string
  label: string
  small?: boolean
  onEdit: () => void
  onDelete: () => void
}) {
  return (
    <div className={'flex items-center gap-2 ' + (small ? 'text-xs' : 'text-sm')}>
      <span className="flex-1 text-gray-700">
        {icon} {label}
      </span>
      <button onClick={onEdit} className="text-gray-500 hover:text-black">Modifier</button>
      <button onClick={onDelete} className="text-gray-500 hover:text-danger">Supprimer</button>
    </div>
  )
}
