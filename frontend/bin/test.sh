#!/bin/sh
to nie działa jak sie nie myle :(
if [ -z "$1" ]; then
	echo "Użycie: $0 <nazwa_testu> [--debug]"
	exit 1
fi

if [ "$2" = "--debug" ] || [ "$2" = "-d" ]; then
	DEBUG=1
	echo "Tryb debugowania włączony"
else
	DEBUG=0
fi
if [ $DEBUG -eq 1 ]; then
	NODE_OPTIONS="--inspect-brk=0.0.0.0:9229" npx vitest run -t "$1"
else
	npx vitest run -t "$1"
fi