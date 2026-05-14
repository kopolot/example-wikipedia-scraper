#!/bin/sh
if [ -z "$1" ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi
migrate create -ext sql -dir ./migrations -seq $1