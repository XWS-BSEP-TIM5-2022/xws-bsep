package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/square/go-jose.v2/jwt"
	"html"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type PostHandler struct {
	postClientAddress       string
	connectionClientAddress string
	userClientAddress       string
	CustomLogger            *CustomLogger
}

func NewPostHandler(postClientAddress, connectionClientAddress, userClientAddress string) Handler {
	CustomLogger := NewCustomLogger()
	return &PostHandler{
		postClientAddress:       postClientAddress,
		connectionClientAddress: connectionClientAddress,
		userClientAddress:       userClientAddress,
		CustomLogger:            CustomLogger,
	}
}

func (handler *PostHandler) Init(mux *runtime.ServeMux) {
	fmt.Println("Hello from api gateway")

	err := mux.HandlePath("GET", "/api/feed/{userID}", handler.GetPosts) // prikaz postova od strane zapracenog profila
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Feed not found")
		panic(err)
	}

	err = mux.HandlePath("GET", "/api/feed/public", handler.GetPublicPosts)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Feed for unregistered user not found")
		panic(err)
	}
}

func (handler *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	/* sanitizacija - prevencija log injection-a */
	id := pathParams["userID"]
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	if id == "" {
		handler.CustomLogger.ErrorLogger.Error("User with ID: " + id + " is non-existent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(id) != 24 {
		handler.CustomLogger.ErrorLogger.Error("User with ID: " + id + " is non-existent")
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

	err = handler.AuthorizeUser(w, r, id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with can not access feed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	handler.CustomLogger.SuccessLogger.Info("User with ID: " + id + " can access feed")

	posts := &domain.Posts{}
	users := &domain.Users{}

	err = handler.getAllConnections(users, html.EscapeString(checkId)) /** EscapeString **/
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all connections for user with ID: " + id + " unsuccessful")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(users.UsersDetails)) + " connections for user with ID: " + id)

	err = handler.addPosts(posts, users)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get feed for user with ID: " + id + " unsuccessful")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Marshal posts is unsuccessful")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(posts.AllPosts)) + " posts for feed")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) AuthorizeUser(w http.ResponseWriter, r *http.Request, requestId string) error {
	jwtToken := r.Header.Get("Authorization")
	jwtToken = jwtToken[7:]
	var claims map[string]interface{}
	token, _ := jwt.ParseSigned(jwtToken)
	_ = token.UnsafeClaimsWithoutVerification(&claims)
	username := claims["username"]

	userClient := services.NewUserClient(handler.userClientAddress)
	idByUsername, err := userClient.GetIdByUsername(context.TODO(), &user.GetIdByUsernameRequest{Username: fmt.Sprint(username)})
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Can not find ID of user with name: " + fmt.Sprint(username))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return err
	}
	if idByUsername.Id != requestId {
		handler.CustomLogger.ErrorLogger.Error("User with ID: " + idByUsername.Id + " can not access feed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		return err
	}
	return nil
}

func (handler *PostHandler) GetPublicPosts(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	postClient := services.NewPostClient(handler.postClientAddress)
	posts, err := postClient.GetAll(context.TODO(), &post.GetAllRequest{})
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all public posts unsuccessful")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	allPosts := &domain.Posts{}
	for _, post := range posts.Posts {
		isPublic, err := handler.isUserPublic(post.UserId)
		if err != nil {
			handler.CustomLogger.ErrorLogger.Error("Is user public unsuccessful")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if isPublic {
			if post.IsJobOffer {
				newPost := domain.Post{
					Id:          post.Id,
					Text:        post.Text,
					Image:       post.Image,
					Links:       post.Links,
					DateCreated: post.DateCreated.AsTime(),
					UserId:      post.UserId,
					JobOffer: domain.JobOffer{
						Id:              post.JobOffer.Id,
						JobDescription:  post.JobOffer.JobDescription,
						Preconditions:   post.JobOffer.Preconditions,
						DailyActivities: post.JobOffer.DailyActivities,
						Position: domain.Position{
							Id:   post.JobOffer.Position.Id,
							Name: post.JobOffer.Position.Name,
							Pay:  post.JobOffer.Position.Pay,
						},
					},
					IsJobOffer: post.IsJobOffer,
					Company: domain.Company{
						Id:          post.Company.Id,
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
					Image:       post.Image,
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
		handler.CustomLogger.ErrorLogger.Error("Marshal public posts is unsuccessful")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(posts.Posts)) + " public posts")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) getAllConnections(users *domain.Users, userId string) error {
	// sanitizacija userId uradjena pre poziva same fje
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	connections, err := connectionClient.GetConnections(context.TODO(), &connection.GetRequest{UserID: userId})
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get connections for user with ID: " + userId + " unsuccessful")
		return err
	}

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
		postsByUser, err := postClient.GetAllByUser(context.TODO(), &post.GetRequest{Id: user.Id})
		if err != nil {
			handler.CustomLogger.ErrorLogger.Error("Get posts by user with ID: " + user.Id + " unsuccessful")
			return err
		}
		for _, post := range postsByUser.Posts {
			if post.IsJobOffer {
				newPost := domain.Post{
					Id:          post.Id,
					Text:        post.Text,
					Image:       post.Image,
					Links:       post.Links,
					DateCreated: post.DateCreated.AsTime(),
					UserId:      post.UserId,
					JobOffer: domain.JobOffer{
						Id:              post.JobOffer.Id,
						JobDescription:  post.JobOffer.JobDescription,
						Preconditions:   post.JobOffer.Preconditions,
						DailyActivities: post.JobOffer.DailyActivities,
						Position: domain.Position{
							Id:   post.JobOffer.Position.Id,
							Name: post.JobOffer.Position.Name,
							Pay:  post.JobOffer.Position.Pay,
						},
					},
					IsJobOffer: post.IsJobOffer,
					Company: domain.Company{
						Id:          post.Company.Id,
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
					Image:       post.Image,
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
		handler.CustomLogger.ErrorLogger.Error("Get all public users unsuccessful")
		return false, err
	}

	for _, user := range users.Users {
		if user.Id == id {
			return true, nil
		}
	}
	return false, nil
}
