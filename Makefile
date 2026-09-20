include .env
export

.PHONY: run db-up db-down db-status db-logs db-shell \
        migrate-up migrate-down migrate-version migration-new

run:
	go run ./cmd/api

db-up:
	docker compose up -d

db-down:
	docker compose down

db-status:
	docker compose ps

db-logs:
	docker compose logs -f postgres

db-shell:
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version

migration-new:
	migrate create -ext sql -dir migrations -seq $(name)