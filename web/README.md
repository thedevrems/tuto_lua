# Lua Academy — plateforme de cours Lua / FiveM

Plateforme web interactive pour apprendre Lua : chaque chapitre contient un **cours**
(le sujet) et des **exercices** avec un éditeur de code, une **console**, des **tests
automatiques** et une **réponse** (solution).

Le code Lua de l'élève s'exécute **directement dans le navigateur** (Lua 5.4 via
[wasmoon](https://github.com/ceifa/wasmoon), WebAssembly) — aucun serveur n'est nécessaire.

## Stack

- **React 18 + TypeScript + Vite**
- **Tailwind CSS** — thème noir & blanc moderne
- **wasmoon** — VM Lua 5.4 en WebAssembly (exécution client)
- **CodeMirror 6** — éditeur de code (coloration Lua)
- **react-markdown** — rendu du contenu des cours

## Démarrer

```bash
cd web
npm install
npm run dev      # http://localhost:5173
```

Autres commandes :

```bash
npm run build    # build de production (tsc --noEmit + vite build) -> dist/
npm run preview  # sert le build de production en local
```

Le site est entièrement **statique** : `dist/` se déploie sur GitHub Pages, Netlify,
Vercel, etc.

## Organisation du contenu

Tout le programme vit dans [`src/content/curriculum.ts`](src/content/curriculum.ts),
structuré en **Module → Chapitre → Items**. Un item est soit une **leçon**, soit un
**exercice**.

```
Module 1 — Fondamentaux
 └─ Chapitre 1 — Introduction
     ├─ Cours        (leçon : markdown)
     ├─ Exercice 1   (éditeur + console + tests + solution)
     ├─ Exercice 2
     └─ ...
```

Les textes de cours volumineux sont des fichiers Markdown dans
[`src/content/lessons/`](src/content/lessons/), importés en `?raw`.

### Ajouter une leçon

1. Créez `src/content/lessons/mon-cours.md`.
2. Dans `curriculum.ts` :

```ts
import monCours from './lessons/mon-cours.md?raw'

// ... dans items: [
{ kind: 'lesson', id: 'm1c3-cours', title: 'Cours', content: monCours },
```

### Ajouter un exercice

```ts
{
  kind: 'exercise',
  id: 'm1c3-ex1',                 // identifiant unique (sert à sauvegarder le code)
  title: 'Exercice 1 — Mon titre',
  difficulty: 'facile',            // 'facile' | 'moyen' | 'difficile'
  statement: 'Énoncé en **markdown**…',
  starter: 'local x = ',           // code pré-rempli dans l'éditeur
  solution: 'local x = 10\nprint(x)',
  tests: [
    { name: 'x vaut 10', code: 'assert(printed("10"), "10 attendu")' },
  ],
  hints: ['Un indice optionnel'],  // optionnel, révélé un par un
}
```

### Écrire des tests

Un test est un **chunk Lua** exécuté dans le même environnement, juste après le code
de l'élève. Il doit `assert(condition, "message")`. S'il lève une erreur → test échoué.

Helpers disponibles dans les tests :

| Helper | Description |
| --- | --- |
| `printed(sub)` | `true` si une ligne affichée contient `sub` |
| `output_text()` | toute la sortie, lignes jointes par `\n` |
| `output_lines()` | table des lignes affichées |

Les **fonctions et variables globales** définies par l'élève (ex. `function Foo()…end`)
sont visibles dans les tests, qui peuvent donc les appeler :

```ts
tests: [
  { name: 'Foo(2) == 4', code: 'assert(Foo(2) == 4)' },
]
```

> ⚠️ Les variables **locales** (`local`) de l'élève ne sont pas visibles dans les
> tests. Pour les exercices basés sur des valeurs, testez la **sortie** avec
> `printed(...)` ; pour tester de la logique, demandez une **fonction globale**.

Les tests sont **optionnels** : un exercice sans `tests` n'affiche que la console.

## Notes techniques

- **Garde anti-boucle infinie** : le VM s'interrompt après ~8 M d'instructions
  (`debug.sethook`), voir [`src/lib/lua.ts`](src/lib/lua.ts).
- **Progression & code** : le code saisi et les exercices réussis sont sauvegardés
  dans le `localStorage` du navigateur.
- Le `.wasm` de wasmoon est importé via `wasmoon/dist/glue.wasm?url` pour un chargement
  fiable côté navigateur.
