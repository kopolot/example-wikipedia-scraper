#!/bin/sh
ctrl_c() {
  echo "Przerwanie debugowania..."
  kill -15 $TEST_PID 2>/dev/null
  exit 0
}
trap ctrl_c INT
go mod tidy
dlv debug cmd/notify/main.go \
    --listen=:2345 \
    --log \
    --headless=true \
    --accept-multiclient &
TEST_PID=$!
sleep infinity