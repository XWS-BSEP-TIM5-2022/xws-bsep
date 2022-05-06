package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Get(id primitive.ObjectID) (*User, error)
	GetByUsername(Username string) (*User, error)
	Insert(user *User) (*User, error)
	GetAll() ([]*User, error)
	GetAllPublic() ([]*User, error)
	DeleteAll()
	Update(user *User) (string, error)
	GetById(userId string) (*User, error)
}
