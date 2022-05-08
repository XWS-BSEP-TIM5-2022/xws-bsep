package config

type Config struct {
	//Port             string
	//ConnectionDBHost string
	//ConnectionDBPort string
	//ConnectionDBName string
	//ConnectionDBUser string
	//ConnectionDBPass string
	Port          string
	Host          string
	Neo4jUri      string
	Neo4jUsername string
	Neo4jPassword string
	PublicKey     string
}

func NewConfig() *Config {
	return &Config{
		Port:          "8084",
		Host:          "localhost",
		Neo4jUri:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "password",
		PublicKey:     "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",

		//Port:             os.Getenv("CONNECTION_SERVICE_PORT"),
		//ConnectionDBHost: os.Getenv("CONNECTION_DB_HOST"),
		//ConnectionDBPort: os.Getenv("CONNECTION_DB_PORT"),
		//ConnectionDBName: os.Getenv("CONNECTION_DB_NAME"),
		//ConnectionDBUser: os.Getenv("CONNECTION_DB_USER"),
		//ConnectionDBPass: os.Getenv("CONNECTION_DB_PASS"),
	}

}
