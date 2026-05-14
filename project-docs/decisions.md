# Decisões relevantes

Registro de decisões que afetam código, contrato ou operação. Formato breve: **contexto → decisão → consequência**.

---

## ADR-001 — Layout sob `internal/` alinhado ao mesa-mestre

- **Contexto:** manter referência ao projeto **mesa-mestre** (`../mesa-mestre`) e pacotes privados em Go.
- **Decisão:** `internal/domain` (entidades, erros, caso de uso + interface de repositório); `internal/app` e `internal/app/v1` (HTTP Gin); `internal/gateway/postgres` (repositórios, sqlc, queries).
- **Consequência:** imports `post-con-back/internal/...`; espelhar Makefile/estilo de teste no mesa quando fizer sentido.

---

## ADR-002 — Modelo inicial da tabela `reviews` (MVP)

- **Contexto:** começar cedo com persistência sem fechar dimensões, veículo e ponderação.
- **Decisão:** colunas `place_id` (texto, Google), `user_id` (UUID NOT NULL), `rating` (double, CHECK 1–5), `created_at` / `updated_at`; sem FK para `users` até a tabela existir.
- **Consequência:** evolução por novas migrações; `user_id` confiável só na aplicação até haver FK.

---

## ADR-003 — Postgres local na porta 5433

- **Contexto:** evitar conflito com outro Postgres na máquina (ex.: 5432 ocupado por outro projeto).
- **Decisão:** `docker-compose.yaml` expõe **5433** no host; default de desenvolvimento documentado nessa porta.
- **Consequência:** `.env` / defaults locais apontam para `5433`; CI usa serviço na **5432** (padrão do container no GitHub Actions).

---

## ADR-004 — sqlcgen versionado e `make sqlc-gen`

- **Contexto:** time pode não ter o binário `sqlc` instalado em toda máquina.
- **Decisão:** manter `internal/gateway/postgres/sqlcgen` versionado, gerado a partir de `queries/**/*.sql` e `sqlc.yaml`; expor `make sqlc-gen`.
- **Consequência:** após mudar SQL de queries ou schema de referência do sqlc, regenerar e commitar o diff quando aplicável.

---

## ADR-005 — `reviews` e ponderação em aberto

- **Contexto:** score justo para o posto exige ponderação (recência, dimensões, etc.) ainda não fechada com produto.
- **Decisão:** tratar o schema atual de `reviews` como **temporário**; evoluir com migrações quando a regra e os dados necessários estiverem definidos.
- **Consequência:** possíveis alterações em domínio, sqlc e API; preservar compatibilidade com ambientes que já aplicaram migrações antigas (cuidado ao alterar `000002` em produção).

---

## ADR-006 — Anti-fraude e Google Maps / Places

- **Contexto:** mitigar reviews falsos; plano de risco já cita o tema.
- **Decisão:** reservar espaço no roadmap para regras de confiança e uso **complementar** de Maps/Places (ex.: validação contextual de `place_id`), sem violar termos, cotas nem LGPD.
- **Consequência:** novos use cases e integração em `internal/integration/...` quando especificado; não acoplar tudo ao Maps de uma vez.

---

## ADR-007 — `DATABASE_URL` com `?=` e CI

- **Contexto:** no Make, `DATABASE_URL=...` no Makefile **sobrescrevia** a variável de ambiente do GitHub Actions; `migrate` tentava `5433` no CI e falhava.
- **Decisão:** usar **`DATABASE_URL ?= ...`** no Makefile (default local `…:5433/…`); workflow define **`DATABASE_URL`** para `localhost:5432` no job de testes.
- **Consequência:** um único `make migrate-up` funciona em dev e no CI; quem rodar migrate localmente sem `.env` ainda usa o default da 5433.

---

## ADR-008 — `fmt-check` com `git ls-files`

- **Contexto:** `find . -name '*.go'` podia incluir ruído ou divergir do que o CI realmente versiona; `fmt-check` falhava em `errors.go` sem mudança aparente no diff do autor.
- **Decisão:** `fmt-check` roda `gofmt -l` apenas em **`git ls-files '*.go'`** (arquivos rastreados).
- **Consequência:** CI e máquina local alinhados ao mesmo conjunto de arquivos; arquivos não rastreados não entram no check até serem `git add`.

---

## ADR-009 — Validação na borda HTTP, não no use case de criação

- **Contexto:** evitar duplicar regras entre Gin `binding` e domínio.
- **Decisão:** validação de formato/range de entrada no **handler** (binding + trim de `place_id`); `ReviewCreatorUseCase` delega ao repositório após a borda aceitar.
- **Consequência:** novos callers (CLI, jobs) que chamem o use case direto precisam garantir entrada válida ou introduzir outra borda de validação.

---

## ADR-010 — `station` no write path da criação de review

- **Contexto:** tabela `station` exige `name` NOT NULL; ainda não há integração Places no backend; o posto precisa de linha agregada ao criar review.
- **Decisão:** no upsert após criar review, `name` no **INSERT** recebe o mesmo texto que `place_id` (placeholder); no **ON CONFLICT** atualizar só `total_score`, `review_count` e `updated_at` (não sobrescrever `name` já preenchido). Média baseada nas últimas reviews do posto (janela configurada no domínio).
- **Consequência:** nomes iguais entre postos distintos são permitidos (unicidade é `place_id`); nome amigável e endereço vêm depois com Places ou fluxo dedicado.
