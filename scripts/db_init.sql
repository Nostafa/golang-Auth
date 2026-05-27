-- Run against the postgres maintenance database (not your app DB):
--   psql "postgres://postgres:PASSWORD@HOST:PORT/postgres?sslmode=disable" -f scripts/db_init.sql
--
-- Or use: make db-init

CREATE DATABASE social;
