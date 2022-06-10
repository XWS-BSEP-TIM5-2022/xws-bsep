package config

type Config struct {
	Port            string
	PostDBHost      string
	PostDBPort      string
	PublicKey       string
	UserServicePort string
	UserServiceHost string
	AuthServicePort string
	AuthServiceHost string
	LogsFolder      string
	InfoLogsFile    string
	DebugLogsFile   string
	ErrorLogsFile   string
	SuccessLogsFile string
	WarningLogsFile string
}

func NewConfig() *Config {
	return &Config{
		//Port:       os.Getenv("POST_SERVICE_PORT"),
		//PostDBHost: os.Getenv("POST_DB_HOST"),
		//PostDBPort: os.Getenv("POST_DB_PORT"),
		//PublicKey:  os.Getenv("PUBLIC_KEY"),
		Port:            "8083",
		PostDBHost:      "localhost",
		PostDBPort:      "27017",
		PublicKey:       "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		UserServiceHost: "localhost",
		UserServicePort: "8081",
		AuthServiceHost: "localhost",
		AuthServicePort: "8082",
		LogsFolder:      "logs",
		InfoLogsFile:    "/info.txt",
		DebugLogsFile:   "/debug.txt",
		ErrorLogsFile:   "/error.txt",
		SuccessLogsFile: "/success.txt",
		WarningLogsFile: "/warning.txt",
	}
}
