DATABASE_URL=postgres://postgres:postgres@localhost:5433/post_confiavel?sslmode=disable

run:
	go run ./cmd/server

.PHONY: run

.PHONY: fmt-check
fmt-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))" || (gofmt -l $$(find . -name '*.go' -not -path './vendor/*') >&2; exit 1)

.PHONY: lint
lint: fmt-check
	go vet ./...
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: sqlc-gen
sqlc-gen:
	@sqlc generate -f internal/gateway/postgres/sqlc.yaml

.PHONY: migrate-create
migrate-create:
	@read -p "Nome da migração: " name; \
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
