package config

import "os"

type Config struct {
	Port           string
	UserHost       string
	UserPort       string
	AuthHost       string
	AuthPort       string
	PostHost       string
	PostPort       string
	ConnectionPort string
	ConnectionHost string
}

func NewConfig() *Config {
	return &Config{
		//Port:           "8080",
		//ConnectionHost: "localhost",
		//ConnectionPort: "8001",

		Port:     os.Getenv("GATEWAY_PORT"),
		UserHost: os.Getenv("USER_SERVICE_HOST"),
		UserPort: os.Getenv("USER_SERVICE_PORT"),
		AuthHost: os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort: os.Getenv("AUTH_SERVICE_PORT"),
		PostHost: os.Getenv("POST_SERVICE_HOST"),
		PostPort: os.Getenv("POST_SERVICE_PORT"),
	}
}
