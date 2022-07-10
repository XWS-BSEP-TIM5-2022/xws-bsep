package api

import (
	"context"
	"errors"
	"log"
	"net/mail"
	"regexp"
	"strconv"
	"unicode"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	service      *application.UserService
	CustomLogger *CustomLogger
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service *application.UserService) *UserHandler {
	CustomLogger := NewCustomLogger()
	return &UserHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	users, err := handler.service.GetAll(ctx)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(users)) + " public users")
	return response, nil
}

func (handler *UserHandler) GetAllPublic(ctx context.Context, request *pb.GetAllPublicRequest) (*pb.GetAllPublicResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllPublic")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	handler.CustomLogger.InfoLogger.Info("Getting all public accounts")
	users, err := handler.service.GetAllPublic(ctx)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Found " + strconv.Itoa(len(users)) + " public users")
		return nil, err
	}
	response := &pb.GetAllPublicResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(users)) + " public users")
	return response, nil
}

func (handler *UserHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Search")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestCriteria := re.ReplaceAllString(request.Criteria, " ")
	criteria := removeMalicious(requestCriteria)
	users, err := handler.service.Search(ctx, criteria)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Search error")
		return nil, err
	}

	response := &pb.SearchResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Number of users found after search: " + strconv.Itoa(len(users)))
	return response, nil
}
func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := checkEmailCriteria(request.User.Email)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User not inserted. Email invalid")
		return nil, err
	}

	err = checkUsernameCriteria(request.User.Username)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User not inserted. Username invalid")
		return nil, err
	}

	re, err := regexp.Compile(`[^\w\.\+\@]`)
	if err != nil {
		log.Fatal(err)
	}
	requestUserEmail := re.ReplaceAllString(request.User.Email, " ")
	handler.CustomLogger.InfoLogger.WithField("email", requestUserEmail).Info("User with email: " + requestUserEmail + " is updating profile")

	user := mapInsertUser(request.User)
	user, err = handler.service.Insert(user)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User is not inserted")
		return nil, err
	} else {
		handler.CustomLogger.SuccessLogger.Info("User inserted successfully")
		return &pb.InsertResponse{
			Id: user.Id.Hex(),
		}, nil
	}
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Update")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := checkEmailCriteria(request.User.Email)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User not updated. Email invalid")
		return nil, err
	}

	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestUserId := re.ReplaceAllString(request.User.Id, " ")

	objectId, err := primitive.ObjectIDFromHex(request.User.Id)
	handler.CustomLogger.InfoLogger.WithField("id", requestUserId).Info("User with ID: " + requestUserId + " is updating profile")
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created")
		return nil, err
	}
	oldUser, err := handler.service.Get(ctx, objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapUpdateUser(mapUser(oldUser), request.User)

	success, err := handler.service.Update(ctx, user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("User with ID: " + user.Id.Hex() + "updated successfully")
	return response, err
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	id := removeMalicious(request.Id)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")

	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Getting user by id: " + requestId)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	user, err := handler.service.Get(ctx, objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	userPb := mapUser(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	handler.CustomLogger.SuccessLogger.Info("User by ID:" + objectId.Hex() + " received successfully")
	return response, nil
}

func (handler *UserHandler) GetLoggedInUserInfo(ctx context.Context, request *pb.GetAllRequest) (*pb.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetLoggedInUserInfo")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	handler.CustomLogger.InfoLogger.WithField("id", userId).Info("Getting logged in user infos by id: " + userId)
	user, err := handler.service.GetById(ctx, userId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + userId + " not found")
		return nil, err
	}
	pbUser := mapUser(user)
	handler.CustomLogger.SuccessLogger.Info("User received successfully" + pbUser.Email)
	return pbUser, nil
}

func (handler *UserHandler) UpdateBasicInfo(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateBasicInfo")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	handler.CustomLogger.InfoLogger.WithField("id", id).Info("Updating basic info by user with ID: " + id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(ctx, objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapBasicInfo(mapUser(oldUser), request.User)

	success, err := handler.service.Update(ctx, user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Basic info updated successfully")
	return response, err
}

func (handler *UserHandler) UpdateExperienceAndEducation(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateExperienceAndEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	handler.CustomLogger.InfoLogger.WithField("id", id).Info("Updating experience and education by user with ID: " + id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(ctx, objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapExperienceAndEducation(mapUser(oldUser), request.User)
	success, err := handler.service.Update(ctx, user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Experience and education updated successfully")
	return response, err
}

func (handler *UserHandler) UpdateSkillsAndInterests(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateSkillsAndInterests")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	handler.CustomLogger.InfoLogger.WithField("id", id).Info("Updating skilss and interests by user with ID: " + id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(ctx, objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapSkillsAndInterests(mapUser(oldUser), request.User)
	success, err := handler.service.Update(ctx, user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Skills and interests updated successfully")
	return response, err
}

func (handler *UserHandler) GetEmail(ctx context.Context, request *pb.GetRequest) (*pb.GetEmailResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetEmail")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	id := removeMalicious(request.Id)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")

	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Get informations about email by user with ID: " + requestId)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	user, err := handler.service.Get(ctx, objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	if !user.IsActive {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + request.Id + " is not activated")
		return nil, errors.New("Account is not activated")
	}
	handler.CustomLogger.SuccessLogger.Info("User email received successfully")
	response := &pb.GetEmailResponse{
		Email: user.Email,
	}
	return response, nil
}
func (handler *UserHandler) UpdateIsActiveById(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateIsActiveById")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")
	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Checking active status by user with id: " + requestId)

	err = handler.service.UpdateIsActiveById(ctx, removeMalicious(request.Id))
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + request.Id + " not activated")

		return &pb.ActivateAccountResponse{
			Success: err.Error(),
		}, err
	}
	handler.CustomLogger.SuccessLogger.Info("User by ID:" + request.Id + " received successfully")
	return &pb.ActivateAccountResponse{
		Success: "Success",
	}, nil
}

func (handler *UserHandler) GetIsActive(ctx context.Context, request *pb.GetRequest) (*pb.IsActiveResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIsActive")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")
	user, err := handler.service.GetById(ctx, removeMalicious(request.Id))
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + requestId + " not found")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("User by id received successfully")
	return &pb.IsActiveResponse{
		IsActive: user.IsActive,
	}, nil
}

func (handler *UserHandler) GetIdByEmail(ctx context.Context, request *pb.GetIdByEmailRequest) (*pb.InsertResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIdByEmail")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := checkEmailCriteria(request.Email)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User not found. Email invalid")
		return nil, err
	}

	re, err := regexp.Compile(`[^\w\.\+\@]`)
	if err != nil {
		log.Fatal(err)
	}
	requestEmail := re.ReplaceAllString(request.Email, " ")
	userId, err := handler.service.GetIdByEmail(ctx, request.Email)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with email: " + requestEmail + " not found")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("User by email received successfully")
	return &pb.InsertResponse{
		Id: userId,
	}, nil
}

func (handler *UserHandler) GetIdByUsername(ctx context.Context, request *pb.GetIdByUsernameRequest) (*pb.InsertResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIdByUsername")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	err := checkUsernameCriteria(request.Username)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User not found. Username invalid")
		return nil, err
	}
	re, err := regexp.Compile(`[^\w\.]`)
	if err != nil {
		log.Fatal(err)
	}
	requestUsername := re.ReplaceAllString(request.Username, " ")
	user, err := handler.service.GetByUsername(ctx, request.Username)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with username: " + requestUsername + " not found")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("User received successfully")
	return &pb.InsertResponse{
		Id: user.Id.Hex(),
	}, nil
}

func (handler *UserHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user := mapInsertUserSaga(request)
	err := handler.service.Create(user, request.Username, request.Password)
	if err != nil {
		return &pb.RegisterResponse{
			StatusCode: "500",
			Message:    "Something wrong, please try again",
		}, err
	}
	return &pb.RegisterResponse{
		StatusCode: "200",
		Message:    "OK",
	}, nil
}

func checkEmailCriteria(email string) error {
	if len(email) == 0 {
		return errors.New("Email should not be empty")
	}
	_, err := mail.ParseAddress(email)

	if err != nil {
		return errors.New("Email is invalid.")
	}
	return nil
}

func checkUsernameCriteria(username string) error {
	if len(username) == 0 {
		return errors.New("Username should not be empty")
	}

	for _, char := range username {

		if unicode.IsSpace(int32(char)) {
			return errors.New("Username should not contain any spaces")
		}
		if char == '$' {
			return errors.New("Username should not contain '$'")
		}
	}
	return nil
}

func (handler *UserHandler) UpdatePostNotification(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdatePostNotification")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	objectId, err := primitive.ObjectIDFromHex(request.User.Id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created")
		return nil, err
	}
	oldUser, err := handler.service.Get(ctx, objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapUpdateNotificationUser(mapUser(oldUser), request.User)

	success, err := handler.service.Update(ctx, user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("User with ID: " + user.Id.Hex() + "updated successfully")
	return response, err
}

func (handler *UserHandler) UpdatePrivacy(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdatePrivacy")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)
	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	handler.CustomLogger.InfoLogger.WithField("id", id).Info("Updating privacy info by user with ID: " + id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(ctx, objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapPrivacyInfo(mapUser(oldUser), request.User)

	success, err := handler.service.UpdatePrivacy(ctx, user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Basic info updated successfully")
	return response, err
}
