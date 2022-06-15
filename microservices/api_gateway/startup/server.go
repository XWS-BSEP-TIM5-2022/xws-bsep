package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/api"
	"github.com/gorilla/handlers"
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
	config       *cfg.Config
	mux          *runtime.ServeMux
	CustomLogger *api.CustomLogger
}

func NewServer(config *cfg.Config) *Server {
	CustomLogger := api.NewCustomLogger()
	server := &Server{
		config: config,
		mux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(customMatcher),
		),
		CustomLogger: CustomLogger,
	}
	server.initHandlers()
	server.initCustomHandlers()
	server.CustomLogger.SuccessLogger.Info("Starting api gateway successfully, PORT: ", config.Port) // TODO: ostaviti port ?
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

	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	err := userGw.RegisterUserServiceHandlerFromEndpoint(context.TODO(), server.mux, userEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("User service registration failed, PORT: ", server.config.UserPort, ", HOST: ", server.config.UserHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("User service registration successful") // TODO: dodati port i host ?

	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	err = authGw.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Auth service registration failed PORT: ", server.config.AuthPort, ", HOST: ", server.config.AuthHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Auth service registration successful") // TODO: dodati port i host ?

	connectionEndPoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndPoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Connection service registration failed PORT: ", server.config.ConnectionPort, ", HOST: ", server.config.ConnectionHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Connection service registration successful") // TODO: dodati port i host ?

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Post service registration failed PORT: ", server.config.PostPort, ", HOST: ", server.config.PostHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Post service registration successful") // TODO: dodati port i host ?
}

func (server *Server) initCustomHandlers() {
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	postsHandler := api.NewPostHandler(postEndpoint, connectionEndpoint, userEndpoint, authEndpoint)
	postsHandler.Init(server.mux)
}

func (server *Server) Start() {
	crtPath, _ := filepath.Abs("../server.crt")
	keyPath, _ := filepath.Abs("../server.key")
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), muxMiddleware(server)))

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://localhost:4200", "https://localhost:4200/**", "http://localhost:4200", "http://localhost:4200/**", "http://localhost:8080/**",
			"http://localhost:3000/**", "http://localhost:3000", "https://localhost:3000/**", "https://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Authorization", "Access-Control-Allow-Origin", "*"}),
		handlers.AllowCredentials(),
	)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", server.config.Port), crtPath, keyPath, cors(muxMiddleware(server))))
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), cors(muxMiddleware(server))))
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
