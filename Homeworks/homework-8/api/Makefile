include .env

MIGRATION_FOLDER=internal/pkg/db/migrations

.PHONY: compose-up
compose-up:
	docker-compose build
	docker-compose up

.PHONY: compose-rm
compose-rm:
	docker-compose down

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down

.PHONY: unit-tests
unit-tests:
	go test ./internal/pkg/server -tags=unit

.PHONY: integration-tests
integration-tests:
	go test ./test/ -tags=integration

generate:
	protoc --proto_path=proto/ --go_out=internal/pb --go-grpc_out=internal/pb proto/*.proto


