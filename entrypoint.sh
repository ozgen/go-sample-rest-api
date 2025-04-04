#!/bin/sh

echo "Running database migrations..."
/migrate up  # ✅ Use the prebuilt migrate binary

echo "Starting API..."
exec /api  # ✅ Ensure API runs as the main process
