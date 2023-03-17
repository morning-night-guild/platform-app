#!/bin/bash -eu
#
# ex) useage
# ./entrypoint.sh export 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'
# 

go run main.go $1 $2
