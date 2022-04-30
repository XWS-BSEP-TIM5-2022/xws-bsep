package domain

type UserStore interface {
	Insert(product *User) error
	GetAll() (*[]User, error)
}
