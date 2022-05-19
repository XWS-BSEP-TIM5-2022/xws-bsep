package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type PostHandler struct {
	postClientAddress       string
	connectionClientAddress string
}

func NewPostHandler(postClientAddress, connectionClientAddress string) Handler {
	return &PostHandler{
		postClientAddress:       postClientAddress,
		connectionClientAddress: connectionClientAddress,
	}
}

func (handler *PostHandler) Init(mux *runtime.ServeMux) {
	fmt.Println("uslo 1")

	err := mux.HandlePath("GET", "/api/feed/{userID}", handler.GetPosts) // prikaz postova od strane zapracenog profila
	if err != nil {
		panic(err)
	}
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

	err := handler.getAllConnections(users, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = handler.addPosts(posts, users)
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

func (handler *PostHandler) getAllConnections(users *domain.Users, userId string) error {
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	conections, err := connectionClient.GetConnections(context.TODO(), &connection.GetRequest{UserID: userId})
	if err != nil {
		return err
	}

	for _, user := range conections.Users {
		newUser := domain.User{
			Id: user.UserID,
		}
		users.UsersDetails = append(users.UsersDetails, newUser)
	}
	return nil
}

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
