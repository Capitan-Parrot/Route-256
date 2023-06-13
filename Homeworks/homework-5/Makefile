ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=123 dbname=ozon host=localhost port=5432 sslmode=disable
endif

MIGRATION_FOLDER=internal/pkg/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down