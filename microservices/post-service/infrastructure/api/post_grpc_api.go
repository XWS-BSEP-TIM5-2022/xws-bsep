package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// implementacije gRPC servera koji smo definisali u okviru common paketa

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service      *application.PostService
	CustomLogger *CustomLogger
}

func NewPostHandler(service *application.PostService) *PostHandler {
	CustomLogger := NewCustomLogger()
	return &PostHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *PostHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	/* sanitizacija */
	id := request.Id
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post with ID: " + id + " not found")
		return nil, err
	}
	postPb := mapPost(post) // prepakujemo iz domenskog modela u protobuf oblik
	response := &pb.GetResponse{
		Post: postPb,
	}
	handler.CustomLogger.SuccessLogger.Info("Post with ID: " + id + " received successfully")
	return response, nil
}

func (handler *PostHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	posts, err := handler.service.GetAll()
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all posts unsuccessful")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Posts: []*pb.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(posts)) + " posts")
	return response, nil
}

func (handler *PostHandler) GetAllByUser(ctx context.Context, request *pb.GetRequest) (*pb.GetAllResponse, error) {
	/* sanitizacija unosa */
	id := request.Id
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	posts, err := handler.service.GetAllByUser(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all by userId: " + id)
		return nil, err
	}
	response := &pb.GetAllResponse{
		Posts: []*pb.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(posts)) + " posts created by user with ID: " + id)
	return response, nil
}

func (handler *PostHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	post, err := mapInsertPost(request.InsertPost)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post was not mapped")
		return nil, err
	}

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post.UserId = userId
	success, err := handler.service.Insert(post)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post was not inserted")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Post with ID: " + post.Id.Hex() + " created by user with ID: " + post.UserId)
	return response, err
}

// TODO:  prebaciti u servis ?
func (handler *PostHandler) InsertJobOffer(ctx context.Context, request *pb.InsertJobOfferRequest) (*pb.InsertResponse, error) {
	post, err := mapInsertJobOfferPost(request.InsertJobOfferPost)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post was not mapped")
		return nil, err
	}

	/* sanitizacija unosa */
	id := request.InsertJobOfferPost.UserId
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, "")
	post.UserId = id

	success, err := handler.service.Insert(post)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post was not inserted")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Job offer post with ID: " + post.Id.Hex() + " created by user with ID: " + id)
	return response, nil
}

func (handler *PostHandler) LikePost(ctx context.Context, request *pb.InsertLike) (*pb.InsertResponse, error) {
	/* sanitizacija */
	id := request.PostId
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	objectId, err := primitive.ObjectIDFromHex(id)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post with ID: " + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)

	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
			handler.CustomLogger.ErrorLogger.Error("User with ID: " + userId + " already liked post with ID: " + id)
			return &pb.InsertResponse{
				Success: "error",
			}, err
		}
	}

	flag := false
	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
			handler.CustomLogger.InfoLogger.Info("Deleting dislike on post with ID: " + id + " by user with ID: " + userId)
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
		handler.CustomLogger.ErrorLogger.Error("Post with ID: " + post.Id.Hex() + " was not liked by user with ID: " + userId)
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Post with ID: " + post.Id.Hex() + " liked by user with ID: " + post.UserId)
	return response, err
}

func (handler *PostHandler) DislikePost(ctx context.Context, request *pb.InsertDislike) (*pb.InsertResponse, error) {
	/* sanitizacija */
	id := request.PostId
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	objectId, err := primitive.ObjectIDFromHex(id)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post with ID: " + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)

	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
			handler.CustomLogger.ErrorLogger.Error("User with ID: " + userId + " already disliked post with ID: " + id)
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
			handler.CustomLogger.InfoLogger.Info("Deleting like on post with ID: " + id + " by user with ID: " + userId)
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
		handler.CustomLogger.ErrorLogger.Error("Post with ID: " + post.Id.Hex() + " was not disliked by user with ID: " + userId)
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Post with ID: " + post.Id.Hex() + " disliked by user with ID: " + post.UserId)
	return response, err
}

func (handler *PostHandler) CommentPost(ctx context.Context, request *pb.InsertComment) (*pb.InsertResponse, error) {
	/* sanitizacija */
	id := request.PostId
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post with ID:" + objectId.Hex() + " not found")
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	success, err := handler.service.CommentPost(post, userId, strings.TrimSpace(request.Text)) // Trim - function to remove leading and trailing whitespace
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post with ID: " + post.Id.Hex() + " was not commented by user with ID: " + userId)
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Post with ID: " + post.Id.Hex() + " was commented by user with ID: " + userId)
	return response, err
}

func (handler *PostHandler) NeutralPost(ctx context.Context, request *pb.InsertNeutralReaction) (*pb.InsertResponse, error) {
	/* sanitizacija */
	id := request.PostId
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Post with ID:" + objectId.Hex() + " not found")
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
			handler.CustomLogger.InfoLogger.Info("Deleting like on post with ID: " + id + " by user with ID: " + userId)
			fmt.Println("user already dislikes selected post - neutral")
			flagDisliked = true
		}
	}

	flagLiked := false
	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
			handler.CustomLogger.InfoLogger.Info("Deleting dislike on post with ID: " + id + " by user with ID: " + userId)
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
		handler.CustomLogger.ErrorLogger.Error("Neutral reaction on post with ID: " + post.Id.Hex() + "  by user with ID: " + userId + " was not successful")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Neutral reaction on post with ID: " + post.Id.Hex() + " by user with ID: " + userId)
	return response, err
}

func (handler *PostHandler) UpdateCompanyInfo(ctx context.Context, request *pb.UpdateCompanyInfoRequest) (*pb.InsertResponse, error) {
	company, err := mapCompanyInfo(request.CompanyInfoDTO)
	oldName := request.CompanyInfoDTO.OldName
	/* sanitizacija unosa - prevencija log injection - u logove nece biti upisani specijalni karakteri */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	oldName = re.ReplaceAllString(oldName, " ")

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Company with name: " + "'" + oldName + "'" + " not found")
		return nil, err
	}

	success, err := handler.service.UpdateCompanyInfo(company, request.CompanyInfoDTO.OldName)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Company with name: " + "'" + oldName + "'" + " was not updated")
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Company with name: " + "'" + oldName + "'" + " updated")
	return response, err
}
