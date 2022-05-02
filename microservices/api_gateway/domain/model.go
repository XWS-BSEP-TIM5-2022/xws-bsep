package domain

type Users struct {
	UsersDetails []User
}

type User struct {
	Id   string
	Name string
	// Email string
}
