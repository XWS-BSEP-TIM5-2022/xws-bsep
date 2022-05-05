package config

type Config struct {
	//Port       string
	//UserDBHost string
	//UserDBPort string
	//UserDBName string
	//UserDBUser string
	//UserDBPass string
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
		Neo4jPassword: "connection",
		//Port:       os.Getenv("CONNECTION_SERVICE_PORT"),
		//UserDBHost: os.Getenv("CONNECTION_DB_HOST"),
		//UserDBPort: os.Getenv("CONNECTION_DB_PORT"),
		//UserDBName: os.Getenv("CONNECTION_DB_NAME"),
		//UserDBUser: os.Getenv("CONNECTION_DB_USER"),
		//UserDBPass: os.Getenv("CONNECTION_DB_PASS"),
	}

}
