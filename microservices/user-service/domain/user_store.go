package domain

type UserStore interface {
	Get(id string) (*User, error)
	Insert(user *User) (string, error)
	GetAll() (*[]User, error)
	DeleteAll()
	Update(user *User) (string, error)
}
