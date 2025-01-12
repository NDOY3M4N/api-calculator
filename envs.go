package main

import "os"

type Envs struct {
	DBString  string
	JWTSecret string
}

var envs = initEnv()

func initEnv() Envs {
	return Envs{
		DBString:  getEnv("DBSTRING", "./foo.db"),
		JWTSecret: getEnv("JWT_SECRET", "my-jwt-secret"),
	}
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}
