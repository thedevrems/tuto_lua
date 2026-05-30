import { useEffect, useState } from 'react'
import { api, type ApiTicket, type User } from '../../lib/api'
import TicketConversation from '../TicketConversation'

/** Admin view of all support reports: list + conversation with admin controls. */
export default function TicketsPanel() {
  const [tickets, setTickets] = useState<ApiTicket[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [selected, setSelected] = useState<string | null>(null)

  const load = () => api.admin.tickets('report').then(setTickets).catch(() => setTickets([]))
  useEffect(() => {
    load()
    api.admin.users().then(setUsers).catch(() => {})
  }, [])

  return (
    <div className="grid gap-6 lg:grid-cols-[320px_1fr]">
      <aside>
        <ul className="space-y-1">
          {tickets.map((t) => (
            <li key={t.id}>
              <button
                onClick={() => setSelected(t.id)}
                className={
                  'w-full rounded-md px-3 py-2 text-left text-sm transition-colors ' +
                  (t.id === selected ? 'bg-black text-white' : 'hover:bg-gray-100')
                }
              >
                <span className="block font-medium">{t.subject}</span>
                <span className={'text-xs ' + (t.id === selected ? 'text-gray-300' : 'text-gray-500')}>
                  {t.creatorName} · {t.status === 'open' ? 'Ouvert' : 'Fermé'}
                </span>
              </button>
            </li>
          ))}
          {tickets.length === 0 && <li className="px-3 py-2 text-sm text-gray-400">Aucun report.</li>}
        </ul>
      </aside>

      <section className="card min-h-[480px]">
        {selected ? (
          <TicketConversation
            ticketId={selected}
            users={users}
            onChanged={load}
          />
        ) : (
          <div className="grid h-full place-items-center text-sm text-gray-400">Sélectionnez un report.</div>
        )}
      </section>
    </div>
  )
}
