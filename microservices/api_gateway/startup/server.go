package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxprom "gitlab.com/msvechla/mux-prometheus/pkg/middleware"

	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/startup/config"
	authGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	connectionGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	jobOfferGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	messageGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	notificationGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	postGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	userGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
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

	messageEndPoint := fmt.Sprintf("%s:%s", server.config.MessageHost, server.config.MessagePort)
	err = messageGw.RegisterMessageServiceHandlerFromEndpoint(context.TODO(), server.mux, messageEndPoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Message service registration failed PORT: ", server.config.MessagePort, ", HOST: ", server.config.MessageHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Message service registration successful") // TODO: dodati port i host ?

	notificationEndPoint := fmt.Sprintf("%s:%s", server.config.NotificationHost, server.config.NotificationPort)
	err = notificationGw.RegisterNotificationServiceHandlerFromEndpoint(context.TODO(), server.mux, notificationEndPoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Notification service registration failed PORT: ", server.config.NotificationPort, ", HOST: ", server.config.NotificationHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Notification service registration successful") // TODO: dodati port i host ?

	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err = postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Post service registration failed PORT: ", server.config.PostPort, ", HOST: ", server.config.PostHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Post service registration successful") // TODO: dodati port i host ?

	jobOfferEndpoint := fmt.Sprintf("%s:%s", server.config.JobOfferHost, server.config.JobOfferPort)
	err = jobOfferGw.RegisterJobOfferServiceHandlerFromEndpoint(context.TODO(), server.mux, jobOfferEndpoint, opts)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Job offer service registration failed PORT: ", server.config.JobOfferPort, ", HOST: ", server.config.JobOfferHost)
		panic(err)
	}
	server.CustomLogger.SuccessLogger.Info("Job offer service registration successful")
}

func (server *Server) initCustomHandlers() {
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	userEndpoint := fmt.Sprintf("%s:%s", server.config.UserHost, server.config.UserPort)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	jobOfferEndpoint := fmt.Sprintf("%s:%s", server.config.JobOfferHost, server.config.JobOfferPort)
	postsHandler := api.NewPostHandler(postEndpoint, connectionEndpoint, userEndpoint, authEndpoint)
	jobOfferHandler := api.NewJobOfferHandler(postEndpoint, connectionEndpoint, userEndpoint, authEndpoint, jobOfferEndpoint)
	postsHandler.Init(server.mux)
	jobOfferHandler.Init(server.mux)
}

func (server *Server) Start() {
	crtPath, _ := filepath.Abs("server.crt") // TODO
	keyPath, _ := filepath.Abs("server.key") // TODO
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), muxMiddleware(server)))

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://localhost:4200", "https://localhost:4200/**", "http://localhost:4200", "http://localhost:4200/**", "http://localhost:8080/**",
			"http://localhost:3000/**", "http://localhost:3000", "https://localhost:3000/**", "https://localhost:3000",
			"http://localhost:3001/**", "http://localhost:3001", "https://localhost:3001/**", "https://localhost:3001",
			"http://localhost:9090"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Authorization", "Access-Control-Allow-Origin", "*"}),
		handlers.AllowCredentials(),
	)
	r := mux.NewRouter()
	instrumentation := muxprom.NewDefaultInstrumentation()
	r.Use(instrumentation.Middleware)
	r.Path("/metrics").Handler(promhttp.Handler())
	r.PathPrefix("/").Handler(cors(muxMiddleware(server)))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", server.config.Port), crtPath, keyPath, r))
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
		server.mux.ServeHTTP(w, r)
	})
}

func AccessibleEndpoints() map[string]string {
	const authService = "/api/auth"
	const userService = "/api/user"
	const postService = "/api/post"
	const connectionService = "/api/connection"
	const jobOfferService = "/api/jobOffer"
	const messageService = "/api/message"
	const notificationService = "/api/notification"

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
