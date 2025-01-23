ent-install:
	@go get entgo.io/ent/cmd/ent

ent-new:
	@go run -mod=mod entgo.io/ent/cmd/ent new $(name)

ent-gen:
	@go generate ./ent

db-status:
	@atlas migrate status \
		--dir "file://ent/migrate/migrations" \
		--url "postgresql://postgres:postgres@localhost:5432/test_temp?search_path=public"

db-migrate:
	@atlas migrate diff $(name) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "postgresql://postgres:postgres@localhost:5432/test_temp?search_path=public"

db-upgrade:
	@atlas migrate apply \
		--dir "file://ent/migrate/migrations" \
		--url "postgresql://postgres:postgres@localhost:5432/test_temp?search_path=public"

db-downgrade:
	@atlas migrate down \
		--dir "file://ent/migrate/migrations" \
		--url "postgresql://postgres:postgres@localhost:5432/test_temp?search_path=public" \
		--dev-url "postgresql://postgres:postgres@localhost:5432/test_temp?search_path=public"

db-clean:
	@atlas schema clean -u "postgresql://postgres:postgres@localhost:5432/test_temp?search_path=public"

build:
	@go build -o bin/main ./cmd/api/main.go

run: build
	@./bin/main
