# Estado atual

## Em andamento

Evolução da API **post-con-back** (PostoConfiável): review criável via HTTP, Postgres com migrações, CI no padrão do **mesa-mestre**. Próximos passos em `next-tasks.md`; visão de produto e roadmap em `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`.

---

## Última sessão (consolidado)

### O que foi feito

- **Infra local:** `docker-compose.yaml` (Postgres 15, host **5433** → container 5432, DB `post_confiavel`), `.env.example`, `.gitignore` (`/vendor`, `.env`).
- **Banco:** `extension/database` (conexão por env); migrações golang-migrate (`uuid-ossp`, tabela **`reviews`** MVP: `place_id`, `user_id` NOT NULL, `rating` double com CHECK 1–5, timestamps, índice em `place_id`).
- **API:** Gin em `cmd/server` → `internal/app/router.go` (`/test`, `/health`, grupo `/api/v1`); **`POST /api/v1/review`** em `internal/app/v1/create_review.go` com `binding` Gin, trim de `place_id`, mapeamento de erros de domínio para HTTP.
- **Domínio:** `internal/domain` — erros genéricos (`ErrBadParams`, `ErrBadRequest`, `ErrNotFound`, `ErrConflict`, `ErrUnexpected`); `create_review.go` com entidade `Review`, input, interface do repo e **`ReviewCreatorUseCase`** (sem revalidar o que o Gin já valida na borda).
- **Persistência:** `internal/gateway/postgres/repositories/reviews_repository.go` + **`sqlcgen`** alinhado a `queries/reviews/insert.sql`; `make sqlc-gen` aponta para `internal/gateway/postgres/sqlc.yaml`.
- **Testes:** `internal/app/v1/create_review_test.go` — casos em **uma tabela** + `httptest.NewServer`, `testify`, cobrindo sucesso, binding, erros do use case.
- **Makefile:** `run`, `test`, **`lint`** (`fmt-check` com `gofmt` em `git ls-files '*.go'`, `go vet`, `go build`), `migrate-*`, `sqlc-gen`; **`DATABASE_URL ?= …`** (default porta **5433**; env do CI sobrescreve para **5432**).
- **CI:** `.github/workflows/merge.yaml` (nome **Merge**, como no mesa) — em `push`/`pull_request` para `main` e `master`: job **lint** (`make lint`), job **test** (serviço Postgres, instala `migrate`, `make migrate-up` com `DATABASE_URL` do workflow, `make test`).

### Arquivos e pastas-chave

| Área | Caminhos |
|------|-----------|
| Plano | `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md` |
| Entry | `cmd/server/main.go` |
| HTTP | `internal/app/router.go`, `internal/app/v1/create_review.go`, `internal/app/v1/create_review_test.go` |
| Domínio | `internal/domain/errors.go`, `internal/domain/create_review.go` |
| Postgres | `internal/gateway/postgres/repositories/`, `sqlcgen/`, `queries/reviews/insert.sql`, `sqlc.yaml` |
| Migrações | `extension/database/priv/migrations/` |
| CI | `.github/workflows/merge.yaml` |
| Referência vizinha | `../mesa-mestre` (Makefile, estilo de teste/CI) |

### Problemas conhecidos

- O binário **`migrate`** precisa estar no **PATH** local (ex.: `$(go env GOPATH)/bin`) para `make migrate-up`, como no mesa-mestre.
- **`sqlc` CLI:** opcional para regenerar; versões novas podem exigir Go mais novo que o `go.mod` para `go install`.

### Pontos de atenção

- Tabela **`reviews`** é **provisória**: será revista para **ponderação** e metadados; regra de “justiça” com o posto **ainda não definida**.
- **Anti-fraude** e uso de **Google Maps / Places** estão no roadmap e no plano técnico; implementação futura com termos, cotas e LGPD.
- **`reviews.user_id`** sem FK até existir **`users`**; integridade depende da aplicação.
- **`vendor/`** no `.gitignore`; se o time usar `-mod=vendor`, rodar `go mod vendor` após mudar dependências.
