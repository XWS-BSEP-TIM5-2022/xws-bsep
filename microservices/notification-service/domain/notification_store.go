package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationStore interface {
	GetById(id primitive.ObjectID) (*Notification, error)
	GetAll() ([]*Notification, error)
	Insert(post *Notification) (string, error)
	GetAllByUser(string) ([]*Notification, error)
}
