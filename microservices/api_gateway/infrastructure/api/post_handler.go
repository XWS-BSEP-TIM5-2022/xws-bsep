package api

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
//	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
//	"net/http"
//)
//
//type PostHandler struct {
//	postClientAddress string
//}
//
//func NewPostHandler(postClientAddress string) Handler {
//	return &PostHandler{
//		postClientAddress: postClientAddress,
//	}
//}
//
//func (handler *PostHandler) Init(mux *runtime.ServeMux) {
//	err := mux.HandlePath("GET", "/posts/{postId}", handler.GetDetails) // TODO ???
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (handler *PostHandler) GetDetails(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
//	id := pathParams["postId"]
//	// if id == "" {
//	// 	w.WriteHeader(http.StatusBadRequest)
//	// 	return
//	// }
//	// id := "1"
//	title := "Novi post"
//	dateCrated := "danas"
//	postDetails := &domain.Post{Id: id, Title: title, DateCreated: dateCrated}
//	fmt.Println(postDetails)
//	// userClient := services.NewUserClient(handler.userClientAddress)
//	// users, err := userClient.GetAll(context.TODO(), &user.GetAllRequest{})
//	// fmt.Print(users)
//	// err := handler.addOrderInfo(userDetails)
//	// if err != nil {
//	// 	w.WriteHeader(http.StatusNotFound)
//	// 	return
//	// }
//	response, err := json.Marshal(postDetails)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Write(response)
//}
