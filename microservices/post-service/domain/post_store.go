package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostStore interface { // nacin za dobavljanje podataka iz skladista
	Get(ctx context.Context, id primitive.ObjectID) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	DeleteAll(ctx context.Context)
	Insert(ctx context.Context, post *Post) (string, error)
	Update(ctx context.Context, post *Post) (string, error)
	GetAllByUser(ctx context.Context, id string) ([]*Post, error)
	LikePost(ctx context.Context, post *Post, id string) (string, error)
	DislikePost(ctx context.Context, post *Post, id string) (string, error)
	CommentPost(ctx context.Context, post *Post, id string, text string) (string, error)
	UpdateCompanyInfo(ctx context.Context, company *Company, oldName string) (string, error)
}
