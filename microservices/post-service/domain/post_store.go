package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostStore interface { // nacin za dobavljanje podataka iz skladista
	Get(id primitive.ObjectID) (*Post, error)
	GetAll() ([]*Post, error)
	DeleteAll()
	Insert(user *Post) (string, error)
	Update(user *Post) (string, error)
}
