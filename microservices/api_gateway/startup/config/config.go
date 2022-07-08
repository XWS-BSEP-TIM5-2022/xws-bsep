package config

type Config struct {
	Port             string
	UserHost         string
	UserPort         string
	AuthHost         string
	AuthPort         string
	PostHost         string
	PostPort         string
	ConnectionPort   string
	ConnectionHost   string
	MessageHost      string
	MessagePort      string
	NotificationHost string
	NotificationPort string
	LogsFolder       string
	InfoLogsFile     string
	DebugLogsFile    string
	ErrorLogsFile    string
	SuccessLogsFile  string
	WarningLogsFile  string
	PrivateKey       string
	PublicKey        string
	JobOfferHost     string
	JobOfferPort     string
	// ServerKey         string
	// ServerCertificate string
}

func NewConfig() *Config {
	return &Config{
		//Port:             os.Getenv("GATEWAY_PORT"),
		//UserHost:         os.Getenv("USER_SERVICE_HOST"),
		//UserPort:         os.Getenv("USER_SERVICE_PORT"),
		//AuthHost:         os.Getenv("AUTH_SERVICE_HOST"),
		//AuthPort:         os.Getenv("AUTH_SERVICE_PORT"),
		//PostHost:         os.Getenv("POST_SERVICE_HOST"),
		//PostPort:         os.Getenv("POST_SERVICE_PORT"),
		//MessageHost:      os.Getenv("MESSAGE_SERVICE_HOST"),
		//MessagePort:      os.Getenv("MESSAGE_SERVICE_PORT"),
		//NotificationHost: os.Getenv("NOTIFICATION_SERVICE_HOST"),
		//NotificationPort: os.Getenv("NOTIFICATION_SERVICE_PORT"),
		//ConnectionHost:   os.Getenv("CONNECTION_SERVICE_HOST"),
		//ConnectionPort:   os.Getenv("CONNECTION_SERVICE_PORT"),
		//PrivateKey:       os.Getenv("PRIVATE_KEY"),
		//PublicKey:        os.Getenv("PUBLIC_KEY"),
		//LogsFolder:       os.Getenv("LOGS_FOLDER"),
		//InfoLogsFile:     os.Getenv("INFO_LOGS_FILE"),
		//DebugLogsFile:    os.Getenv("DEBUG_LOGS_FILE"),
		//ErrorLogsFile:    os.Getenv("ERROR_LOGS_FILE"),
		//SuccessLogsFile:  os.Getenv("SUCCESS_LOGS_FILE"),
		//WarningLogsFile:  os.Getenv("WARNING_LOGS_FILE"),
		//JobOfferHost:     os.Getenv("JOB_OFFER_SERVICE_HOST"),
		//JobOfferPort:     os.Getenv("JOB_OFFER_SERVICE_PORT"),

		Port:             "8080",
		UserHost:         "localhost",
		UserPort:         "8081",
		AuthHost:         "localhost",
		AuthPort:         "8082",
		PostHost:         "localhost",
		PostPort:         "8083",
		ConnectionHost:   "localhost",
		ConnectionPort:   "8084",
		JobOfferHost:     "localhost",
		JobOfferPort:     "8089",
		MessageHost:      "localhost",
		MessagePort:      "8085",
		NotificationHost: "localhost",
		NotificationPort: "8086",
		LogsFolder:       "logs",
		InfoLogsFile:     "/info.log",
		DebugLogsFile:    "/debug.log",
		ErrorLogsFile:    "/error.log",
		SuccessLogsFile:  "/success.log",
		WarningLogsFile:  "/warning.log",
		PrivateKey:       "-----BEGIN RSA PRIVATE KEY-----\nMIIEpgIBAAKCAQEA0AzWYJTc9jiPn+RMNjMJhscn8hg/Mt0U22efM6IvM83CyQCi\nFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU14ncF8oN7utDQxdq6ivTsF1tMQtH\nWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EXPNtQtgAO9FF5gKzjkaBu6KzLU2RJ\nC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiVfDBWCkDylBTQTmI1rO9MdayVduuA\nzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksmYW4Bu0RxK5SP2byHj7pv93XsabA+\nQXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0gQIDAQABAoIBAQCq00+Wn8RKOyja\nnUZiXkedLZtD8dq6dcKxYATdDXb6byFXjziF2KyQW5NbKMyckLjRLV1Vd+8zIaz9\n44TQTgyQqXZY5iPIn81roie8PN4k+qE0WtFsAwT7zlNjqj06bFheBhu6ah0YHYrX\nsRMf+xPMeTttJOEihF9iYxK7uOI4EmXc9ISiKPT7lA9t1kLfiwGXJsVtjlYov7ny\ntIZhJPXoH5SAclsVm2UV26wRoJ1Y0qIKgn2q01QuJyDAS6rkHvYNnPBk3SQcPsUg\nXJ9tBdkAkMWNGvXt6Lwl2R2ehErKt5TUTFEzQO0a9PGg3upkJAh7B98O6BZnn4fq\n1Xe/ZQUBAoGBAPAE65o4tLffY/OzEguLGH494AaXko3wmEe1MHmNdnH2HIkf3KiA\n3XU9uxiHpUQoSRg23rrBgkgDZ9kh9ZQrs4R4IrNjIEnXYair1PBQvDyDk5YA49qG\n2fROVZhQNRYndmM969hlBT7zwZj2tf7hhg4hkkfWb6PaseGOsnlkjXeRAoGBAN3n\nArdAaYQN/uOIlBHWRepL15ptWMHYZpTceUck8JW+WaevImWe25A1a0fNw6LgnGhH\naDR/+CAJOYM+LeJQuCXVtTyG8h6T/xyO7+9StAkMojsPfUZ91r9OzlHp1gBdLUBz\nNoZ9RSJUPHoA99twCROPsWpuu/ki2W2ldtccRxXxAoGBAO3bwmxkc9uge1pABMsB\nvnUk9oUx4p/dZdvyWKatJUtMnfzaYX9vrYgJdAecLZC856sifVnQeT7KeTi6Kbf8\nEvxdXe4udwoWcwaHuw+owtKphjHqkeO3LfmpQ7QdEG7zDqTM8ZPSkP9Q63OeUr/T\nWVlZtbCRdrOIAC5Kjt40YumxAoGBALd46RLxbAzmsYgaBiuVWit11+d0X72vGmoc\nvR3o2g9F2sU9lhglt+7NbE1rQUWrp0bFO9CkulhqqCXuxGtqSEfoIjjQbuJ/haBs\nQtBDNl6BFqX0kaU2KNf25bpuuCWG5QJ0AHJEo2PV+Eb8A/No9+g3l/6jXkKI4PO6\nqr8DP3dRAoGBANJ2jlJ3IaC0Sp0y91OrwqUENF1NCn/288+h7Q7RizakdEh836oG\n1t6VlNOkz1+AJyd3ng+1Zb2r4TZhxu35ll3hMlaUPTI1EuKKJeGuJlVY8CbKJAN5\nk45pG2WEBlOn1XIcGuAWnyifLNL5Yyv7erWplyOcHwuUzS5u23GDVhvl\n-----END RSA PRIVATE KEY-----\n",
		PublicKey:        "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
	}
}
