package config

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
	EmailPort       string
	EmailHost       string
	EmailFrom       string
	EmailPassword   string
	FrontendHost    string
	FrontendPort    string
}

func NewConfig() *Config {
	return &Config{
		//Port:            os.Getenv("AUTH_SERVICE_PORT"),
		//AuthDBHost:      os.Getenv("AUTH_DB_HOST"),
		//AuthDBPort:      os.Getenv("AUTH_DB_PORT"),
		//AuthDBName:      os.Getenv("AUTH_DB_NAME"),
		//AuthDBUser:      os.Getenv("AUTH_DB_USER"),
		//AuthDBPass:      os.Getenv("AUTH_DB_PASS"),
		//UserServicePort: os.Getenv("USER_SERVICE_PORT"),
		//UserServiceHost: os.Getenv("USER_SERVICE_HOST"),
		//PrivateKey:      os.Getenv("PRIVATE_KEY"),
		//PublicKey:       os.Getenv("PUBLIC_KEY"),

		Port:            "8082",
		AuthDBHost:      "localhost", // ?
		AuthDBPort:      "5432",
		AuthDBName:      "auth",
		AuthDBUser:      "postgres",
		AuthDBPass:      "admin", // admin //password
		UserServiceHost: "localhost",
		UserServicePort: "8081",
		PrivateKey:      "-----BEGIN RSA PRIVATE KEY-----\nMIIEpgIBAAKCAQEA0AzWYJTc9jiPn+RMNjMJhscn8hg/Mt0U22efM6IvM83CyQCi\nFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU14ncF8oN7utDQxdq6ivTsF1tMQtH\nWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EXPNtQtgAO9FF5gKzjkaBu6KzLU2RJ\nC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiVfDBWCkDylBTQTmI1rO9MdayVduuA\nzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksmYW4Bu0RxK5SP2byHj7pv93XsabA+\nQXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0gQIDAQABAoIBAQCq00+Wn8RKOyja\nnUZiXkedLZtD8dq6dcKxYATdDXb6byFXjziF2KyQW5NbKMyckLjRLV1Vd+8zIaz9\n44TQTgyQqXZY5iPIn81roie8PN4k+qE0WtFsAwT7zlNjqj06bFheBhu6ah0YHYrX\nsRMf+xPMeTttJOEihF9iYxK7uOI4EmXc9ISiKPT7lA9t1kLfiwGXJsVtjlYov7ny\ntIZhJPXoH5SAclsVm2UV26wRoJ1Y0qIKgn2q01QuJyDAS6rkHvYNnPBk3SQcPsUg\nXJ9tBdkAkMWNGvXt6Lwl2R2ehErKt5TUTFEzQO0a9PGg3upkJAh7B98O6BZnn4fq\n1Xe/ZQUBAoGBAPAE65o4tLffY/OzEguLGH494AaXko3wmEe1MHmNdnH2HIkf3KiA\n3XU9uxiHpUQoSRg23rrBgkgDZ9kh9ZQrs4R4IrNjIEnXYair1PBQvDyDk5YA49qG\n2fROVZhQNRYndmM969hlBT7zwZj2tf7hhg4hkkfWb6PaseGOsnlkjXeRAoGBAN3n\nArdAaYQN/uOIlBHWRepL15ptWMHYZpTceUck8JW+WaevImWe25A1a0fNw6LgnGhH\naDR/+CAJOYM+LeJQuCXVtTyG8h6T/xyO7+9StAkMojsPfUZ91r9OzlHp1gBdLUBz\nNoZ9RSJUPHoA99twCROPsWpuu/ki2W2ldtccRxXxAoGBAO3bwmxkc9uge1pABMsB\nvnUk9oUx4p/dZdvyWKatJUtMnfzaYX9vrYgJdAecLZC856sifVnQeT7KeTi6Kbf8\nEvxdXe4udwoWcwaHuw+owtKphjHqkeO3LfmpQ7QdEG7zDqTM8ZPSkP9Q63OeUr/T\nWVlZtbCRdrOIAC5Kjt40YumxAoGBALd46RLxbAzmsYgaBiuVWit11+d0X72vGmoc\nvR3o2g9F2sU9lhglt+7NbE1rQUWrp0bFO9CkulhqqCXuxGtqSEfoIjjQbuJ/haBs\nQtBDNl6BFqX0kaU2KNf25bpuuCWG5QJ0AHJEo2PV+Eb8A/No9+g3l/6jXkKI4PO6\nqr8DP3dRAoGBANJ2jlJ3IaC0Sp0y91OrwqUENF1NCn/288+h7Q7RizakdEh836oG\n1t6VlNOkz1+AJyd3ng+1Zb2r4TZhxu35ll3hMlaUPTI1EuKKJeGuJlVY8CbKJAN5\nk45pG2WEBlOn1XIcGuAWnyifLNL5Yyv7erWplyOcHwuUzS5u23GDVhvl\n-----END RSA PRIVATE KEY-----\n",
		PublicKey:       "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
		EmailPort:       "587",
		EmailHost:       "smtp.gmail.com",
		EmailFrom:       "dislinkt.e2@gmail.com",
		EmailPassword:   "Dislinkt123*",
		FrontendHost:    "localhost",
		FrontendPort:    "4200",
	}
}
