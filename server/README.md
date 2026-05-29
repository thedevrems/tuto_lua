# Lua Academy — Backend (Go)

API REST en **Go** (net/http + [chi](https://github.com/go-chi/chi)) avec base **SQLite**
(pure-Go, sans CGO). Elle gère les comptes, l'authentification, la progression,
les cours stockés en base, les paiements Stripe et l'administration.

## 🚀 Démarrage

```bash
cd server
cp .env.example .env        # ajustez si besoin
go run ./cmd/api            # http://localhost:8080
```

> Le **premier compte créé devient automatiquement administrateur** (bootstrap).

## 🧪 Tests

```bash
go test ./...               # toute la suite
go test -cover ./...        # avec couverture
```

## 🔐 Sécurité

- Mots de passe **hachés** avec **bcrypt** (jamais stockés ni renvoyés en clair).
- Transport chiffré par **HTTPS/TLS** (à terminer côté reverse-proxy en prod).
- Sessions **stateless** via **JWT** signés HMAC-SHA256.

## 📂 Architecture

Règles de code : une fonction = une responsabilité, ≤ 60 lignes/fonction,
≤ 300 lignes/fichier, zéro duplication.

```
cmd/api/main.go        Point d'entrée : config → DB → router → serveur
internal/
  config/              Chargement de la configuration (env)
  database/            Connexion SQLite + migrations idempotentes
  models/              Structures de données métier
  crypto/              Hachage bcrypt des mots de passe
  token/               Émission / vérification des JWT
  validate/            Règles de validation (username, email, mot de passe)
  httpx/               Helpers JSON (réponses, erreurs, décodage)
  store/               Couche d'accès aux données (toutes les requêtes SQL)
  auth/                Service register/login + middlewares (auth, admin)
  handlers/            Endpoints HTTP
  router/              Montage des routes + CORS
```

## 🌐 Endpoints (état actuel)

| Méthode | Route                | Accès    | Description                          |
| ------- | -------------------- | -------- | ------------------------------------ |
| GET     | `/api/health`         | public   | Sonde de disponibilité               |
| POST    | `/api/auth/register`  | public   | Création de compte (+ token)         |
| POST    | `/api/auth/login`     | public   | Connexion par e-mail ou username     |
| GET     | `/api/auth/me`        | connecté | Profil de l'utilisateur authentifié  |
| GET     | `/api/courses`        | public   | Catalogue des cours publiés          |
| GET     | `/api/courses/{slug}` | public\* | Arbre complet (gratuit, ou inscrit/admin pour les payants) |
| GET     | `/api/progress`       | connecté | Progression de l'utilisateur         |
| PUT     | `/api/progress/{exerciseId}` | connecté | Sauvegarde code + complétion  |
| GET     | `/api/enrollments`    | connecté | Cours auxquels l'utilisateur a accès |
| GET     | `/api/admin/users`    | admin    | Liste des comptes                    |
| GET     | `/api/admin/users/{userId}/progress` | admin | Dernier code (« push ») d'un élève |
| POST    | `/api/admin/enrollments` | admin | Donner l'accès à un cours            |
| GET     | `/api/admin/courses`  | admin    | Tous les cours (brouillons inclus)   |
| POST    | `/api/admin/courses`  | admin    | Créer un cours                       |
| POST    | `/api/admin/courses/{courseId}/chapters`   | admin | Créer un chapitre       |
| POST    | `/api/admin/chapters/{chapterId}/lessons`  | admin | Créer une leçon         |
| POST    | `/api/admin/chapters/{chapterId}/exercises`| admin | Créer un exercice       |
| POST    | `/api/admin/exercises/{exerciseId}/tests`  | admin | Créer un test           |

\* `/api/courses/{slug}` accepte un token optionnel : les cours gratuits sont
ouverts à tous, les cours payants exigent une inscription (achat) ou un rôle admin.

## 📚 Contenu (cours, exercices, tests)

Les cours, chapitres, leçons, exercices, **tests** et indices sont stockés en base.
Au **premier démarrage** (base vide), le catalogue est **semé** depuis
`internal/seed/curriculum.json` — fichier généré à partir du curriculum frontend :

```bash
cd web && node scripts/export-curriculum.mjs   # régénère internal/seed/curriculum.json
```

Chaque **module** devient un **cours** achetable (Module 1 gratuit, les autres payants).

L'**administrateur** (premier compte créé) peut, via `/admin` côté frontend :
créer des cours / chapitres / leçons / exercices / tests, donner l'accès à un
cours, et consulter le **dernier code poussé** par chaque élève.

> Le **paiement Stripe** (déblocage automatique après achat) est la dernière phase ;
> il nécessite des clés Stripe (`STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`).
