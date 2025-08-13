#!/usr/bin/env sh
set -e

if [ "$RUN_MIGRATIONS" = "true" ]; then
  make migrate
  exit 0
else
  echo "Skipping DB migrations."
fi

echo "Starting app…"
exec ./standup-logger
