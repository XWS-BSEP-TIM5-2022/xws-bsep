package job_offer_service

import (
	"fmt"

	startup "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/startup"
	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/startup/config"
)

func main() {
	fmt.Println("Hello world from job offer")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
