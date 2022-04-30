package api

import (
	"api_gateway/domain"
	"api_gateway/infrastructure/services"
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices_demo/common/proto/user-service"
	//user "common/proto/user_service"
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
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userDetails := &domain.User{Id: id}

	err := handler.addOrderInfo(userDetails)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(userDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *UserHandler) addOrderInfo(userDetails *domain.User) error {
	orderingClient := services.NewUserClient(handler.userClientAddress)
	userInfo, err := orderingClient.Get(context.TODO(), &user.GetRequest{Id: userDetails.Id})
	if err != nil {
		return err
	}
	userDetails.Id = userInfo.User.Id
	userDetails.Name = userInfo.User.Name.String()
	return nil
}
