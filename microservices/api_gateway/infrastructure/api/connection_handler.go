package api

import (
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type ConnectionHandler struct {
	connectionClientAddress string
}

func NewConnectionHandler(connectionClientAddress string) Handler {
	return &ConnectionHandler{
		connectionClientAddress: connectionClientAddress,
	}
}

func (handler *ConnectionHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/connections/{userId}", handler.GetDetails) // TODO ???
	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionHandler) GetDetails(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["userId"]
	// if id == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// id := "1"
	userAid := "1"
	userBid := "2"
	connectionDetails := &domain.Connection{Id: id, UserAid: userAid, UserBid: userBid}
	fmt.Println(connectionDetails)
	// userClient := services.NewUserClient(handler.userClientAddress)
	// users, err := userClient.GetAll(context.TODO(), &user.GetAllRequest{})
	// fmt.Print(users)
	// err := handler.addOrderInfo(userDetails)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	response, err := json.Marshal(connectionDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
