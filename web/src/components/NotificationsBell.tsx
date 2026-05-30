import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { api, type ApiNotification } from '../lib/api'

/** Bell with unread badge and a dropdown of the current user's notifications. */
export default function NotificationsBell() {
  const navigate = useNavigate()
  const [items, setItems] = useState<ApiNotification[]>([])
  const [unread, setUnread] = useState(0)
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  const load = () =>
    api.notifications
      .list()
      .then((r) => {
        setItems(r.notifications)
        setUnread(r.unread)
      })
      .catch(() => {})

  useEffect(() => {
    load()
    const timer = setInterval(load, 60000) // light polling
    return () => clearInterval(timer)
  }, [])

  useEffect(() => {
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', onClick)
    return () => document.removeEventListener('mousedown', onClick)
  }, [])

  const openNotif = async (n: ApiNotification) => {
    if (!n.read) {
      await api.notifications.markRead(n.id).catch(() => {})
      load()
    }
    if (n.link) {
      setOpen(false)
      navigate(n.link)
    }
  }

  const markAll = async () => {
    await api.notifications.markAllRead().catch(() => {})
    load()
  }

  return (
    <div className="relative" ref={ref}>
      <button
        onClick={() => {
          setOpen((o) => !o)
          load()
        }}
        className="relative px-2 py-1 text-lg leading-none text-gray-600 hover:text-black"
        aria-label="Notifications"
      >
        🔔
        {unread > 0 && (
          <span className="absolute -right-0.5 -top-0.5 grid h-4 min-w-[16px] place-items-center rounded-full bg-danger px-1 text-[10px] font-bold text-white">
            {unread}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-40 mt-2 w-80 rounded-lg border border-gray-200 bg-white shadow-lg">
          <div className="flex items-center justify-between border-b border-gray-200 px-4 py-2">
            <span className="text-sm font-semibold text-black">Notifications</span>
            {unread > 0 && (
              <button onClick={markAll} className="text-xs text-gray-500 hover:text-black">
                Tout marquer lu
              </button>
            )}
          </div>
          <div className="max-h-96 overflow-y-auto">
            {items.length === 0 ? (
              <div className="px-4 py-6 text-center text-sm text-gray-400">Aucune notification</div>
            ) : (
              items.map((n) => (
                <button
                  key={n.id}
                  onClick={() => openNotif(n)}
                  className={'block w-full border-b border-gray-100 px-4 py-3 text-left hover:bg-gray-50 ' + (n.read ? '' : 'bg-gray-50')}
                >
                  <div className="flex items-center gap-2">
                    {!n.read && <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-danger" />}
                    <span className="text-sm font-medium text-black">{n.title}</span>
                  </div>
                  {n.body && <p className="mt-0.5 text-xs text-gray-600">{n.body}</p>}
                  <span className="mt-1 block text-[10px] text-gray-400">
                    {new Date(n.createdAt).toLocaleString('fr-FR')}
                  </span>
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
