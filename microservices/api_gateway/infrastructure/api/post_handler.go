package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"html"
	"net/http"
)

type PostHandler struct {
	postClientAddress       string
	connectionClientAddress string
	userClientAddress       string
}

func NewPostHandler(postClientAddress, connectionClientAddress, userClientAddress string) Handler {
	return &PostHandler{
		postClientAddress:       postClientAddress,
		connectionClientAddress: connectionClientAddress,
		userClientAddress:       userClientAddress,
	}
}

func (handler *PostHandler) Init(mux *runtime.ServeMux) {
	fmt.Println("Hello from api gateway")

	err := mux.HandlePath("GET", "/api/feed/{userID}", handler.GetPosts) // prikaz postova od strane zapracenog profila
	if err != nil {
		panic(err)
	}

	err = mux.HandlePath("GET", "/api/feed/public", handler.GetPublicPosts)
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

	if len(id) != 24 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	/** Escape '$' - Prevent NoSQL Injection **/
	var checkId = ""
	for i := 0; i < len(id); i++ {
		char := string(id[i])
		if char != "$" {
			checkId = checkId + char
		}
	}

	posts := &domain.Posts{}
	users := &domain.Users{}

	err := handler.getAllConnections(users, html.EscapeString(checkId)) /** EscapeString **/
	if err != nil {
		fmt.Println("error 1")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = handler.addPosts(posts, users)
	if err != nil {
		fmt.Println("error 2")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("error 3")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Authorization", "Bearer " )	// TODO
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) GetPublicPosts(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	postClient := services.NewPostClient(handler.postClientAddress)
	posts, err := postClient.GetAll(context.TODO(), &post.GetAllRequest{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	allPosts := &domain.Posts{}

	for _, post := range posts.Posts {
		isPublic, err := handler.isUserPublic(post.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		if isPublic {
			if post.IsJobOffer {
				newPost := domain.Post{
					Id:          post.Id,
					Text:        post.Text,
					Images:      post.Images,
					Links:       post.Links,
					DateCreated: post.DateCreated.AsTime(),
					UserId:      post.UserId,
					JobOffer: domain.JobOffer{
						//Id: string(post.JobOffer.Id),
						JobDescription:  post.JobOffer.JobDescription,
						Preconditions:   post.JobOffer.Preconditions,
						DailyActivities: post.JobOffer.DailyActivities,
						Position: domain.Position{
							//Id:
							Name: post.JobOffer.Position.Name,
							Pay:  post.JobOffer.Position.Pay,
						},
					},
					IsJobOffer: post.IsJobOffer,
					Company: domain.Company{
						Name:        post.Company.Name,
						Description: post.Company.Description,
						PhoneNumber: post.Company.PhoneNumber,
						IsActive:    post.Company.IsActive,
					},
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
				allPosts.AllPosts = append(allPosts.AllPosts, newPost)
			} else {
				newPost := domain.Post{
					Id:          post.Id,
					Text:        post.Text,
					Images:      post.Images,
					Links:       post.Links,
					DateCreated: post.DateCreated.AsTime(),
					UserId:      post.UserId,
					IsJobOffer:  post.IsJobOffer,
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

				allPosts.AllPosts = append(allPosts.AllPosts, newPost)
			}
		}
	}

	response, err := json.Marshal(allPosts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) getAllConnections(users *domain.Users, userId string) error {
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	connections, err := connectionClient.GetConnections(context.TODO(), &connection.GetRequest{UserID: userId})
	if err != nil {
		return err
	}

	fmt.Println("konekcije", connections) // TODO greska
	for _, user := range connections.Users {
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

		fmt.Println("77777777777777777")
		// TODO: ovo je greska

		postsByUser, err := postClient.GetAllByUser(context.TODO(), &post.GetRequest{Id: user.Id})

		fmt.Println("ove je uslo jjjjjjjjjjjjeeeej", postsByUser)
		if err != nil {
			fmt.Println("desio se error!")
			return err
		}
		for _, post := range postsByUser.Posts {
			if post.IsJobOffer {
				newPost := domain.Post{
					Id:          post.Id,
					Text:        post.Text,
					Images:      post.Images,
					Links:       post.Links,
					DateCreated: post.DateCreated.AsTime(),
					UserId:      post.UserId,
					JobOffer: domain.JobOffer{
						//Id: string(post.JobOffer.Id),
						JobDescription:  post.JobOffer.JobDescription,
						Preconditions:   post.JobOffer.Preconditions,
						DailyActivities: post.JobOffer.DailyActivities,
						Position: domain.Position{
							//Id:
							Name: post.JobOffer.Position.Name,
							Pay:  post.JobOffer.Position.Pay,
						},
					},
					IsJobOffer: post.IsJobOffer,
					Company: domain.Company{
						Name:        post.Company.Name,
						Description: post.Company.Description,
						PhoneNumber: post.Company.PhoneNumber,
						IsActive:    post.Company.IsActive,
					},
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
			} else {
				newPost := domain.Post{
					Id:          post.Id,
					Text:        post.Text,
					Images:      post.Images,
					Links:       post.Links,
					DateCreated: post.DateCreated.AsTime(),
					UserId:      post.UserId,
					IsJobOffer:  post.IsJobOffer,
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

				posts.AllPosts = append(posts.AllPosts, newPost)
			}
		}
	}
	return nil
}

func (handler *PostHandler) isUserPublic(id string) (bool, error) {
	userClient := services.NewUserClient(handler.userClientAddress)
	users, err := userClient.GetAllPublic(context.TODO(), &user.GetAllPublicRequest{})
	if err != nil {
		fmt.Println("isUserPublic vratilo gresku")
		return false, err
	}

	for _, user := range users.Users {
		if user.Id == id {
			return true, nil
		}
	}
	return false, nil
}
