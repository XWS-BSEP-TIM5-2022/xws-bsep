package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gopkg.in/square/go-jose.v2/jwt"
	"html"
	"log"
	"net/http"
	"regexp"
)

type JobOfferHandler struct {
	postClientAddress       string
	connectionClientAddress string
	userClientAddress       string
	authClientAddress       string
	jobOfferAddress         string
	CustomLogger            *CustomLogger
}

func NewJobOfferHandler(postClientAddress, connectionClientAddress, userClientAddress string, authClientAddress string, jobOfferAddress string) Handler {
	CustomLogger := NewCustomLogger()
	return &JobOfferHandler{
		postClientAddress:       postClientAddress,
		connectionClientAddress: connectionClientAddress,
		userClientAddress:       userClientAddress,
		authClientAddress:       authClientAddress,
		CustomLogger:            CustomLogger,
		jobOfferAddress:         jobOfferAddress,
	}
}

func (handler *JobOfferHandler) Init(mux *runtime.ServeMux) {
	fmt.Println("Hello from api gateway")

	err := mux.HandlePath("GET", "/api/jobOfferRecommendations/{userID}", handler.GetRecommendations)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Feed not found")
		panic(err)
	}
}

func (handler *JobOfferHandler) GetRecommendations(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

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

	result := handler.AuthorizeUser(w, r, id)
	if result == "error" {
		handler.CustomLogger.ErrorLogger.Error("Access to feed denied")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	handler.CustomLogger.SuccessLogger.Info("User with ID: " + id + " can access feed")

	user := &domain.User{}
	posts := &domain.Posts{}
	err = handler.getUser(user, html.EscapeString(checkId)) /** EscapeString **/
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + id + " not found")
		return
	}

	err = handler.findJobOffers(user, posts)

}

func (handler *JobOfferHandler) AuthorizeUser(w http.ResponseWriter, r *http.Request, requestId string) string {
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
		return "error"
	}
	if idByUsername.Id != requestId {
		handler.CustomLogger.ErrorLogger.Error("User with ID: " + idByUsername.Id + " can not access feed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		return "error"
	}
	return "success"
}

func (handler *JobOfferHandler) getUser(loggedUser *domain.User, userId string) error {
	userClient := services.NewUserClient(handler.userClientAddress)
	foundUser, err := userClient.Get(context.TODO(), &user.GetRequest{Id: userId})
	fmt.Println(err)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get user with ID: " + userId + " unsuccessful")
		return err
	}
	fmt.Println(foundUser)

	newUser := domain.User{
		Id: foundUser.User.Id,
	}

	for _, expirience := range foundUser.User.Experience {
		newExp := domain.Experience{
			Id:       expirience.Id,
			Headline: expirience.Headline,
		}
		newUser.Experience = append(newUser.Experience, newExp)
	}

	for _, skill := range foundUser.User.Skills {
		newSkill := domain.Skill{
			Id:   skill.Id,
			Name: skill.Name,
		}
		newUser.Skills = append(newUser.Skills, newSkill)
	}

	loggedUser.Id = newUser.Id
	loggedUser.Experience = newUser.Experience
	loggedUser.Skills = newUser.Skills

	handler.CustomLogger.ErrorLogger.Error("Pronadjen uspjesno")
	handler.CustomLogger.ErrorLogger.Error(loggedUser)

	return nil
}

func (handler *JobOfferHandler) findJobOffers(loggedUser *domain.User, posts *domain.Posts) error {

	postClient := services.NewPostClient(handler.postClientAddress)
	allPosts, err := postClient.GetAll(context.TODO(), &post.GetAllRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all posts unsuccessful")
		return err
	}

	for _, post := range allPosts.Posts {
		/*
			Ako je post ponuda za posao i ako je user public, prolazim kroz sve skills na postu i sve skills od usera i ako su jednaki, dodajem post u listu buducih preporuka
		*/
		if post.IsJobOffer {
			userClient := services.NewUserClient(handler.userClientAddress)
			foundUser, _ := userClient.Get(context.TODO(), &user.GetRequest{Id: loggedUser.Id})

			if foundUser.User.IsPublic {
				for _, skill := range loggedUser.Skills {
					if skill.Name == post.JobOffer.Preconditions {
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
						posts.AllPosts = append(posts.AllPosts, newPost)
					}
				}

				for _, exp := range loggedUser.Experience {
					if exp.Headline == post.JobOffer.Position.Name {
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
						posts.AllPosts = append(posts.AllPosts, newPost)
					}
				}
			}
		}
	}

	return nil
}
