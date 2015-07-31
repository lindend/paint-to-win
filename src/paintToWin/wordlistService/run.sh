#!/bin/sh
if make; then
	../../../bin/wordlist -db "user=p2wuser password=devpassword host=localhost port=5432 dbname=paint2win sslmode=disable"\
		ApiPort=8007 \
		WordlistRoot=../../../../../paint-to-win-wordlists/ \
		Address=$(hostname -I)
fi
