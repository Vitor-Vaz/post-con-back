# AI Context

## Projeto

**PostoConfiável** (repo **post-con-back**): API backend em Go para SaaS de avaliações de postos de combustível no Brasil. Visão de produto, modelo de dados e roadmap estão em `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`.

## Repositório vizinho (mesa-mestre)

Este repo costuma ficar **lado a lado** com **mesa-mestre** (mesmo pasta pai), usado como referência de estrutura e tooling.

- Relativo: `../mesa-mestre`
- Exemplo: `.../projetos-pessoais/mesa-mestre` junto de `.../projetos-pessoais/post-con-back`

## Pendências de produto e domínio (não fechadas)

- **Modelo `reviews`**: será **revisto**; a tabela atual é MVP. Serão necessários campos ou tabelas de apoio para a **lógica de ponderação** dos scores (critérios de equidade com o posto **ainda a definir**).
- **Ponderação**: fórmula e pesos serão definidos depois, com foco no que for **mais justo** na relação com o posto (alinhado ao plano técnico de score composto e recência).
- **Anti-fraude / reviews falsos**: haverá regras dedicadas; parte da estratégia usará **Google Maps / Places** (validação contextual, sinais de confiança, limites de API e compliance).

Detalhes e fases também constam em `next-tasks.md` e no roadmap em `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md`.

## Stack

- Go (módulo `post-con-back`, `go 1.22` no `go.mod`)
- HTTP: **Gin**
- Banco: **PostgreSQL** (driver `lib/pq`)
- Config: `caarlos0/env` + `joho/godotenv`
- Migrações: **golang-migrate** (CLI `migrate`; alvos no `Makefile`)
- Persistência gerada: **sqlc** (código em `internal/gateway/postgres/sqlcgen/`; queries em `internal/gateway/postgres/queries/`; `make sqlc-gen`)
- Docker local: `docker-compose.yaml` (Postgres 15, host **5433** → container 5432, DB `post_confiavel`)

## Arquitetura (padrão alinhado ao mesa-mestre, tudo sob `internal/`)

- **`internal/domain`**: entidades, erros de domínio, validação, **casos de uso** e **interfaces de repositório** (ex.: `ReviewCreatorUseCase`, `ReviewCreatorRepository`).
- **`internal/app`**: composição HTTP (`router.go` com Gin), health e teste.
- **`internal/app/v1`**: handlers da API v1 (ex.: criar review).
- **`internal/gateway/postgres`**: `repositories/` (implementações), `sqlcgen/`, `queries/`, `sqlc.yaml`.
- **`extension/database`**: conexão `*sql.DB` (env `DATABASE_*`).
- **`cmd/server`**: entrypoint (`go run ./cmd/server` ou `make run`).

Fluxo típico: **HTTP (app/v1) → caso de uso (domain) → repositório (gateway/postgres)**.

## Contratos HTTP relevantes

- `GET /test` → `200`, corpo `ok`
- `GET /health` → `200` se Postgres responder ping; senão `503`
- `POST /api/v1/review` → JSON `place_id`, `user_id` (UUID string), `rating` (float 1–5); `201` com review criada ou `4xx`/`5xx` conforme erro

## Objetivo atual

Evoluir a API de reviews e integrações (usuários, Places, agregados) de forma incremental, mantendo migrações e contratos estáveis.

## Regras importantes

- não alterar contratos públicos da API sem alinhamento
- preservar compatibilidade com **migrations já aplicadas** em ambientes existentes
- preferir mudanças pequenas e incrementais
- evitar refactors amplos sem necessidade
- seguir padrões do repositório: sem comentários no código (código autoexplicativo), testes em inglês, commits semânticos em inglês
- `vendor/` está no `.gitignore`; após mudanças de dependências, rodar `go mod vendor` se o time usar build com vendor

## Como responder

Ao propor mudanças:

1. explicar impacto
2. apontar arquivos
3. sugerir patch pequeno
4. mencionar riscos
