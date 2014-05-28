package main

type Config struct {
	WebsocketPort int
	ApiPort       int

	DbConnectionString string
	RedisAddress       string
}

func LoadConfig() Config {
	return Config{}
}
