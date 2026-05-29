import { Link } from 'react-router-dom'
import SiteLayout from '../components/layout/SiteLayout'
import { curriculum } from '../content/curriculum'

const features = [
  {
    title: 'Code dans le navigateur',
    body: 'Un éditeur Lua intégré qui exécute votre code instantanément (Lua 5.4 en WebAssembly). Aucune installation.',
  },
  {
    title: 'Tests automatiques',
    body: 'Chaque exercice est validé par des tests. Vous savez immédiatement si votre solution est correcte.',
  },
  {
    title: 'Pensé pour FiveM',
    body: 'Des fondamentaux jusqu’aux natives et à l’architecture client-serveur de l’écosystème FiveM.',
  },
  {
    title: 'Progression sauvegardée',
    body: 'Votre code et votre avancement sont enregistrés sur votre compte, accessibles partout.',
  },
]

export default function HomePage() {
  return (
    <SiteLayout>
      {/* ---- Hero ---- */}
      <section className="container-page py-24 text-center sm:py-32">
        <span className="badge mb-6">Lua 5.4 · pour FiveM</span>
        <h1 className="mx-auto max-w-3xl text-5xl font-black leading-tight tracking-tight text-black sm:text-6xl">
          Apprenez le Lua, du premier <span className="underline decoration-4 underline-offset-8">print()</span> à
          votre ressource FiveM.
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-gray-600">
          Une plateforme interactive : cours clairs, exercices corrigés et un éditeur de code avec tests
          automatiques. Apprenez en pratiquant, directement dans votre navigateur.
        </p>
        <div className="mt-10 flex flex-wrap items-center justify-center gap-4">
          <Link to="/register" className="btn btn-primary btn-lg">
            Commencer gratuitement
          </Link>
          <Link to="/learn" className="btn btn-secondary btn-lg">
            Découvrir les cours
          </Link>
        </div>
      </section>

      {/* ---- Features ---- */}
      <section className="border-t border-gray-200 bg-white-soft py-20">
        <div className="container-page">
          <h2 className="text-3xl font-bold tracking-tight text-black">Pourquoi Lua Academy ?</h2>
          <p className="mt-3 max-w-2xl text-gray-600">
            Tout ce qu’il faut pour passer de débutant à développeur Lua opérationnel.
          </p>
          <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {features.map((f) => (
              <article key={f.title} className="card card-interactive">
                <h3 className="text-lg font-semibold text-black">{f.title}</h3>
                <p className="mt-2 text-sm leading-relaxed text-gray-600">{f.body}</p>
              </article>
            ))}
          </div>
        </div>
      </section>

      {/* ---- Programme ---- */}
      <section className="container-page py-20">
        <h2 className="text-3xl font-bold tracking-tight text-black">Le programme</h2>
        <p className="mt-3 max-w-2xl text-gray-600">
          Quatre modules progressifs, des fondamentaux du langage à la pratique sur FiveM.
        </p>
        <div className="mt-10 grid gap-6 md:grid-cols-2">
          {curriculum.map((module, i) => (
            <article key={module.id} className="card">
              <div className="flex items-baseline gap-3">
                <span className="font-mono text-sm text-gray-400">{String(i + 1).padStart(2, '0')}</span>
                <h3 className="text-xl font-semibold text-black">{module.title}</h3>
              </div>
              {module.summary && <p className="mt-2 text-sm text-gray-600">{module.summary}</p>}
              <ul className="mt-4 space-y-1.5">
                {module.chapters.map((c) => (
                  <li key={c.id} className="flex items-center gap-2 text-sm text-gray-700">
                    <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-gray-300" />
                    {c.title}
                  </li>
                ))}
              </ul>
            </article>
          ))}
        </div>
      </section>

      {/* ---- CTA ---- */}
      <section className="bg-black py-20">
        <div className="container-page text-center">
          <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Prêt à écrire votre premier script ?
          </h2>
          <p className="mx-auto mt-4 max-w-xl text-gray-300">
            Créez votre compte en quelques secondes et démarrez le premier chapitre.
          </p>
          <Link
            to="/register"
            className="mt-8 inline-flex items-center justify-center rounded-md bg-white px-8 py-4 text-lg font-medium text-black transition-all duration-base hover:-translate-y-px hover:bg-gray-100"
          >
            Créer mon compte
          </Link>
        </div>
      </section>
    </SiteLayout>
  )
}
