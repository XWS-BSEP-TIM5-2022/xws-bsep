package config

import "os"

type Config struct {
	Port                  string
	NotificationDBHost    string
	NotificationDBPort    string
	NotificationDBName    string
	NotificationDBMessage string
	NotificationDBPass    string
	PublicKey             string
	LogsFolder            string
	InfoLogsFile          string
	DebugLogsFile         string
	ErrorLogsFile         string
	SuccessLogsFile       string
	WarningLogsFile       string
	UserServicePort       string
	UserServiceHost       string
	MessageServicePort    string
	MessageServiceHost    string
	PostServicePort       string
	PostServiceHost       string
	ConnectionServicePort string
	ConnectionServiceHost string
}

func NewConfig() *Config {
	return &Config{
		Port:                  os.Getenv("NOTIFICATION_SERVICE_PORT"),
		NotificationDBHost:    os.Getenv("NOTIFICATION_DB_HOST"),
		NotificationDBPort:    os.Getenv("NOTIFICATION_DB_PORT"),
		PublicKey:             os.Getenv("PUBLIC_KEY"),
		LogsFolder:            os.Getenv("LOGS_FOLDER"),
		InfoLogsFile:          os.Getenv("INFO_LOGS_FILE"),
		DebugLogsFile:         os.Getenv("DEBUG_LOGS_FILE"),
		ErrorLogsFile:         os.Getenv("ERROR_LOGS_FILE"),
		SuccessLogsFile:       os.Getenv("SUCCESS_LOGS_FILE"),
		WarningLogsFile:       os.Getenv("WARNING_LOGS_FILE"),
		ConnectionServiceHost: os.Getenv("CONNECTION_SERVICE_HOST"),
		ConnectionServicePort: os.Getenv("CONNECTION_SERVICE_PORT"),
		UserServiceHost:       os.Getenv("USER_SERVICE_HOST"),
		UserServicePort:       os.Getenv("USER_SERVICE_PORT"),
		PostServiceHost:       os.Getenv("POST_SERVICE_HOST"),
		PostServicePort:       os.Getenv("POST_SERVICE_PORT"),
		MessageServiceHost:    os.Getenv("MESSAGE_SERVICE_HOST"),
		MessageServicePort:    os.Getenv("MESSAGE_SERVICE_PORT"),

		//Port:                  "8086",
		//NotificationDBHost:    "localhost",
		//NotificationDBPort:    "27017",
		//PublicKey:             "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		//LogsFolder:            "logs",
		//InfoLogsFile:          "/info.log",
		//DebugLogsFile:         "/debug.log",
		//ErrorLogsFile:         "/error.log",
		//SuccessLogsFile:       "/success.log",
		//WarningLogsFile:       "/warning.log",
		//ConnectionServiceHost: "localhost",
		//ConnectionServicePort: "8084",
		//UserServiceHost:       "localhost",
		//UserServicePort:       "8081",
		//MessageServiceHost:    "localhost",
		//MessageServicePort:    "8085",
		//PostServiceHost:       "localhost",
		//PostServicePort:       "8083",
	}
}
