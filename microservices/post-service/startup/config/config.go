package config

import "os"

type Config struct {
	Port                    string
	PostDBHost              string
	PostDBPort              string
	PublicKey               string
	UserServicePort         string
	UserServiceHost         string
	AuthServicePort         string
	AuthServiceHost         string
	LogsFolder              string
	InfoLogsFile            string
	DebugLogsFile           string
	ErrorLogsFile           string
	SuccessLogsFile         string
	WarningLogsFile         string
	NotificationServicePort string
	NotificationServiceHost string
}

func NewConfig() *Config {
	return &Config{
		Port:                    os.Getenv("POST_SERVICE_PORT"),
		PostDBHost:              os.Getenv("POST_DB_HOST"),
		PostDBPort:              os.Getenv("POST_DB_PORT"),
		PublicKey:               os.Getenv("PUBLIC_KEY"),
		UserServiceHost:         os.Getenv("USER_SERVICE_HOST"),
		UserServicePort:         os.Getenv("USER_SERVICE_PORT"),
		AuthServiceHost:         os.Getenv("AUTH_SERVICE_HOST"),
		AuthServicePort:         os.Getenv("AUTH_SERVICE_PORT"),
		NotificationServiceHost: os.Getenv("NOTIFICATION_SERVICE_HOST"),
		NotificationServicePort: os.Getenv("NOTIFICATION_SERVICE_PORT"),
		LogsFolder:              os.Getenv("LOGS_FOLDER"),
		InfoLogsFile:            os.Getenv("INFO_LOGS_FILE"),
		DebugLogsFile:           os.Getenv("DEBUG_LOGS_FILE"),
		ErrorLogsFile:           os.Getenv("ERROR_LOGS_FILE"),
		SuccessLogsFile:         os.Getenv("SUCCESS_LOGS_FILE"),
		WarningLogsFile:         os.Getenv("WARNING_LOGS_FILE"),

		// Port:            "8083",
		// PostDBHost:      "localhost",
		// PostDBPort:      "27017",
		// PublicKey:       "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		// UserServiceHost: "localhost",
		// UserServicePort: "8081",
		// AuthServiceHost: "localhost",
		// AuthServicePort: "8082",
		// LogsFolder:      "logs",
		//NotificationServiceHost: "localhost",
		//NotificationServicePort: "8086",
		// InfoLogsFile:    "/info.log",
		// DebugLogsFile:   "/debug.log",
		// ErrorLogsFile:   "/error.log",
		// SuccessLogsFile: "/success.log",
		// WarningLogsFile: "/warning.log",
	}
}
