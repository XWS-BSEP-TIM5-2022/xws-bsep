package startup

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"

	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/startup/config"
	authGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	connectionGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	postGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	userGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgo "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	config       *cfg.Config
	mux          *runtime.ServeMux
	CustomLogger *api.CustomLogger
	tracer       otgo.Tracer
	closer       io.Closer
}

const name = "api_gateway"

func NewServer(config *cfg.Config) *Server {
	CustomLogger := api.NewCustomLogger()
	tracer, closer := tracer.Init(name)
	otgo.SetGlobalTracer(tracer)
	server := &Server{
		config: config,
		mux: runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(customMatcher),
		),
		CustomLogger: CustomLogger,
		tracer:       tracer,
		closer:       closer,
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
	server.CustomLogger.SuccessLogger.Info("User service registration successful")
	log.Println("User service registration successful")

	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	err = authGw.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), server.mux, authEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Auth service registration failed PORT: ", server.config.AuthPort, ", HOST: ", server.config.AuthHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Auth service registration successful")

	connectionEndPoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionEndPoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Connection service registration failed PORT: ", server.config.ConnectionPort, ", HOST: ", server.config.ConnectionHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Connection service registration successful")

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Post service registration failed PORT: ", server.config.PostPort, ", HOST: ", server.config.PostHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Post service registration successful")
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
	crtPath, _ := filepath.Abs("server.crt")
	keyPath, _ := filepath.Abs("server.key")
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
		accessiblePermissions := AccessibleEndpoints()
		accessiblePermission, err := accessiblePermissions[r.URL.Path]
		if err {
			authorizationHeader := r.Header.Get("Authorization")
			tokenString := strings.Split(authorizationHeader, " ")

			Authorize(server, accessiblePermission, tokenString)
		}
		log.Println(server.config.AuthHost + ":" + server.config.AuthPort)

		endpointName := r.Method + " " + r.URL.Path
		span := tracer.StartSpanFromRequest(endpointName, server.tracer, r)
		defer span.Finish()

		server.mux.ServeHTTP(w, r)
	})
}

func AccessibleEndpoints() map[string]string {
	const authService = "/api/auth"
	const userService = "/api/user"
	const postService = "/api/post"
	const connectionService = "/api/connection"

	return map[string]string{
		authService + "/update":         "UpdateUsername",
		authService + "/changePassword": "UpdatePassword",
		authService + "/adminsEndpoint": "AdminsEndpoint",

		userService + "":                              "GetAllUsers",
		userService + "/updateBasicInfo":              "UpdateUserProfile",
		userService + "/updateExperienceAndEducation": "UpdateUserProfile",
		userService + "/updateSkillsAndInterests":     "UpdateUserProfile",
		userService + "/info":                         "GetLoggedInUserInfo",

		postService + "":         "CreatePost",
		postService + "/like":    "UpdatePostLikes",
		postService + "/dislike": "UpdatePostDislikes",
		postService + "/comment": "UpdatePostComments",
		postService + "/neutral": "NeutralPost",

		connectionService + "":          "CreateConnection",
		connectionService + "/register": "RegisterConnection",
		connectionService + "/reject":   "RejectConnection",
		connectionService + "/approve":  "ApproveConnection",
	}
}

func Authorize(server *Server, accessiblePermission string, values []string) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to parse public key")
		return
	}
	token, err := jwt.ParseWithClaims(
		values[1],
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				server.CustomLogger.ErrorLogger.Error("Unexpected token signing method")
				return nil, fmt.Errorf("Unexpected token signing method")
			}
			return publicKey, nil
		},
	)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Invalid token")
		return
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		server.CustomLogger.ErrorLogger.Error("Invalid token claims")
		return
	}

	foundPermission := false
	for _, jwtPermission := range claims.Permissions {
		if accessiblePermission == jwtPermission {
			foundPermission = true
		}
	}
	if foundPermission == false {
		server.CustomLogger.ErrorLogger.WithField("user", claims.Username).Error("Unauthorized")
	} else {
		server.CustomLogger.SuccessLogger.WithField("user", claims.Username).Info("Authorized")
	}
}

type UserClaims struct {
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}
