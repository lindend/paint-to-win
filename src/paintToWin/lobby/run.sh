#!/bin/bash

make build
../../../bin/lobby -db "user=p2wuser password=devpassword host=localhost port=5432 dbname=paint2win sslmode=disable"
