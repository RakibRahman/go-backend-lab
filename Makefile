include .env
export

run:
	go run ./cmd/api/main.go

db-up:
	docker compose up -d

db-down:
	docker compose down

db-status:
	docker compose ps

db-logs:
	docker compose logs -f postgres

db-shell:
	docker compose exec postgres psql -U app -d backend_lab

.PHONY: migrate-up migrate-down migrate-version migration-new

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version

migration-new:
	migrate create -ext sql -dir migrations -seq $(name)