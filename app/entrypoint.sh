#!/usr/bin/env sh
set -e

echo "RUN_MIGRATIONS=${RUN_MIGRATIONS:-}"
echo "RUN_MIGRATIONS_EXIT=${RUN_MIGRATIONS_EXIT:-}"

# Ensure we are in /app so relative paths work
cd /app

if [ "${RUN_MIGRATIONS}" = "true" ]; then
  echo "Running DB migrations…"
  make migrate
  echo "Migrations complete."
  if [ "${RUN_MIGRATIONS_EXIT}" = "true" ]; then
    echo "Exiting after migrations as requested."
    exit 0
  fi
else
  echo "Skipping DB migrations."
fi

echo "Starting app…"
exec /app/standup-logger
