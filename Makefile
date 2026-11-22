# Makefile

.PHONY: db-up db-down generate

DEV_DB_URL   ?= ./g7.db

# Generates all sqlc models AND templ HTML files
generate:
	sqlc generate

# Pushes new migrations to the database
db-up:
	goose -dir db/migrations sqlite "$(DEV_DB_URL)" up

# Rolls back the most recent migration
db-down:
	goose -dir db/migrations sqlite "$(DEV_DB_URL)" down
