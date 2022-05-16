package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/api"
	"log"
	"net/http"
	"path/filepath"

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
		mux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(customMatcher),
		),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

// Metoda da bi se header sacuvao u servisima
func customMatcher(key string) (string, bool) {
	switch key {
	case "Authorization":
		return key, true
	default:
		return key, false
	}
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
}

func (server *Server) initCustomHandlers() {
	postEmdpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	connectionEmdpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	postsHandler := api.NewPostHandler(postEmdpoint, connectionEmdpoint)
	postsHandler.Init(server.mux)
}

func (server *Server) Start() {
	crtPath, _ := filepath.Abs("../localhost.crt")
	keyPath, _ := filepath.Abs("../localhost.key")
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", server.config.Port), crtPath, keyPath, muxMiddleware(server)))
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), muxMiddleware(server)))

}

func muxMiddleware(server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(server.config.AuthHost + ":" + server.config.AuthPort)
		server.mux.ServeHTTP(w, r)
	})
}

// func muxMiddleware(server *Server) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println(server.config.AuthHost + " -> " + server.config.AuthPort)

// 		fullPath := r.Method + " " + r.URL.Path
// 		if fullPath == "GET /users/getAllPublic" {
// 			// endpoints za neregistrovane korisnike
// 			fmt.Println("Ovaj zahtev nije potrebno validirati ni generisati token")
// 		} else if fullPath == "POST /user" || fullPath == "GET /login" {
// 			// sign in i login -> generisati token

// 		} else {
// 			// ostali endpoint-i -> potrebno validirati token
// 			if r.Header["Authorization"] == nil {
// 				http.Error(w, "unauthorized", http.StatusUnauthorized)
// 				return
// 			}
// 			authorizationHeader := r.Header.Get("Authorization")
// 			fmt.Println("Auth header " + authorizationHeader)

// 			tokenString := strings.Split(authorizationHeader, " ")[1]
// 			fmt.Println("Token string " + tokenString)

// 			// authEmdpoint := fmt.Sprintf("auth_service:8000")
// 			// userEmdpoint := fmt.Sprintf("user_service:8000")

// 			// authHandler := api.NewAuthHandler(authEmdpoint, userEmdpoint)
// 			// authHandler.Init(muxWithMiddleware.mux)
// 		}

// 		server.mux.ServeHTTP(w, r)
// 	})
// }
