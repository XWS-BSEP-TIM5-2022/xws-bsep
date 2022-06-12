package api

// type AuthHandler struct {
// 	authClientAddress string
// 	userClientAddress string
// }

// func NewAuthHandler(authClientAddress, userClientAddress string) Handler {
// 	return &AuthHandler{
// 		authClientAddress: authClientAddress,
// 		userClientAddress: userClientAddress,
// 	}
// }

// func (handler *AuthHandler) Init(mux *runtime.ServeMux) {
// 	err := mux.HandlePath("GET", "/auth/validate/{userId}", handler.TODO)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func (handler *AuthHandler) TODO(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
// 	id := pathParams["userId"]
// 	if id == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	userDetails := &domain.User{Id: id}
// 	fmt.Println(userDetails)
// 	userClient := services.NewUserClient(handler.userClientAddress)
// 	users, err := userClient.GetAll(context.TODO(), &user.GetAllRequest{})
// 	fmt.Print(users)
// 	// err := handler.addOrderInfo(userDetails)
// 	// if err != nil {
// 	// 	w.WriteHeader(http.StatusNotFound)
// 	// 	return
// 	// }
// 	response, err := json.Marshal(userDetails)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(response)
// }

// func (handler *AuthHandler) addOrderInfo(orderDetails *domain.User) error {
// 	return nil
// }
