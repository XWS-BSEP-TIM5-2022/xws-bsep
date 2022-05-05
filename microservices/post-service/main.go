package main

import (
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/startup"
	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/startup/config"
)

func main() {
	fmt.Println("Hello world from posts")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
