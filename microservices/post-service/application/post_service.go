package application

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	store domain.PostStore
}

func NewPostService(store domain.PostStore) *PostService {
	return &PostService{
		store: store,
	}
}

func (service *PostService) Get(ctx context.Context, id primitive.ObjectID) (*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "Get service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.Get(ctx, id)
}

func (service *PostService) GetAll(ctx context.Context) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAll(ctx)
}

func (service *PostService) Insert(ctx context.Context, post *domain.Post) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	success, err := service.store.Insert(ctx, post)
	return success, err
}

func (service *PostService) Update(ctx context.Context, post *domain.Post) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Update service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	success, err := service.store.Update(ctx, post)
	return success, err
}

func (service *PostService) GetAllByUser(ctx context.Context, id string) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUser service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllByUser(ctx, id)
}

func (service *PostService) LikePost(ctx context.Context, post *domain.Post, id string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "LikePost service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.LikePost(ctx, post, id)
}

func (service *PostService) DislikePost(ctx context.Context, post *domain.Post, id string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "DislikePost service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.DislikePost(ctx, post, id)
}

func (service *PostService) CommentPost(ctx context.Context, post *domain.Post, id string, text string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "CommentPost service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.CommentPost(ctx, post, id, text)
}

func (service *PostService) UpdateCompanyInfo(ctx context.Context, company *domain.Company, oldName string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateCompanyInfo service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.UpdateCompanyInfo(ctx, company, oldName)
}
