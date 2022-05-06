package config

import "os"

type Config struct {
	Port            string
	AuthDBHost      string
	AuthDBPort      string
	AuthDBName      string
	AuthDBUser      string
	AuthDBPass      string
	UserServicePort string
	UserServiceHost string
	PrivateKey      string
	PublicKey       string
}

func NewConfig() *Config {
	return &Config{
		Port:            os.Getenv("AUTH_SERVICE_PORT"),
		AuthDBHost:      os.Getenv("AUTH_DB_HOST"),
		AuthDBPort:      os.Getenv("AUTH_DB_PORT"),
		AuthDBName:      os.Getenv("AUTH_DB_NAME"),
		AuthDBUser:      os.Getenv("AUTH_DB_USER"),
		AuthDBPass:      os.Getenv("AUTH_DB_PASS"),
		UserServicePort: os.Getenv("USER_SERVICE_PORT"),
		UserServiceHost: os.Getenv("USER_SERVICE_HOST"),
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		PublicKey:       os.Getenv("PUBLIC_KEY"),
	}
}
