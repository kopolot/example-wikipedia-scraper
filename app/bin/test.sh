#!/bin/sh
[ -z "$1" -o -z "$2" ] && { echo "Użycie: $0 <ścieżka_testu> <nazwa_testu>"; exit 1; }
if [ "$3" = "--debug" ] || [ "$3" = "-d" ]; then
  DEBUG=1
  echo "Tryb debugowania włączony"
else
  DEBUG=0
fi
ctrl_c() {
  echo "Przerwanie testu..."
  kill -15 $TEST_PID 2>/dev/null
  exit 0
}
trap ctrl_c INT
TEST_PID=
if [ $DEBUG -eq 1 ]; then
  export DELVE_MAX_VARIABLE_RECURSION=5
  export DELVE_MAX_STRING_LEN=512
  export DELVE_MAX_ARRAY_VALUES=500
  export DELVE_MAX_STRUCT_FIELDS=500
  dlv test \
    ./$1 \
    --headless \
    --listen=:2346 \
    --api-version=2 \
    --accept-multiclient \
    -- -test.run "^$2$" \
    -test.v -test.count=1 &
  TEST_PID=$!
else
  go test -count=1 -run "^$2$" ./$1 -v &
  TEST_PID=$!
fi
sleep infinity