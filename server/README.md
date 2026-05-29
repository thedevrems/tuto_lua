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
| GET     | `/api/health`        | public   | Sonde de disponibilité               |
| POST    | `/api/auth/register` | public   | Création de compte (+ token)         |
| POST    | `/api/auth/login`    | public   | Connexion par e-mail ou username     |
| GET     | `/api/auth/me`       | connecté | Profil de l'utilisateur authentifié  |

D'autres routes (cours, progression, paiement, admin) arrivent dans les phases suivantes.
