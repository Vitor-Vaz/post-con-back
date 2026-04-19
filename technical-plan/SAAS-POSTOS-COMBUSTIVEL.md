# PostoConfiável — plano de desenvolvimento (SaaS)

Documento vivo com visão de produto, decisões de escopo, roadmap técnico e integração com mapas. O código do backend e do front ainda não foi iniciado; este arquivo serve como referência para as próximas etapas.

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

Base para tabelas futuras, alinhado aos casos de uso **avaliar**, **ver avaliações**, **ver mapa**:

- **`reviews`**: identificador, `place_id` (Google), usuário (ou política para anônimo com limites), notas por dimensão (ex.: 1–5), **snapshot de contexto** na data da avaliação (combustível, perfil de veículo), timestamps.
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
| Ver avaliações | Listar por `place_id`, paginação, filtros por combustível/perfil; opcional score para o perfil atual. |
| Ver mapa | Combinar busca Places com a área solicitada; decisão de MVP: mapa tende a precisar Places para mostrar postos ainda sem review. |

---

## Roadmap sugerido (fases)

| Fase | Entrega |
|------|---------|
| 0 | Documentação (este arquivo), README, ADRs leves (`place_id` como chave; fórmula de score v1) |
| 1 | Esqueleto Go + Gin: `cmd/server`, config por env, healthcheck, logging, shutdown gracioso, migrações PostgreSQL |
| 2 | Domínio + persistência: `Review`, repositórios, índices em `place_id` e `created_at` |
| 3 | Integração Google Places: cliente em `internal/integration/googlemaps`, timeouts, retries, testes com mocks |
| 4 | API HTTP: handlers finos, use cases, validação, OpenAPI |
| 5 | Agregação: `place_scores`, política de recência |
| 6 | Qualidade: testes de integração (ex.: Postgres em container), CI |
| 7 | Front (futuro): stack a definir; consumo da mesma API e mapa |

---

## Riscos e mitigações

- **Fraude / reviews falsos**: rate limit, sinais de anomalia, moderação futura.
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

- Autenticação: JWT vs OAuth social; reviews anônimos e limites.
- Mapa: apenas bbox de reviews existentes vs sempre consultar Places na área.
- Moderação: manual, reportes, automação em fases.
- Internacionalização e marca: domínio `.com.br`, registro de marca (INPI), voz da marca **PostoConfiável**.

---

## Glossário

| Termo | Significado |
|-------|-------------|
| `place_id` | Identificador estável de um lugar no Google Places; chave externa das avaliações no PostoConfiável. |
| Snapshot de contexto | Combustível e perfil de veículo **no momento** da avaliação, para não misturar situações incompatíveis no agregado. |
| Score composto | Número derivado das dimensões e pesos (recência + similaridade de perfil na visualização). |

---

*Última atualização: documento de planejamento; implementação de código não iniciada.*
