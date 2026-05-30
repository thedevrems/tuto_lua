import { useEffect, useState, type FormEvent } from 'react'
import { useSearchParams } from 'react-router-dom'
import SiteLayout from '../components/layout/SiteLayout'
import TicketConversation from '../components/TicketConversation'
import { api, type ApiTicket } from '../lib/api'

export default function SupportPage() {
  const [params, setParams] = useSearchParams()
  const [tickets, setTickets] = useState<ApiTicket[]>([])
  const [selected, setSelected] = useState<string | null>(params.get('t'))
  const [creating, setCreating] = useState(false)

  const load = () => api.tickets.mine().then(setTickets).catch(() => setTickets([]))
  useEffect(() => {
    load()
  }, [])

  const select = (id: string | null) => {
    setCreating(false)
    setSelected(id)
    setParams(id ? { t: id } : {})
  }

  return (
    <SiteLayout>
      <div className="container-page py-12">
        <h1 className="text-3xl font-black tracking-tight text-black">Support</h1>
        <p className="mt-2 text-gray-600">Posez une question, signalez un bug, suivez vos conversations.</p>

        <div className="mt-8 grid gap-6 lg:grid-cols-[320px_1fr]">
          <aside>
            <button
              onClick={() => {
                setCreating(true)
                setSelected(null)
              }}
              className="btn btn-primary w-full"
            >
              + Nouveau report
            </button>
            <ul className="mt-4 space-y-1">
              {tickets.map((t) => (
                <li key={t.id}>
                  <button
                    onClick={() => select(t.id)}
                    className={
                      'w-full rounded-md px-3 py-2 text-left text-sm transition-colors ' +
                      (t.id === selected ? 'bg-black text-white' : 'hover:bg-gray-100')
                    }
                  >
                    <span className="block font-medium">{t.subject}</span>
                    <span className={'text-xs ' + (t.id === selected ? 'text-gray-300' : 'text-gray-500')}>
                      {t.status === 'open' ? 'Ouvert' : 'Fermé'} · {new Date(t.updatedAt).toLocaleDateString('fr-FR')}
                    </span>
                  </button>
                </li>
              ))}
              {tickets.length === 0 && <li className="px-3 py-2 text-sm text-gray-400">Aucune conversation.</li>}
            </ul>
          </aside>

          <section className="card min-h-[480px]">
            {creating ? (
              <NewReportForm
                onCreated={(t) => {
                  load()
                  select(t.id)
                }}
              />
            ) : selected ? (
              <TicketConversation ticketId={selected} onChanged={load} />
            ) : (
              <div className="grid h-full place-items-center text-sm text-gray-400">
                Sélectionnez une conversation ou créez un report.
              </div>
            )}
          </section>
        </div>
      </div>
    </SiteLayout>
  )
}

function NewReportForm({ onCreated }: { onCreated: (t: ApiTicket) => void }) {
  const [subject, setSubject] = useState('')
  const [category, setCategory] = useState('question')
  const [body, setBody] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const submit = async (e: FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const t = await api.tickets.create(subject, category, body)
      onCreated(t)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Échec de la création')
      setBusy(false)
    }
  }

  return (
    <form onSubmit={submit} className="space-y-4">
      <h3 className="text-lg font-semibold text-black">Nouveau report</h3>
      {error && <div className="rounded-md border border-danger-border bg-danger-bg px-4 py-2 text-sm text-danger">{error}</div>}
      <input className="input" placeholder="Sujet" value={subject} onChange={(e) => setSubject(e.target.value)} required />
      <select className="input" value={category} onChange={(e) => setCategory(e.target.value)}>
        <option value="question">Question</option>
        <option value="bug">Bug</option>
        <option value="autre">Autre</option>
      </select>
      <textarea className="input" rows={5} placeholder="Décrivez votre demande…" value={body} onChange={(e) => setBody(e.target.value)} required />
      <button className="btn btn-primary" disabled={busy}>
        {busy ? 'Envoi…' : 'Créer le report'}
      </button>
    </form>
  )
}
