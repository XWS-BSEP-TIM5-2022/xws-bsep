package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/startup/config"
	authGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	connectionGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	postGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	userGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
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

	connectionEndPoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndPoint, opts)
	if err != nil {
		panic(err)
	}

	postEmdpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEmdpoint, opts)
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

}

// ************** AUTHENTICATION - middleware ******************
type MuxWithMiddleware struct {
	mux *runtime.ServeMux
}

func (muxWithMiddleware *MuxWithMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := r.Method + " " + r.URL.Path
	if fullPath == "GET /users/getAllPublic" {
		// endpoints za neregistrovane korisnike
		fmt.Println("Ovaj zahtev nije potrebno validirati ni generisati token")
	} else if fullPath == "POST /user" || fullPath == "GET /login" {
		// sign in i login -> generisati token

	} else {
		// ostali endpoint-i -> potrebno validirati token
		if r.Header["Authorization"] == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		authorizationHeader := r.Header.Get("Authorization")
		fmt.Println("Auth header " + authorizationHeader)

		tokenString := strings.Split(authorizationHeader, " ")[1]
		fmt.Println("Token string " + tokenString)

		// authEmdpoint := fmt.Sprintf("auth_service:8000")
		// userEmdpoint := fmt.Sprintf("user_service:8000")

		// authHandler := api.NewAuthHandler(authEmdpoint, userEmdpoint)
		// authHandler.Init(muxWithMiddleware.mux)
	}

	muxWithMiddleware.mux.ServeHTTP(w, r)
}

func NewMuxWithMiddleware(handlerToWrap *runtime.ServeMux) *MuxWithMiddleware {
	return &MuxWithMiddleware{handlerToWrap}
}

// *************************************************

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}

func muxMiddleware(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(server.config.AuthHost + " -> " + server.config.AuthPort)

		fullPath := r.Method + " " + r.URL.Path
		if fullPath == "GET /users/getAllPublic" {
			// endpoints za neregistrovane korisnike
			fmt.Println("Ovaj zahtev nije potrebno validirati ni generisati token")
		} else if fullPath == "POST /user" || fullPath == "GET /login" {
			// sign in i login -> generisati token

		} else {
			// ostali endpoint-i -> potrebno validirati token
			if r.Header["Authorization"] == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			authorizationHeader := r.Header.Get("Authorization")
			fmt.Println("Auth header " + authorizationHeader)

			tokenString := strings.Split(authorizationHeader, " ")[1]
			fmt.Println("Token string " + tokenString)

			// authEmdpoint := fmt.Sprintf("auth_service:8000")
			// userEmdpoint := fmt.Sprintf("user_service:8000")

			// authHandler := api.NewAuthHandler(authEmdpoint, userEmdpoint)
			// authHandler.Init(muxWithMiddleware.mux)
		}

		server.mux.ServeHTTP(w, r)
	})
}
