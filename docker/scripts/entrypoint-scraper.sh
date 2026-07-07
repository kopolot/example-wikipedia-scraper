#!/bin/sh
set -e

./wait-for-postgres.sh
echo "Starting scraper..."
exec ./scraper.bin
