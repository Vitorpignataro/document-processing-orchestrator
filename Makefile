ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: \
	up \
	down \
	build \
	logs \
	ps \
	api \
	runner \
	test \
	tidy \
	migration \
	migrate-up \
	migrate-down \
	migrate-version \
	db-reset

MIGRATE_IMAGE=migrate/migrate:v4.18.1

DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp(mariadb:3306)/$(DB_NAME)

up:
	docker compose up -d --build

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

ps:
	docker compose ps

api:
	go run ./cmd/api

runner:
	go run ./cmd/runner

test:
	go test ./...

tidy:
	go mod tidy

migration:
	@test -n "$(name)" || (echo "Usage: make migration name=create_documents" && exit 1)
	docker run --rm \
		-v "$(PWD)/migrations:/migrations" \
		$(MIGRATE_IMAGE) \
		create \
		-ext sql \
		-dir /migrations \
		-seq \
		$(name)

migrate-up:
	docker compose run --rm migrate \
		-path /migrations \
		-database "$(DB_URL)" \
		up

migrate-down:
	docker compose run --rm migrate \
		-path /migrations \
		-database "$(DB_URL)" \
		down 1

migrate-version:
	docker compose exec mariadb \
		mariadb \
		-u$(DB_USER) \
		-p$(DB_PASSWORD) \
		$(DB_NAME) \
		-e "SELECT * FROM schema_migrations;"

db-reset:
	docker compose down -v
	docker compose up -d --build