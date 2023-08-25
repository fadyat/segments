ifneq (,$(wildcard ./.env))
	include .env
	export
endif

MIGRATIONS_DIR = ./migrations/postgres
UP_STEP =
DOWN_STEP = -all
DATABASE_URL = postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

migrate-new:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) $(NAME)

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(DATABASE_URL) up $(UP_STEP)

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(DATABASE_URL) down $(DOWN_STEP)

psql:
	@docker-compose up -d psql

run:
	@go run cmd/segment/*.go

lint:
	@golangci-lint run --issues-exit-code 1 --print-issued-lines=true --config .golangci.yml ./...


.PHONY: psql run