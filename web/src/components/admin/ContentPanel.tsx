import { useEffect, useState } from 'react'
import { api, type ApiCourse } from '../../lib/api'
import ChapterEditor from './ChapterEditor'

function CourseSettings({ course, onSaved, onDeleted }: { course: ApiCourse; onSaved: () => void; onDeleted: () => void }) {
  const [f, setF] = useState({
    slug: course.slug,
    title: course.title,
    summary: course.summary,
    priceEuros: (course.priceCents / 100).toString(),
    position: course.position.toString(),
    published: course.published,
  })
  const set = (k: keyof typeof f, v: string | boolean) => setF((s) => ({ ...s, [k]: v }))

  const save = (e: React.FormEvent) => {
    e.preventDefault()
    api.admin
      .updateCourse(course.id, {
        slug: f.slug,
        title: f.title,
        summary: f.summary,
        priceCents: Math.round(parseFloat(f.priceEuros || '0') * 100),
        currency: course.currency,
        published: f.published,
        position: parseInt(f.position || '0', 10),
      })
      .then(onSaved)
      .catch((err) => alert(err instanceof Error ? err.message : 'Échec'))
  }

  const remove = () => {
    if (!window.confirm(`Supprimer le cours « ${course.title} » et tout son contenu ?`)) return
    api.admin.deleteCourse(course.id).then(onDeleted).catch((err) => alert(err instanceof Error ? err.message : 'Échec'))
  }

  return (
    <form onSubmit={save} className="card space-y-3">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-black">Paramètres du cours</h3>
        <button type="button" onClick={remove} className="btn btn-ghost btn-sm text-danger">Supprimer le cours</button>
      </div>
      <div className="grid gap-3 sm:grid-cols-2">
        <input className="input" placeholder="slug" value={f.slug} onChange={(e) => set('slug', e.target.value)} required />
        <input className="input" placeholder="Titre" value={f.title} onChange={(e) => set('title', e.target.value)} required />
      </div>
      <textarea className="input" placeholder="Résumé" rows={2} value={f.summary} onChange={(e) => set('summary', e.target.value)} />
      <div className="grid gap-3 sm:grid-cols-3">
        <input className="input" type="number" step="0.01" placeholder="Prix (€)" value={f.priceEuros} onChange={(e) => set('priceEuros', e.target.value)} />
        <input className="input" type="number" placeholder="Position" value={f.position} onChange={(e) => set('position', e.target.value)} />
        <label className="flex items-center gap-2 text-sm text-gray-700">
          <input type="checkbox" checked={f.published} onChange={(e) => set('published', e.target.checked)} />
          Publié
        </label>
      </div>
      <button className="btn btn-secondary">Enregistrer les modifications</button>
    </form>
  )
}

export default function ContentPanel() {
  const [courses, setCourses] = useState<ApiCourse[]>([])
  const [slug, setSlug] = useState('')
  const [tree, setTree] = useState<ApiCourse | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const refreshCourses = () => api.admin.courses().then(setCourses).catch((e) => setNotice(e.message))
  useEffect(() => {
    refreshCourses()
  }, [])

  const loadTree = async (s: string) => {
    setSlug(s)
    if (!s) return setTree(null)
    try {
      setTree(await api.courses.tree(s))
    } catch (e) {
      setNotice(e instanceof Error ? e.message : 'Échec du chargement')
    }
  }
  const reloadTree = () => slug && loadTree(slug)

  return (
    <div className="space-y-8">
      {notice && <div className="rounded-md border border-gray-200 bg-gray-100 px-4 py-2 text-sm text-gray-700">{notice}</div>}

      <CreateCourseForm
        onDone={() => {
          setNotice('Cours créé.')
          refreshCourses()
        }}
      />

      <div>
        <label className="label">Éditer un cours</label>
        <select className="input !w-auto" value={slug} onChange={(e) => loadTree(e.target.value)}>
          <option value="">— Sélectionner —</option>
          {courses.map((c) => (
            <option key={c.id} value={c.slug}>
              {c.title} {c.published ? '' : '(brouillon)'}
            </option>
          ))}
        </select>
      </div>

      {tree && (
        <div className="space-y-4">
          <CourseSettings
            course={tree}
            onSaved={() => {
              refreshCourses()
              reloadTree()
            }}
            onDeleted={() => {
              setSlug('')
              setTree(null)
              refreshCourses()
            }}
          />
          <AddChapter courseId={tree.id} position={(tree.chapters ?? []).length} onDone={reloadTree} />
          {(tree.chapters ?? []).map((ch) => (
            <ChapterEditor key={ch.id} chapter={ch} onChange={reloadTree} />
          ))}
        </div>
      )}
    </div>
  )
}

function CreateCourseForm({ onDone }: { onDone: () => void }) {
  const [f, setF] = useState({ slug: '', title: '', summary: '', priceEuros: '0', position: '0', published: true })
  const set = (k: keyof typeof f, v: string | boolean) => setF((s) => ({ ...s, [k]: v }))

  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    api.admin
      .createCourse({
        slug: f.slug,
        title: f.title,
        summary: f.summary,
        priceCents: Math.round(parseFloat(f.priceEuros || '0') * 100),
        currency: 'eur',
        published: f.published,
        position: parseInt(f.position || '0', 10),
      })
      .then(() => {
        setF({ slug: '', title: '', summary: '', priceEuros: '0', position: '0', published: true })
        onDone()
      })
      .catch((err) => alert(err instanceof Error ? err.message : 'Échec'))
  }

  return (
    <form onSubmit={submit} className="card space-y-3">
      <h3 className="text-lg font-semibold text-black">Créer un cours</h3>
      <div className="grid gap-3 sm:grid-cols-2">
        <input className="input" placeholder="slug (ex: m5)" value={f.slug} onChange={(e) => set('slug', e.target.value)} required />
        <input className="input" placeholder="Titre" value={f.title} onChange={(e) => set('title', e.target.value)} required />
      </div>
      <textarea className="input" placeholder="Résumé" rows={2} value={f.summary} onChange={(e) => set('summary', e.target.value)} />
      <div className="grid gap-3 sm:grid-cols-3">
        <input className="input" type="number" step="0.01" placeholder="Prix (€)" value={f.priceEuros} onChange={(e) => set('priceEuros', e.target.value)} />
        <input className="input" type="number" placeholder="Position" value={f.position} onChange={(e) => set('position', e.target.value)} />
        <label className="flex items-center gap-2 text-sm text-gray-700">
          <input type="checkbox" checked={f.published} onChange={(e) => set('published', e.target.checked)} />
          Publié
        </label>
      </div>
      <button className="btn btn-primary">Créer le cours</button>
    </form>
  )
}

function AddChapter({ courseId, position, onDone }: { courseId: string; position: number; onDone: () => void }) {
  const [title, setTitle] = useState('')
  const [summary, setSummary] = useState('')
  const submit = (e: React.FormEvent) => {
    e.preventDefault()
    api.admin
      .createChapter(courseId, { title, summary, position })
      .then(() => {
        setTitle('')
        setSummary('')
        onDone()
      })
      .catch((err) => alert(err instanceof Error ? err.message : 'Échec'))
  }
  return (
    <form onSubmit={submit} className="flex flex-wrap items-end gap-3 rounded-md border border-gray-200 p-3">
      <div className="flex-1">
        <label className="label">Nouveau chapitre</label>
        <input className="input" placeholder="Titre du chapitre" value={title} onChange={(e) => setTitle(e.target.value)} required />
      </div>
      <input className="input flex-1" placeholder="Résumé (optionnel)" value={summary} onChange={(e) => setSummary(e.target.value)} />
      <button className="btn btn-secondary">Ajouter</button>
    </form>
  )
}
