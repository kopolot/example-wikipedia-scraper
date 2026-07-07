#!/bin/sh
set -e

./wait-for-postgres.sh
echo "Starting API..."
exec ./api.bin
