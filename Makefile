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

test:
	@go test -v ./... -coverprofile=coverage.out

fill:
	@curl -s -X POST -H "Content-Type: application/json" -d '{"slug":"test"}' 'http://localhost:8080/api/v1/segment' | jq
	@curl -s -X POST -H "Content-Type: application/json" -d '{"slug":"test2"}' 'http://localhost:8080/api/v1/segment' | jq
	@curl -s -X POST -H "Content-Type: application/json" -d '{"slug":"test3"}' 'http://localhost:8080/api/v1/segment' | jq
	@curl -s -X POST -H "Content-Type: application/json" -d '{"slug":"test4"}' 'http://localhost:8080/api/v1/segment' | jq
	@curl -s -X PUT -H "Content-Type: application/json" -d '{"join": ["test", "test2"]}' 'http://localhost:8080/api/v1/user/322/segment'
	@curl -s -X PUT -H "Content-Type: application/json" -d '{"join": ["test3"], "leave": ["test2"]}' 'http://localhost:8080/api/v1/user/322/segment'
	@curl -s -X POST -H "Content-Type: application/json" -d '[{"slug":"test4", "ttl": 1}]' 'http://localhost:8080/api/v1/user/322/segment/ttl'
	@curl -s -H "Content-Type: application/json" 'http://localhost:8080/api/v1/user/322/segment' | jq
	@curl -s 'http://localhost:8080/api/v1/segment/report?time_range=2023-08&format=json' | jq


.PHONY: psql run