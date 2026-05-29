# Lua Academy — Cours Lua interactif pour FiveM

Plateforme web pour apprendre **Lua** dans le contexte de **FiveM** : chaque chapitre
contient un **cours**, et des **exercices** avec un éditeur de code, une **console**, des
**tests automatiques** et une **solution**.

Le code Lua de l'élève s'exécute **directement dans le navigateur** (Lua 5.4 via
WebAssembly). Un **backend Go** gère les comptes, la progression, le contenu des
cours (stocké en base) et les paiements.

## ✨ Fonctionnalités

- 📚 Cours organisés en **Modules → Chapitres → Cours / Exercices**, **stockés en base**
- 💻 Éditeur de code intégré (coloration syntaxique Lua)
- ▶️ **Console** : exécution du code et sortie de `print()`
- ✅ **Tests automatiques** par exercice (✓ / ✗, compteur de réussite)
- 💡 Indices progressifs et **solution** dépliable
- 👤 **Comptes** (inscription / connexion), mots de passe **hachés bcrypt**, sessions JWT
- 💾 **Progression par utilisateur** synchronisée (code du « dernier push » sauvegardé serveur)
- 💳 **Achat de cours** par carte (Stripe) avec déblocage automatique
- 🛠️ **Espace admin** : créer cours/tests, donner l'accès, voir le code des élèves
- 🎨 Interface **noir & blanc** moderne (charte graphique)

## 🧰 Stack technique

| Domaine | Choix |
| --- | --- |
| Front | React 18 + TypeScript + Vite + React Router |
| Style | Tailwind CSS (thème clair monochrome, charte graphique) |
| Exécution Lua | [wasmoon](https://github.com/ceifa/wasmoon) — Lua 5.4 en WebAssembly (client) |
| Éditeur | CodeMirror 6 |
| Backend | **Go** (net/http + [chi](https://github.com/go-chi/chi)) + **SQLite** (pure-Go) |
| Auth | bcrypt + JWT (HS256) |
| Paiement | Stripe Checkout + webhook |

## 🚀 Démarrage rapide

Deux services à lancer (deux terminaux) :

```bash
# Terminal 1 — backend (API + base SQLite)
cd server
cp .env.example .env
go run ./cmd/api          # http://localhost:8080  (le 1er compte créé = admin)

# Terminal 2 — frontend
cd web
npm install
npm run dev               # http://localhost:5173
```

Autres commandes :

```bash
cd web  && npm run build  # build de production -> web/dist/
cd server && go test ./...# tests unitaires du backend
```

> Détails du backend (endpoints, architecture, Stripe) : [server/README.md](server/README.md).

## 📂 Structure du projet

```
.
├── server/                    # Backend Go (API REST + SQLite)
│   ├── cmd/api/               #   point d'entrée
│   ├── internal/              #   config, database, models, auth, store, handlers, payment…
│   └── README.md              #   endpoints, architecture, Stripe
├── web/                       # Application (plateforme interactive)
│   ├── src/
│   │   ├── content/           #   curriculum source + mapping API (fromApi, useCurriculum)
│   │   ├── pages/             #   Home, Pricing, Login, Register, Learn, Admin
│   │   ├── components/        #   UI (sidebar, éditeur, console, tests, admin…)
│   │   ├── auth/              #   contexte d'authentification
│   │   └── lib/               #   client API + exécution Lua (wasmoon)
│   ├── scripts/               #   export-curriculum.mjs (génère le seed du backend)
│   └── README.md              #   Guide : lancer & ajouter du contenu
└── MODULE_1/                  # Sources Markdown d'origine (legacy)
```

> ✍️ Pour **ajouter une leçon, un exercice ou des tests**, voir
> [web/README.md](web/README.md).

## 🗺️ Programme

### Module 1 — Fondamentaux de Lua
- ✅ **Chapitre 1** — Introduction et premiers pas (variables, types, opérateurs)
- ✅ **Chapitre 2** — Structures de contrôle (conditions, boucles)
- ✅ **Chapitre 3** — Tables (tableaux, dictionnaires, métatables)

### Module 2 — Programmation avancée
- ✅ **Chapitre 4** — Fonctions (paramètres, retours multiples, closures, variadiques)
- 🔜 Chapitre 5 — Gestion des chaînes
- 🔜 Chapitre 6 — Modules et packages

### Module 3 — Concepts avancés
- 🔜 Chapitre 7 — Métatables et POO
- 🔜 Chapitre 8 — Coroutines
- 🔜 Chapitre 9 — Gestion d'erreurs

### Module 4 — Préparation à FiveM
- 🔜 Chapitre 10 — Lua dans l'écosystème FiveM
- 🔜 Chapitre 11 — Natives et API FiveM
- 🔜 Chapitre 12 — Projet pratique

---

*Légende : ✅ disponible · 🔜 à venir.*
