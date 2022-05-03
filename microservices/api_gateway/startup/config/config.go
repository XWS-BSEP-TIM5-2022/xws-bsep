package config

import "os"

type Config struct {
	Port     string
	UserHost string
	UserPort string
	AuthHost string
	AuthPort string
}

func NewConfig() *Config {
	return &Config{
		Port:     os.Getenv("GATEWAY_PORT"),
		UserHost: os.Getenv("USER_SERVICE_HOST"),
		UserPort: os.Getenv("USER_SERVICE_PORT"),
		AuthHost: os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort: os.Getenv("AUTH_SERVICE_PORT"),
	}
}
