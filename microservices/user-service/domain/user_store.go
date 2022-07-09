package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByUsername(ctx context.Context, Username string) (*User, error)
	GetByEmail(Email string) (*User, error)
	Insert(user *User) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	GetAllPublic(ctx context.Context) ([]*User, error)
	DeleteAll()
	Update(ctx context.Context, user *User) (string, error)
	UpdateBasicInfo(ctx context.Context, user *User) (string, error)
	UpdateExperienceAndEducation(ctx context.Context, user *User) (string, error)
	UpdateSkillsAndInterests(ctx context.Context, user *User) (string, error)
	GetById(ctx context.Context, userId string) (*User, error)
	Search(ctx context.Context, criteria string) ([]*User, error)
	UpdateIsActiveById(ctx context.Context, userId string) error
	GetIdByEmail(ctx context.Context, email string) (string, error)
	DeleteUser(userId, email string) error
}
