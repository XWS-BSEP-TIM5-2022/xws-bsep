package api

import (
	"encoding/json"
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"net/http"
)

type UserHandler struct {
	userClientAddress string
}

func NewUserHandler(userClientAddress string) Handler {
	return &UserHandler{
		userClientAddress: userClientAddress,
	}
}

func (handler *UserHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/users/{userId}", handler.GetDetails)
	if err != nil {
		panic(err)
	}
}

func (handler *UserHandler) GetDetails(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["userId"]
	// if id == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// id := "1"
	name := "Pera"
	userDetails := &domain.User{Id: id, Name: name}
	fmt.Println(userDetails)
	// userClient := services.NewUserClient(handler.userClientAddress)
	// users, err := userClient.GetAll(context.TODO(), &user.GetAllRequest{})
	// fmt.Print(users)
	// err := handler.addOrderInfo(userDetails)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	response, err := json.Marshal(userDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
