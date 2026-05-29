import { Link } from 'react-router-dom'
import type { LockedCourse } from '../content/useCurriculum'

function formatPrice(cents: number, currency: string): string {
  return (cents / 100).toLocaleString('fr-FR', { style: 'currency', currency: currency.toUpperCase() })
}

/** Shown in the workspace when the user selects a course they haven't unlocked. */
export default function Paywall({ course }: { course: LockedCourse }) {
  return (
    <div className="grid h-full place-items-center px-6 text-center">
      <div className="card max-w-md">
        <span className="badge mb-4">Cours verrouillé</span>
        <h2 className="text-2xl font-bold text-black">{course.title}</h2>
        {course.summary && <p className="mt-2 text-gray-600">{course.summary}</p>}
        <div className="mt-6 text-3xl font-black tracking-tight text-black">
          {formatPrice(course.priceCents, course.currency)}
        </div>
        <Link to="/pricing" className="btn btn-primary mt-6 w-full">
          Débloquer ce cours
        </Link>
        <p className="mt-3 text-xs text-gray-500">Paiement sécurisé. Accès immédiat après l'achat.</p>
      </div>
    </div>
  )
}
