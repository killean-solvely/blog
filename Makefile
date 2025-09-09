# Load env from file
ifneq (,$(wildcard .env))
    include .env
    export
endif

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## Migration Commands

## migrate-create name=migration_name: create a new migration file
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then echo "Error: name parameter required. Usage: make migrate-create name=create_users_table"; exit 1; fi
	migrate create -ext sql -dir internal/infrastructure/persistence/sqlite/migrations -seq $(name)

## migrate-up: apply all pending migrations
.PHONY: migrate-up
migrate-up:
	migrate -path internal/infrastructure/persistence/sqlite/migrations -database "sqlite3://blog.db" up

## migrate-down: rollback the last migration
.PHONY: migrate-down
migrate-down:
	migrate -path internal/infrastructure/persistence/sqlite/migrations -database "sqlite3://blog.db" down 1

## migrate-down-all: rollback all migrations
.PHONY: migrate-down-all
migrate-down-all:
	migrate -path internal/infrastructure/persistence/sqlite/migrations -database "sqlite3://blog.db" down

## migrate-version: show current migration version
.PHONY: migrate-version
migrate-version:
	migrate -path internal/infrastructure/persistence/sqlite/migrations -database "sqlite3://blog.db" version

## migrate-force version=N: force set migration version (use with caution)
.PHONY: migrate-force
migrate-force:
	@if [ -z "$(version)" ]; then echo "Error: version parameter required. Usage: make migrate-force version=1"; exit 1; fi
	migrate -path internal/infrastructure/persistence/sqlite/migrations -database "sqlite3://blog.db" force $(version)
