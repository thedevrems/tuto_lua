import type { ReactNode } from 'react'
import { Link } from 'react-router-dom'

interface Props {
  title: string
  subtitle?: string
  children: ReactNode
}

/** Centered, minimal layout used by the login and register pages. */
export default function AuthShell({ title, subtitle, children }: Props) {
  return (
    <div className="flex min-h-screen flex-col bg-white-soft">
      <div className="container-page py-6">
        <Link to="/" className="inline-flex items-center gap-2.5">
          <span className="grid h-8 w-8 place-items-center rounded-md bg-black font-mono text-sm font-bold text-white">
            L
          </span>
          <span className="text-[15px] font-bold tracking-tight text-black">Lua Academy</span>
        </Link>
      </div>

      <div className="flex flex-1 items-center justify-center px-6 pb-16">
        <div className="w-full max-w-md animate-fade-up">
          <div className="mb-8 text-center">
            <h1 className="text-3xl font-black tracking-tight text-black">{title}</h1>
            {subtitle && <p className="mt-2 text-gray-600">{subtitle}</p>}
          </div>
          <div className="card shadow-lg">{children}</div>
        </div>
      </div>
    </div>
  )
}
