package domain

type AuthStore interface {
	Create(auth *Authentication) (string, error)
	GetAll() (*[]Authentication, error)
	DeleteAll()
	Insert(auth *Authentication) (string, error)
}
