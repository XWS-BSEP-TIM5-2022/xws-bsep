package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventStore interface {
	GetById(id primitive.ObjectID) (*Event, error)
}
