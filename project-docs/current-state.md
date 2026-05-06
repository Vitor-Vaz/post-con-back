# Estado atual

## Em andamento

Manutenção e evolução da API **post-con-back** conforme plano técnico PostoConfiável; próximos passos sugeridos em `next-tasks.md`.

## Última sessão (consolidado)

### O que foi feito

- Setup inicial Go com servidor local e rota de teste; evolução para Postgres via Docker (`docker-compose.yaml`, porta host **5433**).
- `extension/database` com conexão por env; `Makefile` com `run`, `migrate-up` / `migrate-down` / `migrate-drop` / `migrate-create`, `sqlc-gen`.
- Migrações: extensão `uuid-ossp`; tabela **`reviews`** simplificada (`id`, `place_id`, `user_id` NOT NULL, `rating` double precision com CHECK 1–5, timestamps, índice em `place_id`).
- API **POST `/api/v1/review`** para criar review (Gin).
- Organização **mesa-mestre + `internal/`**: `internal/domain` (inclui caso de uso `ReviewCreatorUseCase`), `internal/app` + `internal/app/v1`, `internal/gateway/postgres` (repositório `ReviewsRepository` + sqlc).
- Ajuste de `Makefile` para usar o binário `migrate` no PATH (como no mesa-mestre); orientação para `GOPATH/bin` quando `migrate` não é encontrado.
- `go mod vendor` utilizado no fluxo de build do projeto quando vendor está presente.

### Arquivos e pastas-chave

- `technical-plan/SAAS-POSTOS-COMBUSTIVEL.md` — plano de produto e arquitetura alvo
- `cmd/server/main.go` — sobe Gin via `internal/app.NewRouter`
- `internal/app/router.go` — rotas globais + `/api/v1`
- `internal/app/v1/create_review.go` — handler POST review
- `internal/domain/` — `errors.go`, `create_review.go` (entidade review, input e caso de uso); validação de entrada na camada HTTP (Gin `binding`)
- `internal/gateway/postgres/repositories/reviews_repository.go` — persistência
- `internal/gateway/postgres/sqlcgen/` — queries compiladas manualmente alinhadas a `queries/reviews/insert.sql` (regenerar com `make sqlc-gen` quando `sqlc` estiver instalado)
- `extension/database/priv/migrations/` — migrações golang-migrate
- `.env.example` — variáveis do banco; `.env` ignorado no git

### Problemas conhecidos

- CLI **`migrate`** precisa estar no `PATH` (ex.: `$(go env GOPATH)/bin`) para `make migrate-up` funcionar como no mesa-mestre.
- **`sqlc`**: se regenerar código, comparar diff com `sqlcgen` versionado; versão recente do `sqlc` pode exigir Go mais novo que a linha do `go.mod` para `go install`.

### Pontos de atenção

- Tabela **`reviews` será revista** para suportar **ponderação** de scores e possivelmente mais metadados; o desenho atual é temporário até fechar a fórmula de “justiça” com o posto.
- **Ponderação** (pesos, recência, dimensões): regra de negócio **ainda não definida**; impactará schema, agregados e API.
- **Anti-fraude / reviews falsos**: será necessária **lógica própria** no backend; parte do desenho prevê uso de **Google Maps / Places** como sinal (ex.: consistência de lugar), sempre dentro de termos e cotas — ver roadmap no plano técnico e `next-tasks.md`.
- Tabela `reviews` ainda **sem FK** para `users` (tabela de usuários ainda não existe); `user_id` é obrigatório na aplicação.
- Plano futuro adicional: dimensões de avaliação, perfil de veículo, `place_scores`, integração Places mais ampla — parcialmente coberto pelo roadmap estendido.
