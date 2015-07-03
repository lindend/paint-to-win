#!/bin/bash

if make build; then
	../../../bin/gameserver -db "user=p2wuser password=devpassword host=localhost port=5432 dbname=paint2win sslmode=disable" \
		"Address=$(hostname)" \
		"GameServerGamePort=8080"\
		"GameServerApiPort=8081"\
		"RedisAddress=localhost:6379"

fi