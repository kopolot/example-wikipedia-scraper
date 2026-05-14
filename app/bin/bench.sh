#!/bin/sh
go test -bench=^$2$ -run=^$ -benchmem  -cpuprofile=cpu.prof -memprofile=mem.prof ./$1 -v