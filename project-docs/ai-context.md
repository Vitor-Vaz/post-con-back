# AI Context

## Projeto

**PostoConfiável** — API Go em **post-con-back**. Plano: `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`. Referência de pastas/tooling: **mesa-mestre** (`../mesa-mestre`).

## Stack

Go 1.22, **Gin**, **PostgreSQL** (`lib/pq`), **env** + **godotenv**, **golang-migrate**, **sqlc** (`internal/gateway/postgres/sqlcgen` + `queries/`, `make sqlc-gen`). Docker: `docker-compose.yaml`, Postgres na host **5433**.

## Arquitetura (`internal/`)

| Pacote | Papel |
|--------|--------|
| `internal/domain` | Entidades, erros de domínio, caso de uso + interface do repo (ex. `ReviewCreatorUseCase`) |
| `internal/app` | `router.go` — Gin, `/test`, `/health`, `/api/v1` |
| `internal/app/v1` | Handlers v1 |
| `internal/gateway/postgres` | `repositories/`, `sqlcgen/`, `queries/<tabela>/<ação>.sql` |
| `extension/database` | `*sql.DB` |
| `cmd/server` | Entry |

Fluxo: **HTTP → domain (use case) → gateway (repo)**.

## HTTP útil

- `GET /test`, `GET /health` (ping DB)
- `POST /api/v1/review` — body: `place_id` (string), `user_id` (UUID JSON), `rating` (1–5); Gin `binding` + trim de `place_id` no handler; erros de domínio mapeados (400/404/409/500)

## Makefile

`run`, `test`, `lint` (gofmt em arquivos `git ls-files '*.go'`, `vet`, `build`), `sqlc-gen`, `migrate-*`. **`DATABASE_URL ?= ...`** default porta **5433**; **env sobrescreve** (ex.: CI na 5432).

## CI

`.github/workflows/merge.yaml` — em PR/push `main|master`: job **lint** (`make lint`), job **test** (Postgres serviço, `migrate`, `make migrate-up` com `DATABASE_URL` na 5432, `make test`).

## Pendências de produto (alto nível)

`reviews` e schema **evoluem** com **ponderação** e **anti-fraude** (parte com **Maps/Places**); ver `next-tasks.md` e plano técnico.

## Regras

Contratos e migrações já aplicadas: mudanças pequenas e alinhadas. Código sem comentários desnecessários; testes e commits em inglês.

## Responder

Impacto → arquivos → patch pequeno → riscos.
