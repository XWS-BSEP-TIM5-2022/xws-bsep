package application

import (
	"context"
	auth "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	store             domain.PostStore
	userServiceClient user.UserServiceClient
	authServiceClient auth.AuthServiceClient
}

func NewPostService(store domain.PostStore, userServiceClient user.UserServiceClient, authServiceClient auth.AuthServiceClient) *PostService {
	return &PostService{
		store:             store,
		userServiceClient: userServiceClient,
		authServiceClient: authServiceClient,
	}
}

func (service *PostService) Get(id primitive.ObjectID) (*domain.Post, error) {
	return service.store.Get(id)
}

func (service *PostService) GetAll() ([]*domain.Post, error) {
	return service.store.GetAll()
}

func (service *PostService) Insert(post *domain.Post) (string, error) {
	success, err := service.store.Insert(post)
	return success, err
}

func (service *PostService) Update(post *domain.Post) (string, error) {
	success, err := service.store.Update(post)
	return success, err
}

func (service *PostService) GetAllByUser(id string) ([]*domain.Post, error) {
	return service.store.GetAllByUser(id)
}

func (service *PostService) LikePost(post *domain.Post, id string) (string, error) {
	return service.store.LikePost(post, id)
}

func (service *PostService) DislikePost(post *domain.Post, id string) (string, error) {
	return service.store.DislikePost(post, id)
}

func (service *PostService) CommentPost(post *domain.Post, id string, text string) (string, error) {
	return service.store.CommentPost(post, id, text)
}

func (service *PostService) GetUsernameByApiToken(ctx context.Context, apiToken string) (*auth.GetUsernameResponse, error) {
	return service.authServiceClient.GetUsernameByApiToken(ctx, &auth.GetUsernameRequest{
		ApiToken: apiToken,
	})
}

func (service *PostService) GetIdByUsername(ctx context.Context, username string) (*user.InsertResponse, error) {
	return service.userServiceClient.GetIdByUsername(ctx, &user.GetIdByUsernameRequest{
		Username: username,
	})
}

func (service *PostService) UpdateCompanyInfo(company *domain.Company, oldName string) (string, error) {
	return service.store.UpdateCompanyInfo(company, oldName)
}
