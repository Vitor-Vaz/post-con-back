# Próximas tarefas

Prioridade sugerida; ajustar conforme produto.

1. **Revisão do modelo `reviews` e ponderação** — alinhar schema (migrações) e domínio com a **lógica de ponderação** a definir (equidade com o posto, recência, dimensões); pode exigir novas colunas/tabelas e migração de dados.
2. **Anti-fraude / reviews falsos** — desenho e implementação: rate limit, sinais de anomalia, regras de confiança; integração com **Google Maps / Places** onde agregar valor (validação de contexto, `place_id`, proximidade), respeitando LGPD e termos da plataforma.
3. **Usuários e autenticação** — tabela `users`, FK em `reviews.user_id`, política de identidade (JWT/OAuth a definir no plano).
4. **Leitura de reviews** — `GET /api/v1/reviews?place_id=...` com paginação; caso de uso + query sqlc + handler.
5. **OpenAPI** — documentar contratos JSON (fase 4 do plano técnico).
6. **Integração Google Places (núcleo)** — cliente em camada de integração (`internal/integration/...` conforme plano), validação de `place_id`, timeouts, retries, mocks de teste (complementa o item de anti-fraude quando for só “dados do lugar”).
7. **Agregados** — `place_scores`, recálculo e política de recência (depende da ponderação).
8. **CI** — lint, testes, migrações em pipeline (Postgres de serviço).
9. **Health/readiness** — enriquecer se necessário para orquestração (sem quebrar `/health` simples).

Registrar conclusões e mudanças de prioridade em `current-state.md` e decisões em `decisions.md`.
