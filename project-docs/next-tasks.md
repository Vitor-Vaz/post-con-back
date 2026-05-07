# Próximas tarefas

1. Criar a tabela **`station`** como overview de posto / station overview, ligada a `place_id`, com nome, endereço, coordenadas, `total_score`, `summary` e timestamps.
2. Revisar modelo **`reviews`** + domínio para **ponderação** (migrações quando a regra existir).
3. **Anti-fraude** + sinais **Maps/Places** (desenho antes de codar pesado).
4. **`users`** + auth + FK em `reviews.user_id`.
5. **GET** reviews por `place_id` (+ paginação, sqlc, testes).
6. **OpenAPI** dos contratos JSON.
7. Cliente **Places** em `internal/integration/…` (timeouts, mocks).
8. **`place_scores`** + recência (depois da ponderação).
9. **Health** mais rico só se precisar para deploy.

Atualizar este arquivo e `current-state.md` ao fechar itens.
