import { useEffect, useState } from 'react'
import { api, type ApiTicket, type User } from '../lib/api'
import { useAuth } from '../auth/AuthContext'

const STATUS_LABEL: Record<string, string> = {
  open: 'Ouvert',
  closed: 'Fermé',
  accepted: 'Accepté',
  refused: 'Refusé',
}

interface Props {
  ticketId: string
  /** When provided (admin), enables the "add member" picker. */
  users?: User[]
  /** Called after any change so parent lists can refresh. */
  onChanged?: () => void
}

/** Shared conversation view for reports and devis (user and admin). */
export default function TicketConversation({ ticketId, users, onChanged }: Props) {
  const { user } = useAuth()
  const [ticket, setTicket] = useState<ApiTicket | null>(null)
  const [reply, setReply] = useState('')
  const [busy, setBusy] = useState(false)

  const load = () => api.tickets.get(ticketId).then(setTicket).catch(() => setTicket(null))
  useEffect(() => {
    setTicket(null)
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ticketId])

  if (!ticket) return <div className="grid h-40 place-items-center text-sm text-gray-400">Chargement…</div>

  const isAdmin = user?.role === 'admin'
  const open = ticket.status === 'open'

  const send = async () => {
    if (!reply.trim()) return
    setBusy(true)
    try {
      await api.tickets.postMessage(ticketId, reply)
      setReply('')
      await load()
      onChanged?.()
    } finally {
      setBusy(false)
    }
  }

  const close = async () => {
    await api.admin.closeTicket(ticketId).catch(() => {})
    await load()
    onChanged?.()
  }

  const addMember = async (userId: string) => {
    if (!userId) return
    await api.admin.addTicketMember(ticketId, userId).catch((e) => alert(e.message))
    await load()
  }

  return (
    <div className="flex h-full flex-col">
      <header className="border-b border-gray-200 pb-3">
        <div className="flex items-center gap-2">
          <h3 className="text-lg font-semibold text-black">{ticket.subject}</h3>
          <StatusBadge status={ticket.status} />
        </div>
        <p className="mt-1 text-xs text-gray-500">
          Ouvert par {ticket.creatorName}
          {ticket.category && ` · ${ticket.category}`}
          {' · '}
          {new Date(ticket.createdAt).toLocaleString('fr-FR')}
        </p>
        <Members members={ticket.members ?? []} />
      </header>

      <div className="flex-1 space-y-3 overflow-y-auto py-4">
        {(ticket.messages ?? []).map((m) => (
          <div key={m.id} className="rounded-md border border-gray-200 p-3">
            <div className="flex items-center justify-between text-xs text-gray-500">
              <span className="font-medium text-black">{m.authorName}</span>
              <span>{new Date(m.createdAt).toLocaleString('fr-FR')}</span>
            </div>
            <p className="mt-1 whitespace-pre-wrap text-sm text-gray-800">{m.body}</p>
          </div>
        ))}
      </div>

      {open ? (
        <div className="border-t border-gray-200 pt-3">
          <textarea
            className="input text-sm"
            rows={3}
            placeholder="Votre réponse…"
            value={reply}
            onChange={(e) => setReply(e.target.value)}
          />
          <div className="mt-2 flex flex-wrap items-center gap-2">
            <button onClick={send} disabled={busy} className="btn btn-primary btn-sm">
              {busy ? 'Envoi…' : 'Envoyer'}
            </button>
            {isAdmin && (
              <>
                <button onClick={close} className="btn btn-secondary btn-sm">Fermer la conversation</button>
                {users && (
                  <select
                    className="input !w-auto !py-1.5 text-sm"
                    defaultValue=""
                    onChange={(e) => {
                      addMember(e.target.value)
                      e.currentTarget.value = ''
                    }}
                  >
                    <option value="" disabled>+ Ajouter un membre</option>
                    {users.map((u) => (
                      <option key={u.id} value={u.id}>{u.username}</option>
                    ))}
                  </select>
                )}
              </>
            )}
          </div>
        </div>
      ) : (
        <div className="border-t border-gray-200 pt-3 text-sm text-gray-500">Cette conversation est {STATUS_LABEL[ticket.status]?.toLowerCase() ?? 'close'}.</div>
      )}
    </div>
  )
}

function StatusBadge({ status }: { status: string }) {
  const cls = status === 'open' ? 'badge-success' : status === 'refused' ? 'badge-danger' : ''
  return <span className={'badge ' + cls}>{STATUS_LABEL[status] ?? status}</span>
}

function Members({ members }: { members: { userId: string; username: string }[] }) {
  if (members.length === 0) return null
  return (
    <div className="mt-2 flex flex-wrap items-center gap-1 text-xs text-gray-500">
      Participants :
      {members.map((m) => (
        <span key={m.userId} className="badge">{m.username}</span>
      ))}
    </div>
  )
}
