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