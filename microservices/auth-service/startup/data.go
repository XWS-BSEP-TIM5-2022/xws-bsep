package startup

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
)

var auths = []*domain.Authentication{
	// {
	// 	Id:               "62778fe0042817b7882ee522",
	// 	Username:         "rankoRankovic",
	// 	Password:         "$2a$10$VN.d.CATdxWu/Kbqv3JQnuq3eG.NyUnqWqOV.VM.V.xvZKrzGR8tK",
	// 	Role:             "User",
	// 	VerificationCode: "0",
	// 	ExpirationTime:   0,
	// },
	//{
	//	Id:       "623b4ac336a1d6fd8c1cf0f6",
	//	Username: "markooom",
	//	Password: "$2a$10$VN.d.CATdxWu/Kbqv3JQnuq3eG.NyUnqWqOV.VM.V.xvZKrzGR8tK",
	//	Role:     "User",
	//},
}

var roles = []*domain.Role{

	{
		ID:          1, //uint(getId()),
		Name:        "User",
		Permissions: userPermissions,
	},
	{
		ID:          2,
		Name:        "Admin",
		Permissions: adminPermissions,
	},
}

var adminPermissions = []*domain.Permission{
	{
		ID:   15,
		Name: "AdminsEndpoint",
	},
}

var userPermissions = []*domain.Permission{
	{ // auth service
		ID:   1,
		Name: "UpdateUsername",
	},
	{
		ID:   2,
		Name: "UpdatePassword",
	},
	{ // user-service
		ID:   3,
		Name: "UpdateUserProfile",
	},
	{
		ID:   4,
		Name: "GetLoggedInUserInfo",
	},
	{
		ID:   5,
		Name: "GetUserById",
	},
	{
		ID:   6,
		Name: "GetAllUsers",
	},
	{ // post service
		ID:   7,
		Name: "GetPostById",
	},
	{
		ID:   8,
		Name: "GetAllPosts",
	},
	{
		ID:   9,
		Name: "UpdatePost",
	},
	{
		ID:   10,
		Name: "CreatePost",
	},
	{
		ID:   11,
		Name: "UpdatePostLikes",
	},
	{
		ID:   12,
		Name: "UpdatePostDislikes",
	},
	{
		ID:   13,
		Name: "UpdatePostComments",
	},
	{
		ID:   14,
		Name: "NeutralPost",
	},
	{ // connection service
		ID:   16,
		Name: "RegisterConnection",
	},
	{
		ID:   17,
		Name: "CreateConnection",
	},
	{
		ID:   18,
		Name: "RejectConnection",
	},
	{
		ID:   19,
		Name: "ApproveConnection",
	},
	{
		ID:   20,
		Name: "CheckConnection",
	},
	{
		ID:   21,
		Name: "BlockUser",
	},
	{
		ID:   22,
		Name: "CreateMessage",
	},
	{
		ID:   23,
		Name: "ChangePrivacy",
	},
}
