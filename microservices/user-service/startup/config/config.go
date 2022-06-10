package config

type Config struct {
	Port            string
	UserDBHost      string
	UserDBPort      string
	UserDBName      string
	UserDBUser      string
	UserDBPass      string
	PublicKey       string
	LogsFolder      string
	InfoLogsFile    string
	DebugLogsFile   string
	ErrorLogsFile   string
	SuccessLogsFile string
	WarningLogsFile string
}

func NewConfig() *Config {
	return &Config{
		// Port:       os.Getenv("USER_SERVICE_PORT"),
		// UserDBHost: os.Getenv("USER_DB_HOST"),
		// UserDBPort: os.Getenv("USER_DB_PORT"),
		// UserDBName: os.Getenv("USER_DB_NAME"),
		// UserDBUser: os.Getenv("USER_DB_USER"),
		// UserDBPass: os.Getenv("USER_DB_PASS"),
		// PublicKey:  os.Getenv("PUBLIC_KEY"),
		// LogsFolder: os.Getenv("LOGS_FOLDER"),
		// LogsFolder:      os.Getenv("LOGS_FOLDER"),
		// InfoLogsFile:    os.Getenv("INFO_LOGS_FILE"),
		// DebugLogsFile:   os.Getenv("DEBUG_LOGS_FILE"),
		// ErrorLogsFile:   os.Getenv("ERROR_LOGS_FILE"),
		// SuccessLogsFile: os.Getenv("SUCCESS_LOGS_FILE"),
		// WarningLogsFile: os.Getenv("WARNING_LOGS_FILE"),

		Port:       "8081",
		UserDBHost: "localhost",
		UserDBPort: "27017",
		//UserDBName: "user_store",
		//UserDBUser: "postgres",
		//UserDBPass: "password",
		PublicKey:       "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		LogsFolder:      "logs",
		InfoLogsFile:    "/info.txt",
		DebugLogsFile:   "/debug.txt",
		ErrorLogsFile:   "/error.txt",
		SuccessLogsFile: "/success.txt",
		WarningLogsFile: "/warning.txt",
	}
}
