package main

import (
	"fmt"

	startup "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/startup"
	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/startup/config"
)

func main() {
	fmt.Println("Hello from message service")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
