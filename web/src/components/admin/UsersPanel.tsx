import { useEffect, useState } from 'react'
import { api, type ApiCourse, type ApiProgress, type User } from '../../lib/api'

export default function UsersPanel() {
  const [users, setUsers] = useState<User[]>([])
  const [courses, setCourses] = useState<ApiCourse[]>([])
  const [grantSel, setGrantSel] = useState<Record<string, string>>({})
  const [notice, setNotice] = useState<string | null>(null)
  const [viewing, setViewing] = useState<{ user: User; rows: ApiProgress[] } | null>(null)

  useEffect(() => {
    Promise.all([api.admin.users(), api.courses.list()])
      .then(([u, c]) => {
        setUsers(u)
        setCourses(c)
      })
      .catch((e) => setNotice(e.message))
  }, [])

  const grant = async (userId: string) => {
    const courseId = grantSel[userId] || courses[0]?.id
    if (!courseId) return
    try {
      await api.admin.grant(userId, courseId)
      setNotice('Accès accordé.')
    } catch (e) {
      setNotice(e instanceof Error ? e.message : 'Échec')
    }
  }

  const viewCode = async (user: User) => {
    setViewing({ user, rows: [] })
    try {
      setViewing({ user, rows: await api.admin.userProgress(user.id) })
    } catch (e) {
      setNotice(e instanceof Error ? e.message : 'Échec')
    }
  }

  return (
    <div className="space-y-6">
      {notice && <div className="rounded-md border border-gray-200 bg-gray-100 px-4 py-2 text-sm text-gray-700">{notice}</div>}

      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-200 text-left text-gray-500">
              <th className="py-2 pr-4 font-medium">Utilisateur</th>
              <th className="py-2 pr-4 font-medium">E-mail</th>
              <th className="py-2 pr-4 font-medium">Rôle</th>
              <th className="py-2 pr-4 font-medium">Donner l'accès à un cours</th>
              <th className="py-2 font-medium">Code</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id} className="border-b border-gray-100">
                <td className="py-2 pr-4 font-medium text-black">{u.username}</td>
                <td className="py-2 pr-4 text-gray-600">{u.email}</td>
                <td className="py-2 pr-4">
                  <span className={'badge ' + (u.role === 'admin' ? 'badge-warning' : '')}>{u.role}</span>
                </td>
                <td className="py-2 pr-4">
                  <div className="flex items-center gap-2">
                    <select
                      className="input !py-1.5 !w-auto text-sm"
                      value={grantSel[u.id] ?? ''}
                      onChange={(e) => setGrantSel((s) => ({ ...s, [u.id]: e.target.value }))}
                    >
                      {courses.map((c) => (
                        <option key={c.id} value={c.id}>
                          {c.title}
                        </option>
                      ))}
                    </select>
                    <button onClick={() => grant(u.id)} className="btn btn-secondary btn-sm whitespace-nowrap">
                      Donner l'accès
                    </button>
                  </div>
                </td>
                <td className="py-2">
                  <button onClick={() => viewCode(u)} className="btn btn-ghost btn-sm">
                    Voir le code
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {viewing && <StudentCode user={viewing.user} rows={viewing.rows} onClose={() => setViewing(null)} />}
    </div>
  )
}

/** Read-only viewer for a student's latest pushed code per exercise. */
function StudentCode({ user, rows, onClose }: { user: User; rows: ApiProgress[]; onClose: () => void }) {
  return (
    <div className="card">
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-lg font-semibold text-black">
          Dernier code de <span className="font-mono">{user.username}</span>
        </h3>
        <button onClick={onClose} className="btn btn-ghost btn-sm">
          Fermer
        </button>
      </div>
      {rows.length === 0 ? (
        <p className="text-sm text-gray-500">Aucune progression enregistrée pour cet utilisateur.</p>
      ) : (
        <div className="space-y-4">
          {rows.map((p) => (
            <div key={p.id}>
              <div className="mb-1 flex items-center gap-2 text-xs text-gray-500">
                <span className="font-mono">exercice {p.exerciseId.slice(0, 8)}</span>
                {p.completed && <span className="badge badge-success">réussi</span>}
                <span className="ml-auto">{new Date(p.updatedAt).toLocaleString('fr-FR')}</span>
              </div>
              <pre className="overflow-x-auto rounded-md bg-black p-3 font-mono text-[12.5px] text-white">
                {p.code || '(vide)'}
              </pre>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
