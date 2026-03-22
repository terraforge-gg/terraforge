# terraforge

A Terraria mod distribution platform. Browse, publish, and manage mods at [terraforge.gg](https://terraforge.gg).

## Getting Started

Install prerequisites:
- [Go](https://go.dev/dl/)
- [Node.js](https://nodejs.org/) (LTS)
- [pnpm](https://pnpm.io/installation)
- [Docker](https://www.docker.com/) & Docker Compose
- [Task](https://taskfile.dev/installation/) (`taskfile.yml` runner)
- [AWS CLI](https://aws.amazon.com/cli/) (for LocalStack interactions)

---

### 1. Start infrastructure

Spin up Postgres, Meilisearch, LocalStack, and Redis with Docker Compose:

```bash
docker compose up -d
```

### 2. Setup environment variables
Copy `.env.example` values into a `.env` for each application

### 3. Run the applications


Each app can be started independently using [Task](https://taskfile.dev/):

*Note: All migrations are located under `./apps/api/internal/database/migrations`. These will be automatically applied on api start through [Goose](https://github.com/pressly/goose)*

```bash
task api:start

task auth:start

task web:start
```

To seed the database with a user run:

```bash
task auth:seed
```
