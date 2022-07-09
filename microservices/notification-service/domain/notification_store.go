package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationStore interface {
	GetById(id primitive.ObjectID) (*Notification, error)
}
