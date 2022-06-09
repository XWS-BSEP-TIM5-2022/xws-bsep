package create_user

import "time"

type UserDetails struct {
	Id           string
	Name         string
	LastName     string
	Email        string
	MobileNumber string
	Gender       GenderEnum
	Birthday     time.Time
	Username     string
	Biography    string
	// IsPublic     bool
	Education  []Education
	Experience []Experience
	Skills     []Skill
	Interests  []Interest
	Role       []string
	// IsActive     bool
	Password string
}

type CreateUserCommandType int8

const (
	CreateUser CreateUserCommandType = iota
	DeleteUser
	CreateAuth
	RollbackUser
	ApproveUser
	UnknownCommand
)

type CreateUserCommand struct {
	User UserDetails
	Type CreateUserCommandType
}

type CreateUserReplyType int8

const (
	UserCreated CreateUserReplyType = iota
	UserNotCreated
	AuthCreated
	AuthNotCreated
	UserRolledBack
	UserDeleted
	UserApproved
	UnknownReply
)

type CreateUserReply struct {
	User UserDetails
	Type CreateUserReplyType
}

type Education struct {
	Id        string //primitive.ObjectID
	Name      string
	Level     EducationEnum
	Place     string
	StartDate time.Time
	EndDate   time.Time
}

type Experience struct {
	Id        string //primitive.ObjectID
	Name      string
	Headline  string
	Place     string
	StartDate time.Time
	EndDate   time.Time
}

type Skill struct {
	Id   string //primitive.ObjectID
	Name string
}

type Interest struct {
	Id          string //primitive.ObjectID
	Name        string
	Description string
}

type GenderEnum int8
type EducationEnum int8

const (
	Female GenderEnum = iota
	Male
)

const (
	Primary EducationEnum = iota
	Secondary
	Bachelor
	Master
	Doctorate
)

type Role struct {
	ID          uint
	Name        string
	Permissions []*Permission
}

type Permission struct {
	ID   uint
	Name string
}
