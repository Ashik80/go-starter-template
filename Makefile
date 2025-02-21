include .env

dbDriver = postgres
dsn = user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} host=${DB_HOST} port=${DB_PORT} sslmode=${DB_SSL_MODE}
migrationsDir = ./pkg/migrations

db-status:
	@GOOSE_DRIVER="$(dbDriver)" GOOSE_DBSTRING="$(dsn)" goose status -dir $(migrationsDir)

db-migrate:
	@GOOSE_DRIVER="$(dbDriver)" GOOSE_DBSTRING="$(dsn)" goose create $(name) sql -dir $(migrationsDir)

db-upgrade:
	@GOOSE_DRIVER="$(dbDriver)" GOOSE_DBSTRING="$(dsn)" goose up -dir $(migrationsDir)

db-downgrade:
	@GOOSE_DRIVER="$(dbDriver)" GOOSE_DBSTRING="$(dsn)" goose down -dir $(migrationsDir)

tailwind-build-dev:
	@npx @tailwindcss/cli -i ./web/css/styles.css -o ./web/css/output.css --watch

tailwind-build:
	@npx @tailwindcss/cli -i ./web/css/styles.css -o ./web/css/output.css

build: tailwind-build
	@go build -o bin/main ./cmd/api/main.go

run: build
	@./bin/main
