DATABASE_URL=postgres://postgres:postgres@localhost:5433/post_confiavel?sslmode=disable

run:
	go run ./cmd/server

.PHONY: run

.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir extension/database/priv/migrations -seq $${name}

.PHONY: migrate-up
migrate-up:
	@migrate -path extension/database/priv/migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path extension/database/priv/migrations -database "$(DATABASE_URL)" down

.PHONY: migrate-drop
migrate-drop:
	@migrate -path extension/database/priv/migrations -database "$(DATABASE_URL)" drop
