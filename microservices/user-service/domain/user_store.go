package domain

type UserStore interface {
	Insert(user *User) (string, error)
	GetAll() (*[]User, error)
	DeleteAll()
}
