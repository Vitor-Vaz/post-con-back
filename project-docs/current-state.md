# Estado atual

## Em andamento

API **post-con-back** (PostoConfiável): criação de review via HTTP atualiza agregados na tabela **`station`** (média e contagem com base nas últimas reviews por `place_id`). Postgres, migrações, **sqlc**, CI alinhados ao **mesa-mestre**. Próximos passos em `next-tasks.md`; visão de produto em `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`.

---

## Última sessão (consolidado)

### O que foi feito

- **Infra local:** `docker-compose.yaml` (Postgres 15, host **5433** → container 5432, DB `post_confiavel`), `.env.example`, `.gitignore` (`/vendor`, `.env`).
- **Banco:** `extension/database` (conexão por env); migrações golang-migrate (`uuid-ossp`, **`reviews`** MVP, **`station`** em `000003_station_table`: `place_id` único, `name` NOT NULL, `address`, coordenadas, `total_score`, `review_count`, `summary`, timestamps).
- **API:** Gin em `cmd/server` → `internal/app/router.go` (`/test`, `/health`, grupo `/api/v1`); **`POST /api/v1/review`** em `internal/app/v1/create_review.go` — body `place_id`, `user_id`, `rating`; trim de `place_id`; erros de domínio → HTTP.
- **Domínio:** `internal/domain/create_review.go` — `ReviewCreatorUseCase` após `InsertReview` chama `GetRecentReviewStats` (até 100 reviews) e **`UpsertStationScore`** na `station` (média e `review_count`). `CreateReviewInput` sem nome de posto na borda; persistência de `station.name` no insert usa o mesmo valor de `place_id` até existir fonte (ex.: Places).
- **Persistência:** `reviews_repository.go` (`InsertReview`, `GetRecentReviewStats`), `station_repository.go` (`UpsertStationScore`); queries em `queries/reviews/` e `queries/station/upsert_score.sql`; **`sqlcgen`** versionado (`make sqlc-gen`).
- **Testes:** `internal/app/v1/create_review_test.go` (handler + stub). `internal/domain/create_review_test.go` em pacote **`domain_test`** (evita ciclo de import) com Postgres real via **`extension/testhelpers`** (`SetupTestDB` usa `DATABASE_URL` do ambiente ou default local **5433** / DB `post_confiavel`).
- **Makefile:** `run`, `test`, **`lint`** (`fmt-check` + `go vet` + `go build`), `migrate-*`, `sqlc-gen`; **`DATABASE_URL ?= …`** (default **5433**; env sobrescreve).
- **CI:** `.github/workflows/merge.yaml` — job **test** com `env` **`DATABASE_URL`** no nível do job (Postgres serviço em **localhost:5432** no runner), `make migrate-up`, `make test`.

### Arquivos e pastas-chave

| Área | Caminhos |
|------|-----------|
| Plano | `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md` |
| Entry | `cmd/server/main.go` |
| HTTP | `internal/app/router.go`, `internal/app/v1/create_review.go`, `internal/app/v1/create_review_test.go` |
| Domínio | `internal/domain/errors.go`, `internal/domain/create_review.go`, `internal/domain/create_review_test.go` |
| Test DB | `extension/testhelpers/db.go` |
| Postgres | `internal/gateway/postgres/repositories/`, `sqlcgen/`, `queries/reviews/`, `queries/station/`, `sqlc.yaml` |
| Migrações | `extension/database/priv/migrations/` |
| CI | `.github/workflows/merge.yaml` |
| Referência | `../mesa-mestre` |

### Problemas conhecidos

- Binário **`migrate`** no **PATH** local (`$(go env GOPATH)/bin`) para `make migrate-up`.
- **`sqlc` CLI:** opcional para regenerar; versões novas podem exigir Go mais novo que o `go.mod` para `go install`.

### Pontos de atenção

- Tabela **`reviews`** e agregação em **`station`** são **MVP**; ponderação, anti-fraude e nome/endereço reais via **Places** ainda em roadmap (`next-tasks.md`, plano técnico).
- `station.name` na primeira linha é preenchido com **`place_id`** só para satisfazer NOT NULL; atualização de scores no `ON CONFLICT` **não altera** `name` já existente.
- **`reviews.user_id`** sem FK até existir **`users`**.
- **`vendor/`** no `.gitignore`; se usar `-mod=vendor`, rodar `go mod vendor` após mudar dependências.
