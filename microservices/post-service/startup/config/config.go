package config

import "os"

type Config struct {
	Port       string
	PostDBHost string
	PostDBPort string
}

func NewConfig() *Config {

	//os.Setenv("POST_SERVICE_PORT", "post_service")
	//os.Setenv("POST_DB_HOST", "post_db")
	//os.Setenv("POST_DB_PORT", "27017")

	return &Config{
		Port:       os.Getenv("POST_SERVICE_PORT"),
		PostDBHost: os.Getenv("POST_DB_HOST"),
		PostDBPort: os.Getenv("POST_DB_PORT"),
	}
}
