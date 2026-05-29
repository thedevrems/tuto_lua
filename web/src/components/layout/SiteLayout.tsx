import type { ReactNode } from 'react'
import Navbar from './Navbar'
import Footer from './Footer'

/** Wraps marketing pages with the shared navbar and footer. */
export default function SiteLayout({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col bg-white">
      <Navbar />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  )
}
