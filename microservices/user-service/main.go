package main

import (
	"fmt"

	startup "github.com/sanjadrinic/test_repo/microservices/user_service/startup"
	cfg "github.com/sanjadrinic/test_repo/microservices/user_service/startup/config"
)

func main() {
	fmt.Println("Hello world from user")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
