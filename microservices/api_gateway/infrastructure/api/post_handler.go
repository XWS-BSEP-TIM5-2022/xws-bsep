package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"time"
)

type PostHandler struct {
	postClientAddress string
	//connectionClientAddress string
}

func NewPostHandler(postClientAddress /*, connectionClientAddress*/ string) Handler {
	return &PostHandler{
		postClientAddress: postClientAddress,
		//connectionClientAddress: connectionClientAddress,
	}
}

func (handler *PostHandler) Init(mux *runtime.ServeMux) {
	fmt.Println("uslo 1")

	err := mux.HandlePath("GET", "/connection/{userID}/posts", handler.GetResult) // prikaz postova od strane zapracenog profila
	if err != nil {
		panic(err)
	}
}

func (handler *PostHandler) GetResult(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["userID"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	posts := &domain.Posts{}
	postClient := services.NewPostClient(handler.postClientAddress)

	postsByUser, err := postClient.GetAllByUser(context.TODO(), &post.GetRequest{Id: id})
	if err != nil {
		fmt.Println("desio se error")
		fmt.Println(err)
		return
	}

	for _, post := range postsByUser.Posts {
		newPost := domain.Post{
			Id:          post.Id,
			Text:        post.Text,
			Images:      post.Images,
			Links:       post.Links,
			DateCreated: post.DateCreated.AsTime(),
			UserId:      post.UserId,
		}

		for _, like := range post.Likes {
			newLike := domain.Like{
				Id:     like.Id,
				UserId: like.UserId,
			}
			newPost.Likes = append(newPost.Likes, newLike)
		}

		for _, dislike := range post.Dislikes {
			newDislike := domain.Dislike{
				Id:     dislike.Id,
				UserId: dislike.UserId,
			}
			newPost.Dislikes = append(newPost.Dislikes, newDislike)
		}

		for _, comment := range post.Comments {
			newComment := domain.Comment{
				Id:     comment.Id,
				UserId: comment.UserId,
				Text:   comment.Text,
			}
			newPost.Comments = append(newPost.Comments, newComment)
		}

		posts.AllPosts = append(posts.AllPosts, newPost) // dodati post u listu postova
	}

	response, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["userID"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	posts := &domain.Posts{}
	users := &domain.Users{}

	// getAllFriends -> lista user-a (zapraceni profili)
	// getAllPostsByUserId

	// TODO: uncomment when docker is resolved
	//err := handler.getAllConnections(users, id)
	//if err != nil {
	//	w.WriteHeader(http.StatusNotFound)
	//	return
	//}

	// TODO: ***
	var listOfUsers = []*domain.User{
		{
			Id:           "623b0cc336a1d6fd8c1cf0f6",
			Name:         "Ranko",
			LastName:     "Rankovic",
			MobileNumber: "0653829384",
			Gender:       domain.Male,
			Birthday:     time.Date(1997, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			Email:        "ranko@gmail.com",
			Biography:    "Vredan od malih nogu",
			Username:     "rankoRankovic",
		},
		{
			Id:           "623b4ac336a1d6fd8c1cf0f6",
			Name:         "Marko",
			LastName:     "Markovic",
			MobileNumber: "06538293354",
			Gender:       domain.Male,
			Birthday:     time.Date(1967, time.Month(8), 21, 1, 10, 30, 0, time.UTC),
			Email:        "markic@gmail.com",
			Biography:    "Rodjen u Novom Sadu",
			Username:     "markooom",
		},
	}

	for _, u := range listOfUsers {
		newUser := domain.User{
			Id: u.Id,
		}
		users.UsersDetails = append(users.UsersDetails, newUser)
	}
	// TODO: ***

	err := handler.addPosts(posts, users)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

//func (handler *PostHandler) getAllConnections(users *domain.Users, userId string) error {
//	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
//	conections, err := connectionClient.GetFriends(context.TODO(), &connection.GetRequest{UserID: userId})
//	if err != nil {
//		return err
//	}
//
//	for _, user := range conections.Users {
//		newUser := domain.User{
//			Id: user.UserID,
//		}
//		users.UsersDetails = append(users.UsersDetails, newUser)
//	}
//	return nil
//}

func (handler *PostHandler) addPosts(posts *domain.Posts, users *domain.Users) error {
	postClient := services.NewPostClient(handler.postClientAddress)

	for _, user := range users.UsersDetails {
		postsByUser, err := postClient.GetAllByUser(context.TODO(), &post.GetRequest{Id: user.Id})
		if err != nil {
			fmt.Println("desio se error!")
			return err
		}

		for _, post := range postsByUser.Posts {
			newPost := domain.Post{
				Id:          post.Id,
				Text:        post.Text,
				Images:      post.Images,
				Links:       post.Links,
				DateCreated: post.DateCreated.AsTime(),
				UserId:      post.UserId,
			}

			for _, like := range post.Likes {
				newLike := domain.Like{
					Id:     like.Id,
					UserId: like.UserId,
				}
				newPost.Likes = append(newPost.Likes, newLike)
			}

			for _, dislike := range post.Dislikes {
				newDislike := domain.Dislike{
					Id:     dislike.Id,
					UserId: dislike.UserId,
				}
				newPost.Dislikes = append(newPost.Dislikes, newDislike)
			}

			for _, comment := range post.Comments {
				newComment := domain.Comment{
					Id:     comment.Id,
					UserId: comment.UserId,
					Text:   comment.Text,
				}
				newPost.Comments = append(newPost.Comments, newComment)
			}

			posts.AllPosts = append(posts.AllPosts, newPost) // dodati post u listu postova
		}
	}
	return nil
}
