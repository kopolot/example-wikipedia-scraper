#!/bin/sh
set -e

host="${DB_HOST:-db}"
port="${DB_PORT:-5432}"
user="${DB_USER:-user}"
retries="${WAIT_RETRIES:-30}"
delay="${WAIT_DELAY:-2}"

echo "Waiting for PostgreSQL at ${host}:${port}..."

attempt=0
while [ "$attempt" -lt "$retries" ]; do
	if nc -z "$host" "$port" 2>/dev/null; then
		echo "PostgreSQL is available."
		exit 0
	fi

	attempt=$((attempt + 1))
	echo "Attempt ${attempt}/${retries}: PostgreSQL is not ready yet."
	sleep "$delay"
done

echo "PostgreSQL did not become ready in time." >&2
exit 1
