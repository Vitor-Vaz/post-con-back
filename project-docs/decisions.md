# Decisões relevantes

Registro curto de decisões que afetam código, contrato ou operação.

## ADR-001 — Tudo sob `internal/` com layout inspirado no mesa-mestre

- **Contexto**: alinhar com projeto de referência (mesa-mestre) mantendo pacotes privados sob `internal/`.
- **Decisão**: `internal/domain` concentra entidades e casos de uso + interfaces de repositório; `internal/app` e `internal/app/v1` para HTTP; `internal/gateway/postgres` para Postgres + sqlc.
- **Consequência**: imports sempre prefixados com `post-con-back/internal/...`.

## ADR-002 — Modelo inicial de `reviews` enxuto

- **Contexto**: MVP antes de fechar dimensões e perfil de veículo.
- **Decisão**: colunas `place_id`, `user_id` (UUID NOT NULL), `rating` (double, CHECK 1–5), timestamps; sem FK para `users` até existir a tabela.
- **Consequência**: integridade de `user_id` depende da aplicação até haver FK.

## ADR-003 — Postgres local na porta 5433

- **Contexto**: evitar conflito com outro Postgres local (ex.: outro projeto na 5432).
- **Decisão**: `docker-compose.yaml` publica `5433:5432`; `DATABASE_URL` / defaults do código alinhados a isso.

## ADR-004 — sqlcgen versionado + `make sqlc-gen`

- **Contexto**: permitir build sem depender do binário `sqlc` na máquina de cada dev.
- **Decisão**: manter pacote `internal/gateway/postgres/sqlcgen` no repositório, gerado a partir de `queries/*.sql` e `sqlc.yaml`; Makefile expõe `sqlc-gen`.
- **Consequência**: após mudar queries ou schema de referência do sqlc, rodar `make sqlc-gen` e commitar o diff quando aplicável.

## ADR-005 — Modelo `reviews` e ponderação em aberto

- **Contexto**: a tabela atual é MVP; a justiça do score para o posto exige **ponderação** ainda não fechada.
- **Decisão**: tratar o schema de `reviews` como **provisório**; evoluir com migrações quando a regra de ponderação e os dados necessários estiverem definidos.
- **Consequência**: possíveis migrações futuras e ajustes em sqlc/use cases; manter compatibilidade com ambientes que já aplicaram migrações anteriores.

## ADR-006 — Anti-fraude com apoio de Google Maps / Places

- **Contexto**: mitigar reviews falsos; o plano já cita risco de fraude.
- **Decisão**: reservar espaço no roadmap para regras de confiança e **uso seletivo** de dados/serviços do Google Maps / Places (ex.: validação contextual), sem violar termos, cotas nem LGPD.
- **Consequência**: novos use cases, integração e possivelmente armazenamento mínimo/cache conforme política; detalhes a especificar antes de implementar.
