#!/bin/bash

if make build; then
	../../../bin/lobby -db "user=p2wuser password=devpassword host=localhost port=5432 dbname=paint2win sslmode=disable" \
		"LobbyApiPort=8083" \
		"RedisAddress=localhost:6379"
fi
