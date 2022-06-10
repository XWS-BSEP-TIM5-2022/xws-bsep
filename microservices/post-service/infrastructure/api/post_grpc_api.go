package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/startup/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// implementacije gRPC servera koji smo definisali u okviru common paketa

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service       *application.PostService
	WarningLogger *logrus.Entry
	InfoLogger    *logrus.Entry //*log.Logger //*zerolog.Logger
	ErrorLogger   *logrus.Entry
	SuccessLogger *logrus.Entry
	DebugLogger   *logrus.Entry
}

func NewPostHandler(service *application.PostService) *PostHandler {
	InfoLogger := setLogrusLogger(config.NewConfig().InfoLogsFile)
	ErrorLogger := setLogrusLogger(config.NewConfig().ErrorLogsFile)
	WarningLogger := setLogrusLogger(config.NewConfig().WarningLogsFile)
	SuccessLogger := setLogrusLogger(config.NewConfig().SuccessLogsFile)
	DebugLogger := setLogrusLogger(config.NewConfig().DebugLogsFile)

	return &PostHandler{
		InfoLogger:    InfoLogger,
		ErrorLogger:   ErrorLogger,
		WarningLogger: WarningLogger,
		SuccessLogger: SuccessLogger,
		DebugLogger:   DebugLogger,
		service:       service,
	}
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()

		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}

func setLogrusLogger(filename string) *logrus.Entry {
	mLog := logrus.New()
	mLog.SetReportCaller(true)

	logsFolderName := config.NewConfig().LogsFolder

	if _, err := os.Stat(logsFolderName); os.IsNotExist(err) {
		os.Mkdir(logsFolderName, 0777)
	}
	file, err := os.OpenFile(logsFolderName+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	mLog.SetOutput(mw)

	mLog.SetFormatter(&logrus.JSONFormatter{ //TextFormatter //JSONFormatter
		CallerPrettyfier: caller(),
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFile: "method",
		},
		// ForceColors: true,
	})
	contextLogger := mLog.WithFields(logrus.Fields{})
	return contextLogger
}

func (handler *PostHandler) logMessage(msg string, loger *log.Logger) {
	// TODO SD: escape karaktera
	loger.Println(strings.ReplaceAll(msg, " ", "_"))
}

func (handler *PostHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	/* TODO */
	str1 := "123#123$123%123^123&123*123(123)-+=|'.,!"

	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	str1 = re.ReplaceAllString(str1, "")
	fmt.Println(str1)
	/*  TODO  */

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("Post with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	postPb := mapPost(post) // prepakujemo iz domenskog modela u protobuf oblik
	response := &pb.GetResponse{
		Post: postPb,
	}
	handler.SuccessLogger.Info("Post by ID:" + objectId.Hex() + " received successfully")
	return response, nil
}

func (handler *PostHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	posts, err := handler.service.GetAll()
	if err != nil {
		handler.ErrorLogger.Error("Get all")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Posts: []*pb.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	handler.SuccessLogger.Info("Found " + strconv.Itoa(len(posts)) + " posts")
	return response, nil
}

func (handler *PostHandler) GetAllByUser(ctx context.Context, request *pb.GetRequest) (*pb.GetAllResponse, error) {
	id := request.Id
	posts, err := handler.service.GetAllByUser(id)
	if err != nil {
		handler.ErrorLogger.Error("Get all by userId: " + id)
		return nil, err
	}
	response := &pb.GetAllResponse{
		Posts: []*pb.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	handler.SuccessLogger.Info("Found " + strconv.Itoa(len(posts)) + " posts by user with ID:" + id)
	return response, nil
}

func (handler *PostHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	post, err := mapInsertPost(request.InsertPost)
	if err != nil {
		handler.ErrorLogger.Error("Post was not mapped")
		return nil, err
	}

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post.UserId = userId
	success, err := handler.service.Insert(post)
	if err != nil {
		handler.ErrorLogger.Error("Post was not inserted")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Post created by user with ID: " + post.UserId)
	return response, err
}

func (handler *PostHandler) InsertJobOffer(ctx context.Context, request *pb.InsertJobOfferRequest) (*pb.InsertResponse, error) {
	apiToken := request.InsertJobOfferPost.ApiToken

	username, err := handler.service.GetUsernameByApiToken(ctx, apiToken)
	if err != nil {
		handler.ErrorLogger.Error("Can not find username by api token")
		return nil, err
	}

	userId, err := handler.service.GetIdByUsername(ctx, username.Username)
	if err != nil {
		handler.ErrorLogger.Error("Can not find id by username: " + username.Username)
		return nil, err
	}

	post, err := mapInsertJobOfferPost(request.InsertJobOfferPost)
	if err != nil {
		handler.ErrorLogger.Error("Post was not mapped")
		return nil, err
	}

	post.UserId = userId.Id
	success, err := handler.service.Insert(post)
	if err != nil {
		handler.ErrorLogger.Error("Post was not inserted")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Post with job offer created by user with ID: " + post.UserId)
	return response, nil
}

func (handler *PostHandler) LikePost(ctx context.Context, request *pb.InsertLike) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("Post with ID: " + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)

	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
			handler.ErrorLogger.Error("User with ID: " + userId + " already liked selected post")
			return &pb.InsertResponse{
				Success: "error",
			}, err
		}
	}

	flag := false
	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
			handler.InfoLogger.Info("Deleting dislike by user with ID: " + userId)
			fmt.Println("user liked selected post, deleting dislike")
			flag = true
		}
	}

	postHelper.Dislikes = nil // prazan niz dislajkova
	if flag == true {
		for _, p := range post.Dislikes {
			if p.UserId != userId { // ubacujemo sve dislajkove osim onog koji je lajkovao
				postHelper.Dislikes = append(postHelper.Dislikes, p)
			}
		}
		post.Dislikes = postHelper.Dislikes
	}

	success, err := handler.service.LikePost(post, userId)
	if err != nil {
		handler.ErrorLogger.Error("Post was not liked by user with ID: " + userId)
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Post liked by user with ID: " + post.UserId)
	return response, err
}

func (handler *PostHandler) DislikePost(ctx context.Context, request *pb.InsertDislike) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("Post with ID: " + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)

	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
			handler.ErrorLogger.Error("User with ID: " + userId + " already disliked selected post") // TODO: specificirati post, a ne da ostane selected
			fmt.Println("user already dislikes selected post")
			return &pb.InsertResponse{
				Success: "error",
			}, err
		}
	}

	flag := false
	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
			handler.InfoLogger.Info("Deleting like by user with ID: " + userId)
			fmt.Println("user liked selected post, deleting like")
			flag = true
		}
	}

	postHelper.Likes = nil // prazan niz lajkova
	if flag == true {
		for _, p := range post.Likes {
			if p.UserId != userId { // ubacujemo sve lajkove osim onog koji je dislajkovao
				postHelper.Likes = append(postHelper.Likes, p)
			}
		}
		post.Likes = postHelper.Likes
	}

	success, err := handler.service.DislikePost(post, userId)
	if err != nil {
		handler.ErrorLogger.Error("Post was not disliked by user with ID: " + userId)
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Post disliked by user with ID: " + post.UserId)
	return response, err
}

func (handler *PostHandler) CommentPost(ctx context.Context, request *pb.InsertComment) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("Post with ID:" + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	success, err := handler.service.CommentPost(post, userId, strings.TrimSpace(request.Text)) // Trim - function to remove leading and trailing whitespace
	if err != nil {
		handler.ErrorLogger.Error("Post was not commented by user with ID: " + userId)
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Post commented by user with ID: " + userId)
	return response, err
}

func (handler *PostHandler) NeutralPost(ctx context.Context, request *pb.InsertNeutralReaction) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("Post with ID:" + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)

	flagDisliked := false
	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
			handler.InfoLogger.Info("Deleting like on post by user with ID: " + userId)
			fmt.Println("user already dislikes selected post - neutral")
			flagDisliked = true
		}
	}

	flagLiked := false
	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
			handler.InfoLogger.Info("Deleting dislike on post by user with ID: " + userId)
			fmt.Println("user true likes selected post - neutral")
			flagLiked = true
		}
	}

	postHelper.Likes = nil
	if flagLiked == true {
		for _, p := range post.Likes {
			if p.UserId != userId {
				postHelper.Likes = append(postHelper.Likes, p)
			}
		}
		post.Likes = postHelper.Likes
	}

	postHelper.Dislikes = nil
	if flagDisliked == true {
		for _, p := range post.Dislikes {
			if p.UserId != userId {
				postHelper.Dislikes = append(postHelper.Dislikes, p)
			}
		}
		post.Dislikes = postHelper.Dislikes
	}

	success, err := handler.service.Update(post)
	if err != nil {
		handler.ErrorLogger.Error("Neutral reaction on post by user with ID: " + userId + " was not successful")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Neutral reaction on post by user with ID: " + userId)
	return response, err
}

func (handler *PostHandler) UpdateCompanyInfo(ctx context.Context, request *pb.UpdateCompanyInfoRequest) (*pb.InsertResponse, error) {
	company, err := mapCompanyInfo(request.CompanyInfoDTO)
	if err != nil {
		handler.ErrorLogger.Error("Company with ID: " + request.CompanyInfoDTO.Id + " not found")
		return nil, err
	}

	success, err := handler.service.UpdateCompanyInfo(company, request.CompanyInfoDTO.OldName)
	if err != nil {
		handler.ErrorLogger.Error("Company with ID: " + request.CompanyInfoDTO.Id + " was not updated")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Company by name: " + company.Name + " updated")
	return response, err
}
