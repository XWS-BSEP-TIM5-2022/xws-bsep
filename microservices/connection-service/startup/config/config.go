package config

import "os"

type Config struct {
	Port                    string
	Host                    string
	ConnectionDBHost        string
	ConnectionDBPort        string
	ConnectionDBName        string // ?
	ConnectionDBUser        string
	ConnectionDBPass        string
	PublicKey               string
	LogsFolder              string
	InfoLogsFile            string
	DebugLogsFile           string
	ErrorLogsFile           string
	SuccessLogsFile         string
	WarningLogsFile         string
	NotificationServicePort string
	NotificationServiceHost string
	UserServicePort         string
	UserServiceHost         string
	// TODO SD: obrisati
	// Neo4jUri        string "bolt://localhost:7687",
	// Neo4jUsername   string  "neo4j",
	// Neo4jPassword   string  "password",
}

func NewConfig() *Config {
	return &Config{
		Port:                    os.Getenv("CONNECTION_SERVICE_PORT"),
		Host:                    os.Getenv("CONNECTION_SERVICE_HOST"),
		ConnectionDBHost:        os.Getenv("CONNECTION_DB_HOST"),
		ConnectionDBPort:        os.Getenv("CONNECTION_DB_PORT"),
		ConnectionDBUser:        os.Getenv("CONNECTION_DB_USER"),
		ConnectionDBPass:        os.Getenv("CONNECTION_DB_PASS"),
		PublicKey:               os.Getenv("PUBLIC_KEY"),
		LogsFolder:              os.Getenv("LOGS_FOLDER"),
		InfoLogsFile:            os.Getenv("INFO_LOGS_FILE"),
		DebugLogsFile:           os.Getenv("DEBUG_LOGS_FILE"),
		ErrorLogsFile:           os.Getenv("ERROR_LOGS_FILE"),
		SuccessLogsFile:         os.Getenv("SUCCESS_LOGS_FILE"),
		WarningLogsFile:         os.Getenv("WARNING_LOGS_FILE"),
		NotificationServiceHost: os.Getenv("NOTIFICATION_SERVICE_HOST"),
		NotificationServicePort: os.Getenv("NOTIFICATION_SERVICE_PORT"),
		UserServiceHost:         os.Getenv("USER_SERVICE_HOST"),
		UserServicePort:         os.Getenv("USER_SERVICE_PORT"),

		//Port:                    "8084",
		//Host:                    "localhost",
		//ConnectionDBHost:        "localhost",
		//ConnectionDBPort:        "7687",
		//ConnectionDBUser:        "neo4j",
		//ConnectionDBPass:        "password",
		//PublicKey:               "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		//LogsFolder:              "logs",
		//InfoLogsFile:            "/info.log",
		//DebugLogsFile:           "/debug.log",
		//ErrorLogsFile:           "/error.log",
		//SuccessLogsFile:         "/success.log",
		//WarningLogsFile:         "/warning.log",
		//NotificationServiceHost: "localhost",
		//NotificationServicePort: "8086",
		//UserServiceHost:         "localhost",
		//UserServicePort:         "8081",
	}

}
