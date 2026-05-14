#!/bin/sh
go mod tidy
air -c .air.notify.toml
