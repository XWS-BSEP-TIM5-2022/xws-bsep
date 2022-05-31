package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

// implementacije gRPC servera koji smo definisali u okviru common paketa

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service *application.PostService
}

func NewPostHandler(service *application.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (handler *PostHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	postPb := mapPost(post) // prepakujemo iz domenskog modela u protobuf oblik
	response := &pb.GetResponse{
		Post: postPb,
	}
	return response, nil
}

func (handler *PostHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	posts, err := handler.service.GetAll()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Posts: []*pb.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

func (handler *PostHandler) GetAllByUser(ctx context.Context, request *pb.GetRequest) (*pb.GetAllResponse, error) {
	id := request.Id
	posts, err := handler.service.GetAllByUser(id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Posts: []*pb.Post{},
	}
	for _, post := range posts {
		current := mapPost(post)
		response.Posts = append(response.Posts, current)
	}
	return response, nil
}

func (handler *PostHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	post, err := mapInsertPost(request.InsertPost)
	if err != nil {
		return nil, err
	}

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	post.UserId = userId
	success, err := handler.service.Insert(post)
	if err != nil {
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

func (handler *PostHandler) InsertJobOffer(ctx context.Context, request *pb.InsertJobOfferRequest) (*pb.InsertResponse, error) {
	//api_token := request.InsertJobOfferPost.ApiToken;

	// iz auth servisa na osnovu tokena izvuci usera ili username
	// iz user servisa na osnovu username izvuci usera
	// kreirati post u ime tog izvucenog usera

	// 6295f438b2b4c212b7c369e4 user_id
	jobOffer, err := mapInsertJobOffer(request.InsertJobOfferPost.JobOffer)
	if err != nil {
		return nil, err
	}

	position, err := mapInsertPosition(request.InsertJobOfferPost.JobOffer.Position)
	if err != nil {
		return nil, err
	}

	company, err := mapInsertCompany(request.InsertJobOfferPost.JobOffer.Company)
	if err != nil {
		return nil, err
	}

	jobOffer.Position.Id = position.Id
	jobOffer.Position.Name = position.Name
	jobOffer.Position.Pay = position.Pay
	jobOffer.Company.Id = company.Id
	jobOffer.Company.Name = company.Name
	jobOffer.Company.Description = company.Description
	jobOffer.Company.PhoneNumber = company.PhoneNumber
	jobOffer.Company.IsActive = company.IsActive

	post, err := mapInsertJobOfferPost(request.InsertJobOfferPost)
	if err != nil {
		return nil, err
	}

	post.JobOffer.Id = jobOffer.Id
	post.JobOffer.Preconditions = jobOffer.Preconditions
	post.JobOffer.JobDescription = jobOffer.JobDescription
	post.JobOffer.DailyActivities = jobOffer.DailyActivities
	post.JobOffer.Position = jobOffer.Position
	post.JobOffer.Company = jobOffer.Company

	post.UserId = "6295f438b2b4c212b7c369e4" // TODO: change

	success, err := handler.service.Insert(post)
	if err != nil {
		return nil, err
	}
	fmt.Println("ovo je request:", request)
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, nil
}

// TODO: delete
func (handler *PostHandler) Test(ctx context.Context, request *pb.TestRequest) (*pb.InsertResponse, error) {
	fmt.Println("ovo je request:", request)
	response := &pb.InsertResponse{
		Success: "hhh",
	}
	return response, nil
}

func (handler *PostHandler) LikePost(ctx context.Context, request *pb.InsertLike) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)

	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
			fmt.Println("user already likes selected post")
			return &pb.InsertResponse{
				Success: "error",
			}, err
		}
	}

	flag := false
	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
			fmt.Println("user disliked selected post, deleting dislike")
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
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

func (handler *PostHandler) DislikePost(ctx context.Context, request *pb.InsertDislike) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	postHelper, err := handler.service.Get(objectId)
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)

	// provera - da li je korisnik vec dislajkovao post
	for _, p := range post.Dislikes {
		if p.UserId == userId {
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
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

func (handler *PostHandler) CommentPost(ctx context.Context, request *pb.InsertComment) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
		return &pb.InsertResponse{
			Success: "error",
		}, err
	}

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	success, err := handler.service.CommentPost(post, userId, strings.TrimSpace(request.Text)) // Trim - function to remove leading and trailing whitespace
	if err != nil {
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

func (handler *PostHandler) NeutralPost(ctx context.Context, request *pb.InsertNeutralReaction) (*pb.InsertResponse, error) {
	id := request.PostId
	objectId, err := primitive.ObjectIDFromHex(id)
	post, err := handler.service.Get(objectId)
	if err != nil {
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
			fmt.Println("user already dislikes selected post - neutral")
			flagDisliked = true
		}
	}

	flagLiked := false
	// provera - da li je korisnik vec lajkovao post
	for _, p := range post.Likes {
		if p.UserId == userId {
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
		return nil, err
	}
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}
