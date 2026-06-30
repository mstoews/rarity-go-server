POSTGRES_DSN ?= postgres://rarity:rarity@localhost:5433/rarity?sslmode=disable

.PHONY: postgres migrate server

postgres:
	docker run --name rarity-pg -e POSTGRES_USER=rarity -e POSTGRES_PASSWORD=rarity \
	  -e POSTGRES_DB=rarity -p 5433:5432 -d postgres:16-alpine

migrate:
	psql "$(POSTGRES_DSN)" -f db/schema.sql

server:
	go run ./cmd/server
