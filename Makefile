.PHONY: generate

run:
	go run ./cmd/app

generate:
	sqlc generate

migrate:
	atlas schema apply --dev-url "docker://postgres" --url "postgresql://app:dev-pass@127.0.0.1:5432/app?sslmode=disable" --to "file://schema.sql"       