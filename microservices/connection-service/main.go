package main

import (
	"fmt"

	startup "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/startup"
	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/startup/config"
)

func main() {
	fmt.Println("Hello world from connection")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
