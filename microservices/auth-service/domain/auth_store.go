package domain

type AuthStore interface {
	Create(auth *Authentication) (*Authentication, error)
	FindByUsername(username string) (*Authentication, error)
	FindAll() (*[]Authentication, error)
}
