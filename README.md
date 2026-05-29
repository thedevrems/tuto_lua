# Lua Academy — Cours Lua interactif pour FiveM

Plateforme web pour apprendre **Lua** dans le contexte de **FiveM** : chaque chapitre
contient un **cours**, et des **exercices** avec un éditeur de code, une **console**, des
**tests automatiques** et une **solution**.

Le code Lua de l'élève s'exécute **directement dans le navigateur** (Lua 5.4 via
WebAssembly) — aucun serveur n'est nécessaire.

## ✨ Fonctionnalités

- 📚 Cours organisés en **Modules → Chapitres → Cours / Exercices**
- 💻 Éditeur de code intégré (coloration syntaxique Lua)
- ▶️ **Console** : exécution du code et sortie de `print()`
- ✅ **Tests automatiques** par exercice (✓ / ✗, compteur de réussite)
- 💡 Indices progressifs et **solution** dépliable
- 💾 Progression et code sauvegardés localement (navigateur)
- 🎨 Interface **noir & blanc** moderne

## 🧰 Stack technique

| Domaine | Choix |
| --- | --- |
| Front | React 18 + TypeScript + Vite |
| Style | Tailwind CSS (thème monochrome) |
| Exécution Lua | [wasmoon](https://github.com/ceifa/wasmoon) — Lua 5.4 en WebAssembly (client) |
| Éditeur | CodeMirror 6 |
| Contenu | Markdown (`react-markdown`) + données TypeScript |

## 🚀 Démarrage rapide

```bash
cd web
npm install
npm run dev      # http://localhost:5173
```

Autres commandes :

```bash
npm run build    # build de production -> web/dist/
npm run preview  # sert le build de production
```

Le site est **statique** : `web/dist/` se déploie sur GitHub Pages, Netlify, Vercel, etc.

## 📂 Structure du projet

```
.
├── web/                       # Application (plateforme interactive)
│   ├── src/
│   │   ├── content/           # Le programme du cours
│   │   │   ├── curriculum.ts  #   structure Modules/Chapitres/Exercices
│   │   │   └── lessons/       #   textes de cours (.md)
│   │   ├── components/        # UI (sidebar, éditeur, console, tests…)
│   │   └── lib/lua.ts         # Exécution Lua + tests (wasmoon)
│   └── README.md              # Guide : lancer & ajouter du contenu
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
