package startup

import (
	"fmt"
	"log"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/tree/feat/user-servicesep/microservices/user_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/tree/feat/user-servicesep/microservices/user_service/startup/config"

	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	fmt.Println(postgresClient)

}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.UserDBHost, server.config.UserDBUser,
		server.config.UserDBPass, server.config.UserDBName,
		server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
