package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostStore interface { // nacin za dobavljanje podataka iz skladista
	Get(id primitive.ObjectID) (*Post, error)
	GetAll() ([]*Post, error)
	DeleteAll()
	Insert(post *Post) (string, error)
	Update(post *Post) (string, error)
	GetAllByUser(string) ([]*Post, error)
	LikePost(post *Post, id string) (string, error)
}
