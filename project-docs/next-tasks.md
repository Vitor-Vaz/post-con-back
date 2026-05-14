# Próximas tarefas

1. Enriquecer **`station`** com nome, endereço e coordenadas reais (ex.: **Google Places** após `place_id`), substituindo o placeholder atual de `name`.
2. Revisar modelo **`reviews`** + domínio para **ponderação** (migrações quando a regra existir).
3. **Anti-fraude** + sinais **Maps/Places** (desenho antes de codar pesado).
4. **`users`** + auth + FK em `reviews.user_id`.
5. **GET** reviews por `place_id` (+ paginação, sqlc, testes).
6. **OpenAPI** dos contratos JSON.
7. Cliente **Places** em `internal/integration/…` (timeouts, mocks).
8. **`place_scores`** + recência (depois da ponderação, se separar de `station`).
9. **Health** mais rico só se precisar para deploy.

Atualizar este arquivo e `current-state.md` ao fechar itens.
