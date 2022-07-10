package config

import "os"

type Config struct {
	Port             string
	Host             string
	ConnectionDBHost string
	ConnectionDBPort string
	ConnectionDBName string
	ConnectionDBUser string
	ConnectionDBPass string
	PublicKey        string
	LogsFolder       string
	InfoLogsFile     string
	DebugLogsFile    string
	ErrorLogsFile    string
	SuccessLogsFile  string
	WarningLogsFile  string
	Neo4jUri         string
	Neo4jUsername    string
	Neo4jPassword    string
}

func NewConfig() *Config {
	return &Config{
		Port:             os.Getenv("JOB_OFFER_SERVICE_PORT"),
		Host:             os.Getenv("JOB_OFFER_SERVICE_HOST"),
		ConnectionDBHost: os.Getenv("JOB_OFFER_DB_HOST"),
		ConnectionDBPort: os.Getenv("JOB_OFFER_DB_PORT"),
		ConnectionDBName: os.Getenv("JOB_OFFER_DB_NAME"),
		ConnectionDBUser: os.Getenv("JOB_OFFER_DB_USER"),
		ConnectionDBPass: os.Getenv("JOB_OFFER_DB_PASS"),
		LogsFolder:       os.Getenv("LOGS_FOLDER"),
		InfoLogsFile:     os.Getenv("INFO_LOGS_FILE"),
		DebugLogsFile:    os.Getenv("DEBUG_LOGS_FILE"),
		ErrorLogsFile:    os.Getenv("ERROR_LOGS_FILE"),
		SuccessLogsFile:  os.Getenv("SUCCESS_LOGS_FILE"),
		WarningLogsFile:  os.Getenv("WARNING_LOGS_FILE"),
		PublicKey:        os.Getenv("PUBLIC_KEY"),

		// Port:             "8089",
		// Host:             "localhost",
		// ConnectionDBHost: "localhost",
		// ConnectionDBPort: "7687",
		// Neo4jUri:         "bolt://localhost:7687",
		// ConnectionDBUser: "neo4j",
		// ConnectionDBPass: "password",
		// PublicKey:        "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		// LogsFolder:       "logs",
		// InfoLogsFile:     "/info.log",
		// DebugLogsFile:    "/debug.log",
		// ErrorLogsFile:    "/error.log",
		// SuccessLogsFile:  "/success.log",
		// WarningLogsFile:  "/warning.log",
	}

}
