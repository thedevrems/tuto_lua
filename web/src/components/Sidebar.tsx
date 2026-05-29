import { useState } from 'react'
import type { Module, Item } from '../types'

interface Props {
  modules: Module[]
  activeId: string
  completed: Set<string>
  onSelect: (item: Item) => void
}

function itemBadge(item: Item): string {
  if (item.kind === 'lesson') return 'Cours'
  return { facile: 'Facile', moyen: 'Moyen', difficile: 'Difficile' }[item.difficulty]
}

export default function Sidebar({ modules, activeId, completed, onSelect }: Props) {
  // Expand the module that contains the active item by default.
  const initiallyOpen = modules
    .filter((m) => m.chapters.some((c) => c.items.some((i) => i.id === activeId)))
    .map((m) => m.id)
  const [open, setOpen] = useState<Set<string>>(new Set(initiallyOpen.length ? initiallyOpen : [modules[0]?.id]))

  const toggle = (id: string) =>
    setOpen((prev) => {
      const next = new Set(prev)
      next.has(id) ? next.delete(id) : next.add(id)
      return next
    })

  return (
    <nav className="h-full overflow-y-auto px-3 py-4">
      {modules.map((module) => {
        const isOpen = open.has(module.id)
        return (
          <div key={module.id} className="mb-1.5">
            <button
              onClick={() => toggle(module.id)}
              className="w-full flex items-center gap-2 text-left px-2 py-2 rounded-md hover:bg-gray-100 transition-colors"
            >
              <span className={'text-gray-400 text-xs transition-transform ' + (isOpen ? 'rotate-90' : '')}>
                ▶
              </span>
              <span className="text-[13px] font-semibold text-black leading-tight">{module.title}</span>
            </button>

            {isOpen && (
              <div className="mt-0.5 ml-3 pl-3 border-l border-gray-200 space-y-3 pb-2">
                {module.chapters.map((chapter) => (
                  <div key={chapter.id}>
                    <div className="flex items-center gap-2 px-2 pt-2 pb-1">
                      <span className="text-[11px] font-medium uppercase tracking-wide text-gray-500 leading-tight">
                        {chapter.title}
                      </span>
                      {chapter.comingSoon && (
                        <span className="text-[10px] px-1.5 py-px rounded-full border border-gray-300 text-gray-400">
                          bientôt
                        </span>
                      )}
                    </div>

                    {chapter.items.map((item) => {
                      const active = item.id === activeId
                      const done = completed.has(item.id)
                      return (
                        <button
                          key={item.id}
                          onClick={() => onSelect(item)}
                          className={
                            'w-full flex items-center gap-2 text-left px-2 py-1.5 rounded-md text-[13px] transition-colors ' +
                            (active
                              ? 'bg-black text-white font-medium'
                              : 'text-gray-600 hover:bg-gray-100 hover:text-black')
                          }
                        >
                          <span
                            className={
                              'w-1.5 h-1.5 rounded-full shrink-0 ' +
                              (done
                                ? active
                                  ? 'bg-white'
                                  : 'bg-success'
                                : active
                                  ? 'bg-gray-400'
                                  : 'bg-gray-300')
                            }
                          />
                          <span className="flex-1 leading-tight">{item.title}</span>
                          <span className={'text-[10px] shrink-0 ' + (active ? 'text-gray-300' : 'text-gray-400')}>
                            {itemBadge(item)}
                          </span>
                        </button>
                      )
                    })}
                  </div>
                ))}
              </div>
            )}
          </div>
        )
      })}
    </nav>
  )
}
