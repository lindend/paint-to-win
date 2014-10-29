package main

type Config struct {
	Address string

	GameServerGamePort int
	GameServerApiPort  int

	RedisAddress string
}
