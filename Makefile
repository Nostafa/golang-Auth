# Load .env for DATABASE_URL (export so migrate sees it)
ifneq (,$(wildcard ./.env))
include .env
export
endif

MIGRATION_PATH := migrations

.PHONY: help db-init migration migrate-up migrate-down migrate-reset migrate-version

help:
	@echo "Targets:"
	@echo "  make db-init              Create 'social' DB (requires psql)"
	@echo "  make migration NAME=foo   Create migrations/<ts>_foo.{up,down}.sql"
	@echo "  make migrate-up           Apply pending migrations"
	@echo "  make migrate-down N=1     Roll back N migrations (default 1)"
	@echo "  make migrate-reset        Drop all tables + re-apply migrations"
	@echo "  make migrate-version      Show current migration version"

db-init:
	@bash scripts/run_db_init.sh

# Timestamped migrations (matches existing 20260525145907_* files; do not use -seq)
migration:
	@test -n "$(NAME)" || (echo "Usage: make migration NAME=create_posts_table" >&2 && exit 1)
	migrate create -ext sql -dir $(MIGRATION_PATH) $(NAME)

migrate-up:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is not set in .env" >&2 && exit 1)
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is not set in .env" >&2 && exit 1)
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down $(or $(N),1)

migrate-reset:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is not set in .env" >&2 && exit 1)
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" drop -f
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-version:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL is not set in .env" >&2 && exit 1)
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version
