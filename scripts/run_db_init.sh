#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

if [ -z "${DATABASE_URL:-}" ]; then
  echo "DATABASE_URL is not set. Add it to .env or export it." >&2
  exit 1
fi

# psql needs a maintenance DB; swap /social or /test for /postgres
ADMIN_URL="${DATABASE_URL//\/social/\/postgres}"
ADMIN_URL="${ADMIN_URL//\/test/\/postgres}"

echo "Creating database from scripts/db_init.sql ..."
psql "$ADMIN_URL" -v ON_ERROR_STOP=1 -f scripts/db_init.sql
