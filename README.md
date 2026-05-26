# 🚀 Project Template

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white) ![React](https://img.shields.io/badge/React-61DAFB?style=flat&logo=react&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white) ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white)

A **full-stack template** featuring:

- **Backend:** Go
- **Database:** PostgreSQL
- **Frontend:** TypeScript + React
- Includes **local development**, **type generation**, and **production-ready Docker workflows**.

---

## 📑 Table of Contents

- [Requirements](#requirements)
- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [Database](#database)
- [Backend](#backend)
- [Frontend](#frontend)
- [Generating API Types](#generating-api-types)
- [Running in Docker (Production)](#running-in-docker-production)
- [Makefile Commands](#makefile-commands)

---

## ⚙️ Requirements

- [Go](https://golang.org/) (backend)
- [Air](https://github.com/air-verse/air) (backend hot reloading)
- [Bun](https://bun.sh/) (frontend)
- [Docker](https://www.docker.com/) (database & production builds)
- [OpenAPI Generator](https://openapi-generator.tech/) (API types)

---

## 🛠 Setup

1. **Clone the repository**:
2. **Copy development environment example**:

```bash
cp api/.env.example.dev api/.env
```

3. **Edit `api/.env`**:

- `API_PORT` — backend port (default `35231`)
- `DB_CONNECTION_STRING` — database connection
- OAuth credentials — insert keys or dummy values for dev
- `IMAGES_PATH` — must be absolute if used
- `OAUTH_JWT_SECRET` & `AUTH_JWT_SECRET` — required for auth

4. **Copy frontend environment variables**:

```bash
cp fe/.env.example.dev fe/.env
```

> **Tip:** If you change `API_PORT`, update `VITE_API_BASE_PATH`:

```env
VITE_API_BASE_PATH=http://127.0.0.1:<API_PORT>/
```

---

## 🌐 Environment Variables

Common variables:

- `GO_ENV` — `development` or `production`
- `DB_CONNECTION_STRING` — PostgreSQL connection string
- `DB_MIGRATIONS_PATH` — path to migrations
- `API_PROD_URL` / `FE_PROD_URL` — production URLs
- OAuth & JWT secrets

For local production testing:

```bash
cp api/.env.example.production api/.env.production
cp fe/.env.example.production fe/.env.production
```

---

## 🗄 Database

Start PostgreSQL via Docker:

```bash
make db        # Start database container
make db-logs   # Follow logs
make db-stop   # Stop container
make db-remove # Remove container
```

**Defaults:**

- Name: `project-template-db`
- Port: `35232`
- Password: `s0m3C0mpl3xP4ss`

> Data persists in `./db-data`.

---

## 🖥 Backend

Install and run:

```bash
make api-install  # Install Go modules
make api          # Start backend (hot reload)
make api-update   # Update dependencies + format
```

- Default dev port: `35231`

---

## 💻 Frontend

Install and run:

```bash
make fe-install  # Install packages
make fe          # Start dev server
make fe-update   # Update packages + lint
```

- Default dev port: `35230` (Vite)

---

## 📦 Generating API Types

Backend exposes an OpenAPI spec. Generate TypeScript types:

```bash
make gen-types
```

- Output: `./shared/fe/api-client/src/generated`
- Automatically installs and builds the package

---

## 🐳 Running in Docker (Production)

Build and run production images:

```bash
make prod-build v1.0.0         # Build Docker image
make prod-locally v1.0.0       # Run locally
make prod-locally-logs          # Follow logs
make prod-locally-stop          # Stop container
```

- Default local production API port: `35230`
- Ensure database is running (`make db`)

---

## 📋 Makefile Commands

### Backend

- `make api` — run backend (hot reload)
- `make api-install` — install Go modules
- `make api-update` — update dependencies

### Frontend

- `make fe` — start frontend dev server
- `make fe-install` — install packages
- `make fe-update` — update + lint

### Database

- `make db` — start DB
- `make db-logs` — follow logs
- `make db-stop` — stop DB
- `make db-remove` — remove container

### API Types

- `make gen-types` — generate TS types from OpenAPI

### Production

- `make prod-build vX.X.X` — build Docker image
- `make prod-locally vX.X.X` — run locally
- `make prod-locally-logs` — follow logs
- `make prod-locally-stop` — stop container

### Combined

- `make install` — install frontend + backend + generate types
- `make update` — update frontend + backend + re-install
