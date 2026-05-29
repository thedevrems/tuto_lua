import { Link } from 'react-router-dom'
import SiteLayout from '../components/layout/SiteLayout'
import { useAuth } from '../auth/AuthContext'

interface Tier {
  name: string
  price: string
  period: string
  description: string
  features: string[]
  highlighted?: boolean
}

const tiers: Tier[] = [
  {
    name: 'Découverte',
    price: '0 €',
    period: 'pour toujours',
    description: 'Goûtez à la plateforme avec le premier module.',
    features: ['Module 1 — Fondamentaux', 'Éditeur de code & console', 'Tests automatiques', 'Progression locale'],
  },
  {
    name: 'Complet',
    price: '49 €',
    period: 'accès à vie',
    description: 'Tous les modules, mises à jour incluses.',
    features: [
      'Tous les modules (1 à 4)',
      'Exercices & solutions complets',
      'Progression synchronisée',
      'Spécial FiveM : natives & projet',
      'Mises à jour à vie',
    ],
    highlighted: true,
  },
  {
    name: 'Équipe',
    price: 'Sur devis',
    period: 'par organisation',
    description: 'Pour les serveurs et studios FiveM.',
    features: ['Tout le contenu', 'Comptes multiples', 'Suivi de progression d’équipe', 'Support prioritaire'],
  },
]

export default function PricingPage() {
  const { user } = useAuth()
  const ctaTarget = user ? '/learn' : '/register'

  return (
    <SiteLayout>
      <section className="container-page py-20 text-center">
        <span className="badge mb-6">Tarifs</span>
        <h1 className="text-5xl font-black tracking-tight text-black">Un prix simple, un accès complet</h1>
        <p className="mx-auto mt-4 max-w-2xl text-lg text-gray-600">
          Payez une fois, apprenez à votre rythme. Pas d’abonnement caché.
        </p>
      </section>

      <section className="container-page pb-24">
        <div className="grid items-start gap-6 lg:grid-cols-3">
          {tiers.map((tier) => (
            <article
              key={tier.name}
              className={
                'card flex flex-col ' +
                (tier.highlighted ? 'border-black shadow-lg lg:-translate-y-2' : '')
              }
            >
              {tier.highlighted && <span className="badge badge-success mb-4 self-start">Le plus populaire</span>}
              <h2 className="text-xl font-bold text-black">{tier.name}</h2>
              <p className="mt-1 text-sm text-gray-600">{tier.description}</p>
              <div className="mt-6 flex items-baseline gap-2">
                <span className="text-4xl font-black tracking-tight text-black">{tier.price}</span>
                <span className="text-sm text-gray-500">{tier.period}</span>
              </div>
              <ul className="mt-6 flex-1 space-y-3">
                {tier.features.map((f) => (
                  <li key={f} className="flex items-start gap-2.5 text-sm text-gray-700">
                    <span className="mt-0.5 text-success">✓</span>
                    {f}
                  </li>
                ))}
              </ul>
              <Link
                to={ctaTarget}
                className={'btn mt-8 w-full ' + (tier.highlighted ? 'btn-primary' : 'btn-secondary')}
              >
                {tier.price === 'Sur devis' ? 'Nous contacter' : user ? 'Accéder aux cours' : 'Commencer'}
              </Link>
            </article>
          ))}
        </div>

        <p className="mx-auto mt-12 max-w-xl text-center text-sm text-gray-500">
          Paiement sécurisé par carte bancaire (Stripe). L’accès au cours est débloqué automatiquement après l’achat.
        </p>
      </section>
    </SiteLayout>
  )
}
