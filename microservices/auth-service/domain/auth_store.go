package domain

type AuthStore interface {
	Insert(user *Authentication) (string, error)
	GetAll() (*[]Authentication, error)
	DeleteAll()
}
