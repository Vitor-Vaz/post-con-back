# PostoConfiável — plano de desenvolvimento (SaaS)

Documento vivo com visão de produto, decisões de escopo, roadmap técnico e integração com mapas. Este arquivo serve como referência para as próximas etapas (backend **post-con-back** em evolução).

---

## Resumo executivo

**PostoConfiável** é um SaaS de avaliações de postos de combustível no Brasil, inspirado na ideia de “Glassdoor”, mas com foco em **confiança**, **contexto do veículo** e **informação que envelhece bem** (recência e dimensões explícitas), em contraste com avaliações genéricas e estáticas em plataformas como o Google Maps.

O problema central: há muitos postos com combustível adulterado ou bombas duvidosas, e hoje não existe um lugar que una **percepção estruturada**, **perfil de uso** (carro, combustível) e **atualização no tempo** para orientar quem abastece.

**Nome do produto (decidido):** **PostoConfiável** — alinhado à busca e à mensagem para o público e para comunicação institucional.

---

## Personas (resumo)

| Persona | Necessidade |
|--------|-------------|
| Motorista urbano (flex / gasolina comum) | Saber onde abastecer com menor risco de problema no carro e bom custo-benefício percebido. |
| Motorista com veículo mais exigente (ex.: turbo, injeção direta) | Avaliações que **não punam** o posto por “gasolina comum” quando o problema é mismatch de combustível; e vice-versa: problema em carro simples em posto “de nome” é sinal forte. |
| Viajante | Mapa + histórico recente por região, não só nota antiga. |

---

## Visão do produto

- **Problema**: Avaliações genéricas não refletem mudanças no posto nem o **encaixe** entre posto/combustível e **perfil do veículo**.
- **Diferencial PostoConfiável**: Avaliações **estruturadas** e uma **média ponderada** (ou score derivado) que considera:
  - tipo de combustível usado naquele abastecimento;
  - características do veículo (campos enumerados no MVP, evitando texto livre demais para o agregado);
  - dimensões explícitas (ex.: qualidade percebida do combustível, honestidade do litro/bomba, atendimento, pressão comercial);
  - **recência** (avaliações antigas pesam menos).
- **Mapa / estações**: tela com mapa e várias **`station`**; cliente envia referência geográfica + **range**; API devolve estações da região (detalhe de schema e contrato no backlog e no roadmap).
- **Pontuação**: **MVP** com nota única **1–5**; depois **detalhar avaliações** (dimensões + snapshot veículo/combustível) e evoluir para **score** composto / `place_scores` / recência (ver fases **6–7** do roadmap).

Sistema de **5 estrelas** como base, com regras de ponderação documentadas e versionáveis (ajustes sem reescrever todo o modelo de dados).

---

## Escopo técnico acordado

| Área | Decisão |
|------|---------|
| Backend | Go, framework **Gin**, **PostgreSQL** |
| Arquitetura | Clean Architecture + SOLID: `cmd/server`, `internal/app` (HTTP), `internal/domain` (casos de uso), `internal/gateway/postgres` (sqlc + repositórios), `extension/database` (migrações); `internal/integration` (Google) **a criar** |
| Front | Ainda não definido; a API deve expor contratos estáveis (JSON / OpenAPI) |
| Postos | Não armazenar o catálogo nacional como fonte da verdade; **vincular** tudo ao **`place_id`** do Google Places |
| Mapas | Integração com **Google Maps Platform** (Places para busca/detalhes; Maps JavaScript API no front quando existir) |

---

## Modelo de dados (conceitual)

Base para tabelas futuras, alinhado aos casos de uso **avaliar**, **ver avaliações**, **ver mapa**, **listar stations**:

- **`station`** *(implementada — migração `000003`)*: `id` (UUID), `place_id` (único), `name`, `address`, `latitude`, `longitude`, `total_score`, `review_count`, `summary`, `created_at`, `updated_at`. Hoje `total_score` / `review_count` são atualizados no **POST review** (média das últimas 100 avaliações do `place_id`; nome inicial do posto = `place_id` até integração Places). Coordenadas e endereço existem no schema, mas ainda não são preenchidos pela API.
- **`reviews`** *(implementada — migração `000002`)*: `id`, `place_id`, `user_id` (UUID, sem FK em `users`), `rating` (1–5), timestamps. **Provisória** para evolução de ponderação, dimensões e snapshot veículo/combustível.
- **`users`**: quando houver login; autenticação (JWT, OAuth, etc.) a decidir.
- **`vehicle_profile`**: no MVP pode ser colunas enumeradas ou JSON enxuto ligado à review (motor aspirado/turbo, faixa de exigência, flex predominante álcool vs gasolina, etc.).
- **`place_scores`**: agregados por `place_id` — médias por dimensão, score composto, contagem de reviews, `last_computed_at`; recálculo no write ou job assíncrono.

**Ponderação (regras de negócio):**

- Peso por **recência** (decaimento ou buckets, ex.: 30 / 90 / 180 dias).
- Na visualização “nota para o meu perfil”: filtrar ou ponderar reviews com **perfil similar** (MVP: mesmo tipo de combustível + mesma faixa de exigência do motor).

---

## Google Maps / Places — estratégia

1. **Chave canônica**: `place_id` em reviews e agregados.
2. **Mapa / listagem**: Nearby Search ou Text Search com viewport (centro + raio ou bounds); o backend não precisa ingerir todos os postos do país.
3. **Detalhes**: Place Details para validar posto, nome, endereço, coordenadas, status; opcional **cache** (TTL) por `place_id` para custo e latência — cache de **consultados**, não cópia nacional completa.
4. **Quotas**: limites por IP/usuário, cache de Details, evitar persistir mais do que o necessário vindos do Places.
5. **Compliance**: [Termos do Google Maps / Places](https://developers.google.com/maps/terms) — atribuição, políticas de cache e uso de dados.

---

## Casos de uso (backend)

| Caso de uso | Descrição | Status |
|-------------|-----------|--------|
| Avaliar | Criar review com `place_id`; atualizar agregados na `station`. | **Feito (MVP)** — `POST /api/v1/review`; sem validação Places nem dimensões extras. |
| Listar stations | Lista paginada de postos já persistidos. | **Feito** — `GET /api/v1/stations` (plural = coleção). |
| Obter station | Detalhe de um posto por `place_id`. | **Feito** — `GET /api/v1/station/:place_id` (singular = um recurso). |
| Ver avaliações | `GET` por `place_id` com preview e paginação. | **Pendente** |
| Ver mapa | Postos por região (geo + range), possivelmente Places. | **Pendente** |

### Convenção de rotas HTTP (decidida)

| Recurso | Padrão | Exemplo |
|---------|--------|---------|
| Coleção | substantivo **plural** | `GET /api/v1/stations` |
| Item único | substantivo **singular** + `place_id` no path | `GET /api/v1/station/:place_id` |

Chave canônica do posto na API: **`place_id`** (Google Places).

---

## API HTTP — entregue até agora

Rotas operacionais além de `/test` e `/health` (ping Postgres):

| Método | Caminho | Descrição |
|--------|---------|-----------|
| POST | `/api/v1/review` | Body JSON: `place_id`, `user_id` (UUID), `rating` (1–5). Cria review e faz upsert de `total_score` / `review_count` na `station`. |
| GET | `/api/v1/stations` | Lista postos; query `page` (opcional, default `1`); **10 itens por página**; resposta `{ data, pagination }`. |
| GET | `/api/v1/station/:place_id` | Um posto; resposta com todos os campos persistidos (`id`, `place_id`, `name`, `address`, coordenadas, scores, `summary`, timestamps). `404` se não existir. |

**Camadas:** handler (`internal/app/v1`) → caso de uso (`internal/domain`) → repositório + sqlc (`internal/gateway/postgres`). Validação de entrada na borda HTTP; erros de domínio mapeados para status HTTP.

**Persistência:** queries sqlc em `internal/gateway/postgres/queries/*.sql`; `make sqlc-gen` após alterar SQL.

---

## Estado do roadmap (progresso real)

| Fase | Entrega | Status |
|------|---------|--------|
| 0 | Documentação, README, plano técnico | **Em curso** (este arquivo + `AGENTS.md`) |
| 1 | Esqueleto Go + Gin, healthcheck, migrações Postgres, Docker local | **Feito** |
| 2 | Domínio + persistência: `reviews`, `station`, repositórios, sqlc | **Feito (MVP)** |
| 3 | Integração Google Places | **Pendente** |
| 4 | API HTTP: lista + detalhe `station`, GET `reviews`, OpenAPI | **Parcial** — lista e detalhe **feitos**; reviews e OpenAPI **pendentes** |
| 5 | Mapa: GET `station` por região (range + geo) | **Pendente** |
| 6 | Agregação: `place_scores`, recência | **Pendente** (hoje agregado simplificado na própria `station`) |
| 7 | Revisar `reviews` para ponderação | **Pendente** |
| 8 | Anti-fraude e confiança | **Pendente** |
| 9 | Front + qualidade ampliada | **Pendente** |

---

## Roadmap sugerido (fases) — visão de entregas

Referência de escopo futuro; progresso detalhado na tabela **Estado do roadmap** acima.

| Fase | Entrega |
|------|---------|
| 0 | Documentação (este arquivo), README, ADRs leves (`place_id` como chave; fórmula de score v1) |
| 1 | Esqueleto Go + Gin: `cmd/server`, config por env, healthcheck, migrações PostgreSQL |
| 2 | Domínio + persistência: `reviews`, `station`, repositórios sqlc |
| 3 | Integração Google Places: cliente em `internal/integration/googlemaps`, timeouts, retries, testes com mocks |
| 4 | API HTTP: **GET** `/api/v1/stations` (lista paginada), **GET** `/api/v1/station/:place_id` (detalhe), **GET** `reviews` por `place_id`, OpenAPI |
| 5 | Mapa na API: postos por **região** (range + geo); índices geográficos; decisão Places vs só base local |
| 6 | Agregação: `place_scores`, política de recência |
| 7 | Modelagem: **revisar `reviews`** para ponderação e contexto de veículo/combustível |
| 8 | Anti-fraude e confiança (incl. sinais Maps/Places dentro dos termos) |
| 9 | Front futuro e evolução de qualidade/deploy |

---

## Riscos e mitigações

- **Fraude / reviews falsos**: rate limit, sinais de anomalia, moderação futura; **camada anti-fraude** com regras explícitas no backend; onde fizer sentido, **Google Maps / Places** como sinal de confiança (ex.: consistência de `place_id` / contexto geográfico), sem depender só de um único sinal e respeitando termos e privacidade.
- **LGPD**: bases legais, privacidade, minimização de dados do veículo.
- **Expectativa legal/científica**: comunicar que a nota é **percepção e contexto**, não laudo ANP; opcionalmente educação sobre canais oficiais (ANP, Procon) sem prometer precisão química.

---

## Checklist Google Cloud / Maps Platform (quando for integrar)

- [ ] Projeto no Google Cloud com faturamento adequado ao uso esperado.
- [ ] APIs habilitadas: **Places API** (e novas APIs Places conforme documentação vigente), **Geocoding** apenas se necessário.
- [ ] Chave de API **restrita** (HTTP referrers para front; IPs para backend em produção).
- [ ] Cotas e alertas de billing configurados.
- [ ] Atribuição e uso conforme termos atuais do Google Maps Platform.

---

## Backlog de decisões pendentes

- **Preencher `station` com dados Places:** nome, endereço e coordenadas reais no create/upsert (hoje `name` = `place_id` no primeiro review).
- **`station` + range:** coluna(s) de **alcance** na tabela vs **range só na rota**; request (raio km, bbox); índices geo.
- **GET stations por região:** parâmetros obrigatórios, teto de linhas, ordenação; postos sem review via Places na mesma resposta ou fluxo separado.
- **Paginação da lista:** hoje offset por `page` fixo em 10 itens; evoluir cursor, `page_size` configurável ou teto máximo documentado.
- **GET reviews:** `limit` máximo; paginação (`cursor` vs offset); ordenação (ex.: `created_at` DESC).
- **Ponderação e justiça com o posto**: fórmula de pesos (recência, dimensões, similaridade de perfil), impacto em `reviews` e em `place_scores`.
- **Anti-fraude**: quais sinais (incluindo Google Maps / Places), limites de uso da API, retenção de dados e experiência do usuário legítimo.
- Autenticação: JWT vs OAuth social; reviews anônimos e limites.
- Mapa: apenas bbox de reviews existentes vs sempre consultar Places na área.
- Moderação: manual, reportes, automação em fases.
- Internacionalização e marca: domínio `.com.br`, registro de marca (INPI), voz da marca **PostoConfiável**.

---

## Glossário

| Termo | Significado |
|-------|-------------|
| `place_id` | Identificador estável de um lugar no Google Places; chave externa das avaliações no PostoConfiável. |
| `station` | Registro do posto na base da aplicação; suporta mapa e detalhe; amarrado a `place_id`. |
| Snapshot de contexto | Combustível e perfil de veículo **no momento** da avaliação, para não misturar situações incompatíveis no agregado. |
| Score composto | Número derivado das dimensões e pesos (recência + similaridade de perfil na visualização). |

---

*Última atualização: maio/2026 — backend com **POST review**, **GET stations** (lista paginada, 10/página), **GET station/:place_id** (detalhe); convenção plural/singular nas rotas; fases 1–2 concluídas no MVP; fase 4 parcial. Próximos passos naturais: **GET reviews**, Places, mapa por região.*
