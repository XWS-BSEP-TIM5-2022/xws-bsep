package config

type Config struct {
	Port            string
	UserHost        string
	UserPort        string
	AuthHost        string
	AuthPort        string
	PostHost        string
	PostPort        string
	ConnectionPort  string
	ConnectionHost  string
	LogsFolder      string
	InfoLogsFile    string
	DebugLogsFile   string
	ErrorLogsFile   string
	SuccessLogsFile string
	WarningLogsFile string
}

func NewConfig() *Config {
	return &Config{

		Port:            "8080",
		UserHost:        "localhost",
		UserPort:        "8081",
		AuthHost:        "localhost",
		AuthPort:        "8082",
		PostHost:        "localhost",
		PostPort:        "8083",
		ConnectionHost:  "localhost",
		ConnectionPort:  "8084",
		LogsFolder:      "logs",
		InfoLogsFile:    "/info.log",
		DebugLogsFile:   "/debug.log",
		ErrorLogsFile:   "/error.log",
		SuccessLogsFile: "/success.log",
		WarningLogsFile: "/warning.log",

		//Port:     os.Getenv("GATEWAY_PORT"),
		//UserHost: os.Getenv("USER_SERVICE_HOST"),
		//UserPort: os.Getenv("USER_SERVICE_PORT"),
		//AuthHost: os.Getenv("AUTH_SERVICE_HOST"),
		//AuthPort: os.Getenv("AUTH_SERVICE_PORT"),
		//PostHost: os.Getenv("POST_SERVICE_HOST"),
		//PostPort: os.Getenv("POST_SERVICE_PORT"),
	}
}
