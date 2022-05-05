package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/api"
	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/startup/config"
	authGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	userGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"

	// "golang.org/x/oauth2/jwt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	userEmdpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEmdpoint, opts)
	if err != nil {
		panic(err)
	}
	authEmdpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	err = authGw.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), server.mux, authEmdpoint, opts)
	if err != nil {
		panic(err)
	}
	// shippingEmdpoint := fmt.Sprintf("%s:%s", server.config.ShippingHost, server.config.ShippingPort)
	// err = shippingGw.RegisterShippingServiceHandlerFromEndpoint(context.TODO(), server.mux, shippingEmdpoint, opts)
	// if err != nil {
	// 	panic(err)
	// }
	// inventoryEmdpoint := fmt.Sprintf("%s:%s", server.config.InventoryHost, server.config.InventoryPort)
	// err = inventoryGw.RegisterInventoryServiceHandlerFromEndpoint(context.TODO(), server.mux, inventoryEmdpoint, opts)
	// if err != nil {
	// 	panic(err)
	// }
}

func (server *Server) initCustomHandlers() {
	userEmdpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)

	userHandler := api.NewUserHandler(userEmdpoint)
	userHandler.Init(server.mux)
}

// ************** AUTHENTICATION - middleware ******************
type MuxWithMiddleware struct {
	mux *runtime.ServeMux
}

func (muxWithMiddleware *MuxWithMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: odvojiti za koje metode i url-ove ne treba autentifikacija
	fullPath := r.Method + " " + r.URL.Path
	fmt.Println(fullPath)
	if fullPath == "POST /user" || fullPath == "GET /user" {
		//GenerateToken()
		// fmt.Println("body")
		// fmt.Println(r.Body)

		fmt.Printf("---------------------")
	} else {
		if r.Header["Authorization"] == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}

	muxWithMiddleware.mux.ServeHTTP(w, r)
}

func NewMuxWithMiddleware(handlerToWrap *runtime.ServeMux) *MuxWithMiddleware {
	return &MuxWithMiddleware{handlerToWrap}
}

// *************************************************

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), NewMuxWithMiddleware(server.mux)))
}
