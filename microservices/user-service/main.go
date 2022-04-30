package main

import (
	"fmt"
	startup "user-service/startup"
	cfg "user-service/startup/config"
)

func main() {
	fmt.Println("Hello world from user")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
