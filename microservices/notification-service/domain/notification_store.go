package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationStore interface {
	GetById(ctx context.Context, id primitive.ObjectID) (*Notification, error)
	GetAll(ctx context.Context) ([]*Notification, error)
	Insert(ctx context.Context, post *Notification) (string, error)
	GetAllByUser(ctx context.Context, id string) ([]*Notification, error)
}
