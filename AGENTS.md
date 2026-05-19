# AGENTS.md — post-con-back (PostoConfiável)

Documento único para **assistentes de IA e desenvolvedores** alinharem comportamento neste repositório. Complementa o plano de produto em `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`. Atualize este arquivo ao fechar tarefas grandes ou decisões estáveis; histórico fino fica em commits e PRs.

---

## Propósito

API **Go** do produto **PostoConfiável**: avaliações de postos de combustível no Brasil, chave externa **`place_id`** (Google Places), roadmap com ponderação, anti-fraude e integração Maps (fases futuras).

---

## Quando usar este guia

Use como referência principal quando a tarefa envolver:

- HTTP (Gin), novos handlers em `internal/app/v1`
- Casos de uso e entidades em `internal/domain`
- Persistência Postgres, **sqlc**, repositórios em `internal/gateway/postgres`
- Migrações em `extension/database/priv/migrations`
- Makefile, CI, Docker local, variáveis de ambiente de banco

---

## Stack

| Área | Escolha |
|------|---------|
| Linguagem | Go 1.22 (`go.mod`) |
| HTTP | Gin |
| Banco | PostgreSQL 15 (`lib/pq`) |
| Config | `caarlos0/env`, `godotenv` |
| Migrações | golang-migrate |
| Queries tipadas | sqlc — código gerado em `internal/gateway/postgres/sqlcgen` (versionado; `make sqlc-gen`) |
| Testes | `testing`, `testify` |

---

## Estrutura de pastas (visão geral)

```text
cmd/server/              # main
internal/
  app/                   # Gin: router, grupos de rota
  app/v1/                # handlers HTTP v1 + testes
  domain/                # entidades, erros de domínio, casos de uso, interfaces de repositório
  gateway/postgres/
    repositories/        # implementações que chamam sqlc
    queries/             # SQL fonte do sqlc (arquivos `*.sql` na raiz, ex. `station_get_stations.sql`)
    sqlcgen/             # código gerado (commitar quando mudar queries)
    sqlc.yaml
extension/
  database/              # abertura de *sql.DB
  database/priv/migrations/
technical-plan/          # visão de produto e roadmap
.github/workflows/     # CI (Merge)
docker-compose.yaml    # Postgres local, porta host 5433 → 5432 no container
```

Repositório de referência de tooling: **`../mesa-mestre`** (Makefile, migrações, estilo de CI).

---

## Fluxo de camadas (obrigatório)

```text
HTTP (Gin handler em internal/app/v1)
  → internal/domain (caso de uso + contratos de repositório)
  → internal/gateway/postgres/repositories (sqlc)
  → PostgreSQL
```

- **Handlers:** binding/validação de entrada na borda (ex.: `binding` tags, trim de `place_id`), mapeamento de erros de domínio para status HTTP. Evitar duplicar no use case a mesma validação que o Gin já aplicou na borda (ver resumo **ADR-009** na seção de decisões abaixo).
- **Domain:** entidades, `var Err…` exportados, structs de input, interfaces que o repositório implementa, orquestração do caso de uso.
- **Repositories:** traduzem tipos de domínio ↔ parâmetros sqlc; mapeiam erros Postgres (`pq`) para erros de domínio quando fizer sentido.

**Proibido:** handler acessar `sql.DB` ou queries SQL diretamente; pular o caso de uso para persistir.

---

## HTTP (contrato atual)

| Método | Caminho | Descrição |
|--------|---------|-----------|
| GET | `/test` | smoke string |
| GET | `/health` | ping Postgres (timeout curto) |
| POST | `/api/v1/review` | corpo JSON: `place_id`, `user_id` (UUID), `rating` (1–5) |
| GET | `/api/v1/stations` | lista postos paginada; query `page` (default 1), 10 itens por página |

Novas rotas: registrar em `internal/app/router.go` (ou sub-rotas por versão em `internal/app/v1`).

---

## Banco e migrações

- URL: variável **`DATABASE_URL`**; no Makefile usa-se **`DATABASE_URL ?= …`** com default **localhost:5433**, base **`post_confiavel`**, para não sobrescrever env no CI (ADR-007).
- Migrações versionadas em `extension/database/priv/migrations/` (extensões, **`reviews`**, tabela **`station`** conforme evolução do schema). Nem toda tabela migrada precisa ter caso de uso exposto ainda; conferir código e router.
- **Não editar** migrações já aplicadas em ambientes compartilhados; criar nova sequência.
- Binário **`migrate`** precisa estar no `PATH` local (ex.: `$(go env GOPATH)/bin`).

---

## sqlc

- Fonte: `internal/gateway/postgres/queries/*.sql` (sqlc 1.31 não varre subpastas; usar `queries/` no `sqlc.yaml`) + `sqlc.yaml`.
- Após alterar SQL: rodar `make sqlc-gen` (requer CLI `sqlc` instalada) e commitar mudanças em `sqlcgen/`.

---

## Comandos (fonte de verdade: Makefile)

```bash
make install-linters  # golangci-lint e staticcheck (versões no Makefile)
make lint             # gofmt, go vet, go build, golangci-lint, staticcheck
make test
make run
make migrate-up | migrate-down | migrate-drop
make sqlc-gen
```

Subir Postgres local: `docker compose up -d` (porta **5433** no host).

---

## Testes

- Handlers: pacote `v1_test` com `httptest`, tabela de casos quando couber (`internal/app/v1/create_review_test.go`).
- Casos de uso com dependência real de Postgres: preferir pacote **`domain_test`** no mesmo diretório do domínio para **evitar ciclo de import** com `internal/gateway/postgres/repositories`; helper de conexão em `extension/testhelpers` quando existir; alinhar **`DATABASE_URL`** com o ambiente (local 5433, CI 5432).
- Nomes de teste e asserts em **inglês** (padrão do repositório).

---

## CI

Arquivo `.github/workflows/merge.yaml`:

- **lint:** `make lint`
- **test:** serviço Postgres 15, instala `migrate`, `make migrate-up` com `DATABASE_URL` apontando para **localhost:5432** no runner, depois `make test`

---

## Git e contribuição

- **Branches:** não trabalhar direto em `main`/`master`; usar padrão `feat/…`, `fix/…`, `chore/…` com escopo curto.
- **Commits:** mensagens em **inglês**, imperativo, estilo convencional (`feat:`, `fix:`, …), até ~72 caracteres na primeira linha.
- **PRs:** descrições claras; CI deve passar.

---

## Segurança

- Não commitar segredos; usar `.env` (local) e variáveis no CI.
- Não logar tokens, senhas ou dados sensíveis.

---

## Código (padrões do time)

- Código e testes em **inglês**.
- Evitar comentários óbvios; nomes claros.
- Diffs pequenos e focados na tarefa.
- Respeitar `fmt-check` (apenas arquivos Go rastreados pelo git).

---

## Produto e pendências (resumo)

- Tabela **`reviews`** é **MVP**; evolução prevista para ponderação, dimensões e contexto de veículo — ver `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`.
- **`place_id`** é a chave canônica do posto.
- **Fila típica (alta nível):** evoluir **`station`** e dados de posto (Places), revisar **`reviews`** para ponderação, anti-fraude + sinais Maps/Places, **`users`** + auth + FK em `reviews`, **GET** reviews por `place_id` (+ paginação), **OpenAPI**, cliente Places em `internal/integration/…`, **`place_scores`** / recência após regras de negócio, health mais rico se precisar para deploy.

---

## Decisões de arquitetura (resumo)

- **Layout `internal/`:** `domain` (entidades, erros, caso de uso, interfaces de repo); `app` / `app/v1` (Gin); `gateway/postgres` (repos, sqlc, queries). Alinhado ao **mesa-mestre**.
- **`reviews` MVP:** `place_id`, `user_id` UUID NOT NULL, `rating` 1–5; sem FK em `users` até existir a tabela; schema evolui por novas migrações.
- **Postgres local:** host **5433** no `docker-compose`; **`DATABASE_URL ?=`** no Makefile para o CI poder sobrescrever (**5432** no GitHub Actions).
- **sqlc:** `sqlcgen/` versionado; após mudar SQL, `make sqlc-gen` e commitar o diff.
- **`fmt-check`:** `gofmt -l` só em `git ls-files '*.go'`.
- **ADR-009:** validação de formato/range na **borda HTTP** (Gin); use case não revalida o que o binding já garantiu (novos callers diretos no use case devem validar ou ter outra borda).

---

## Onde aprofundar

| Conteúdo | Arquivo |
|----------|---------|
| Visão e roadmap de produto | `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md` |
| Fluxo de trabalho e ponte com mesa-mestre | `README.md` (raiz) |

**Ao encerrar uma sessão:** atualizar este `AGENTS.md` e/ou o plano técnico quando escopo ou decisões mudarem; manter commits e PRs descritivos.

---

## Persona sugerida para o agente

- Desenvolvedor Go com Gin e Postgres.
- Respeita camadas acima; não inventa camada paralela (ex.: “service” genérico fora de `domain`) sem alinhar ao repositório.
- Prefere inspecionar código existente antes de refatorar em larga escala.
- Comunica impacto (arquivos, riscos, migrações) de forma objetiva.

---

## Prompts de exemplo

- "Adiciona GET `/api/v1/reviews` por `place_id` com paginação, sqlc e testes de handler."
- "Cria migration para coluna X em `reviews` sem alterar migrations antigas."
- "Ajusta o mapeamento de erro Postgres 235xx para o erro de domínio correto."
- "Atualiza o workflow de CI para também rodar X."
