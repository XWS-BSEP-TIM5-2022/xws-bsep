package main

import (
	"github.com/sanjadrinic/test_repo/microservices/api-gateway/startup"
	"github.com/sanjadrinic/test_repo/microservices/api-gateway/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
