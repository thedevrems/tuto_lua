import { useState } from 'react'
import SiteLayout from '../components/layout/SiteLayout'
import UsersPanel from '../components/admin/UsersPanel'
import ContentPanel from '../components/admin/ContentPanel'

type Tab = 'users' | 'content'

export default function AdminPage() {
  const [tab, setTab] = useState<Tab>('users')

  return (
    <SiteLayout>
      <div className="container-page py-12">
        <h1 className="text-3xl font-black tracking-tight text-black">Administration</h1>
        <p className="mt-2 text-gray-600">Gérez les accès, le contenu et suivez le code des élèves.</p>

        <div className="mt-8 flex gap-1 border-b border-gray-200">
          <TabButton active={tab === 'users'} onClick={() => setTab('users')}>
            Utilisateurs & accès
          </TabButton>
          <TabButton active={tab === 'content'} onClick={() => setTab('content')}>
            Cours & tests
          </TabButton>
        </div>

        <div className="mt-8">{tab === 'users' ? <UsersPanel /> : <ContentPanel />}</div>
      </div>
    </SiteLayout>
  )
}

function TabButton({ active, onClick, children }: { active: boolean; onClick: () => void; children: React.ReactNode }) {
  return (
    <button
      onClick={onClick}
      className={
        'px-4 py-2.5 text-sm font-medium border-b-2 -mb-px transition-colors ' +
        (active ? 'border-black text-black' : 'border-transparent text-gray-500 hover:text-black')
      }
    >
      {children}
    </button>
  )
}
