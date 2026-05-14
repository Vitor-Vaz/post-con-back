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
| Arquitetura | Clean Architecture + SOLID: `cmd/`, `internal/api`, `internal/domain`, `internal/usecase`, `internal/repository` (ou gateway), `internal/integration` (Google), `pkg/` se necessário |
| Front | Ainda não definido; a API deve expor contratos estáveis (JSON / OpenAPI) |
| Postos | Não armazenar o catálogo nacional como fonte da verdade; **vincular** tudo ao **`place_id`** do Google Places |
| Mapas | Integração com **Google Maps Platform** (Places para busca/detalhes; Maps JavaScript API no front quando existir) |

---

## Modelo de dados (conceitual)

Base para tabelas futuras, alinhado aos casos de uso **avaliar**, **ver avaliações**, **ver mapa**, **listar stations**:

- **`station`**: `place_id`, dados de listagem/detalhe, coordenadas quando houver; suporta mapa e **GET** por região (ver backlog: **range** em coluna vs só parâmetro de busca).
- **`reviews`**: identificador, `place_id` (Google), usuário (ou política para anônimo com limites), notas por dimensão (ex.: 1–5), **snapshot de contexto** na data da avaliação (combustível, perfil de veículo), timestamps. **No backend atual** a tabela `reviews` é **provisória** e será **revista** quando a **ponderação** e a política de **justiça com o posto** estiverem definidas (podem surgir colunas ou tabelas auxiliares).
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

| Caso de uso | Descrição |
|-------------|-----------|
| Avaliar | Criar review com `place_id` validado (ex.: Details uma vez), contexto veículo/combustível, dimensões; atualizar agregados. |
| Ver avaliações | `GET` por `place_id`: **preview** com `limit=10`; mesma rota com **`limit` maior + paginação** para lista completa; filtros por combustível/perfil e score por perfil ficam para depois do MVP simples. |
| Ver mapa | Combinar busca Places com a área solicitada; decisão de MVP: mapa tende a precisar Places para mostrar postos ainda sem review. |
| Listar stations | `GET` lista **genérica** no primeiro corte; evoluir para **região + range** (centro/viewport + distância) para alimentar mapa/lista. |
| Obter station | `GET` **detalhe** de uma `station` (id interno ou `place_id`). |

---

## Roadmap sugerido (fases)

| Fase | Entrega |
|------|---------|
| 0 | Documentação (este arquivo), README, ADRs leves (`place_id` como chave; fórmula de score v1) |
| 1 | Esqueleto Go + Gin: `cmd/server`, config por env, healthcheck, logging, shutdown gracioso, migrações PostgreSQL |
| 2 | Domínio + persistência: `Review`, repositórios, índices em `place_id` e `created_at` |
| 3 | Integração Google Places: cliente em `internal/integration/googlemaps`, timeouts, retries, testes com mocks |
| 4 | API HTTP: handlers finos, use cases, validação, OpenAPI — **GET** lista **`station`** (genérica), **GET** detalhe **`station`**, **GET** **`reviews`** por `place_id` (`limit=10` default / preview; mesma rota com paginação e `limit` maior) |
| 5 | Mapa na API: **GET `station` por região** (range + geo); fechar decisão de **range** em `station` vs só query; índices geográficos |
| 6 | Agregação: `place_scores`, política de recência |
| 7 | Modelagem e migrações: **revisar `reviews`** (e correlatos) para suportar **ponderação** e rastreabilidade; fórmula de equidade com o posto **a definir** com produto |
| 8 | **Anti-fraude e confiança**: regras contra reviews falsos; uso complementar de **Google Maps / Places** (validação contextual, sinais de lugar), dentro de termos, cotas e LGPD |
| 9 | Qualidade (integração/CI) e **front** futuro: stack a definir; consumo da mesma API e mapa |

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

- **`station` + range:** coluna(s) de **alcance** na tabela vs **range só na rota**; request (raio km, bbox); índices geo.
- **GET stations por região:** parâmetros obrigatórios, teto de linhas, ordenação; postos sem review via Places na mesma resposta ou fluxo separado.
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

*Última atualização: casos de uso e roadmap alinhados a **GET** `station` (lista, detalhe, por região/range), **GET** `reviews` (preview + paginação); MVP 1–5 com evolução para dimensões e score composto.*
