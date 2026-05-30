import { useEffect, useState, type FormEvent } from 'react'
import SiteLayout from '../components/layout/SiteLayout'
import { useAuth } from '../auth/AuthContext'
import { api, ApiError, type ApiCourse } from '../lib/api'

export default function ProfilePage() {
  const { user } = useAuth()
  const [courses, setCourses] = useState<ApiCourse[]>([])

  useEffect(() => {
    api.account.courses().then(setCourses).catch(() => setCourses([]))
  }, [])

  if (!user) return null

  return (
    <SiteLayout>
      <div className="container-page max-w-3xl py-12">
        <h1 className="text-3xl font-black tracking-tight text-black">Mon profil</h1>

        <section className="card mt-8">
          <h2 className="text-lg font-semibold text-black">Informations</h2>
          <dl className="mt-4 grid gap-3 sm:grid-cols-2">
            <Info label="Nom d'utilisateur" value={user.username} />
            <Info label="E-mail" value={user.email} />
            <Info label="Rôle" value={user.role === 'admin' ? 'Administrateur' : 'Membre'} />
            <Info label="Membre depuis" value={new Date(user.createdAt).toLocaleDateString('fr-FR')} />
          </dl>
        </section>

        <section className="card mt-6">
          <h2 className="text-lg font-semibold text-black">Mes cours</h2>
          {courses.length === 0 ? (
            <p className="mt-2 text-sm text-gray-500">Aucun cours accessible pour le moment.</p>
          ) : (
            <ul className="mt-4 space-y-2">
              {courses.map((c) => (
                <li key={c.id} className="flex items-center gap-2 text-sm text-gray-700">
                  <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-success" />
                  {c.title}
                  {c.priceCents === 0 && <span className="badge">gratuit</span>}
                </li>
              ))}
            </ul>
          )}
        </section>

        <ChangePasswordForm />
      </div>
    </SiteLayout>
  )
}

function Info({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs uppercase tracking-wide text-gray-500">{label}</dt>
      <dd className="mt-0.5 font-medium text-black">{value}</dd>
    </div>
  )
}

function ChangePasswordForm() {
  const [current, setCurrent] = useState('')
  const [next, setNext] = useState('')
  const [confirm, setConfirm] = useState('')
  const [feedback, setFeedback] = useState<{ ok: boolean; msg: string } | null>(null)
  const [busy, setBusy] = useState(false)

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setFeedback(null)
    if (next !== confirm) {
      setFeedback({ ok: false, msg: 'Les nouveaux mots de passe ne correspondent pas.' })
      return
    }
    setBusy(true)
    try {
      await api.account.changePassword(current, next)
      setFeedback({ ok: true, msg: 'Mot de passe mis à jour.' })
      setCurrent('')
      setNext('')
      setConfirm('')
    } catch (err) {
      setFeedback({ ok: false, msg: err instanceof ApiError ? err.message : 'Échec.' })
    } finally {
      setBusy(false)
    }
  }

  return (
    <section className="card mt-6">
      <h2 className="text-lg font-semibold text-black">Changer le mot de passe</h2>
      <form onSubmit={submit} className="mt-4 space-y-4">
        {feedback && (
          <div
            className={
              'rounded-md border px-4 py-2 text-sm ' +
              (feedback.ok
                ? 'border-success-border bg-success-bg text-success'
                : 'border-danger-border bg-danger-bg text-danger')
            }
          >
            {feedback.msg}
          </div>
        )}
        <input className="input" type="password" placeholder="Mot de passe actuel" autoComplete="current-password" value={current} onChange={(e) => setCurrent(e.target.value)} required />
        <input className="input" type="password" placeholder="Nouveau mot de passe" autoComplete="new-password" value={next} onChange={(e) => setNext(e.target.value)} required />
        <input className="input" type="password" placeholder="Confirmer le nouveau mot de passe" autoComplete="new-password" value={confirm} onChange={(e) => setConfirm(e.target.value)} required />
        <button className="btn btn-primary" disabled={busy}>
          {busy ? 'Mise à jour…' : 'Mettre à jour'}
        </button>
      </form>
    </section>
  )
}
