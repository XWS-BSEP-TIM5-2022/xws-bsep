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
}

func NewConfig() *Config {
	return &Config{
		Port:          "8001",
		Host:          "localhost",
		Neo4jUri:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "password",
		//Port:             os.Getenv("CONNECTION_SERVICE_PORT"),
		//ConnectionDBHost: os.Getenv("CONNECTION_DB_HOST"),
		//ConnectionDBPort: os.Getenv("CONNECTION_DB_PORT"),
		//ConnectionDBName: os.Getenv("CONNECTION_DB_NAME"),
		//ConnectionDBUser: os.Getenv("CONNECTION_DB_USER"),
		//ConnectionDBPass: os.Getenv("CONNECTION_DB_PASS"),
	}

}
