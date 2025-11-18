.PHONY: run migrate seed

run:
	go run cmd/main.go

migrate:
	DB_HOST=localhost DB_PORT=5433 go run migration/main.go

seed:
	DB_HOST=localhost DB_PORT=5433 go run data/main.go