package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	jobOffer "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/square/go-jose.v2/jwt"
	"html"
	"log"
	"net/http"
	"regexp"
	"strings"
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

	allJobs := handler.findRecommendations(user, posts)

	response, err := json.Marshal(allJobs)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Marshal posts is unsuccessful")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	fmt.Println(allJobs)
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
		Id:       foundUser.User.Id,
		Email:    foundUser.User.Email,
		Name:     foundUser.User.Name,
		LastName: foundUser.User.LastName,
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
	loggedUser.Name = newUser.Name
	loggedUser.LastName = foundUser.User.LastName
	loggedUser.Email = newUser.Email
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

		//TODO: nadji poslove od konekcija
		//TODO: provjera da ne smije da doda iste postove 2 puta, za skill i exp
		if post.IsJobOffer {
			userClient := services.NewUserClient(handler.userClientAddress)
			foundUser, _ := userClient.Get(context.TODO(), &user.GetRequest{Id: loggedUser.Id})

			if foundUser.User.IsPublic {
				for _, skill := range loggedUser.Skills {
					if strings.ToUpper(skill.Name) == strings.ToUpper(post.JobOffer.Preconditions) || strings.Contains(strings.ToUpper(skill.Name), strings.ToUpper(post.JobOffer.Preconditions)) {
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
					if strings.ToUpper(exp.Headline) == strings.ToUpper(post.JobOffer.Position.Name) || strings.Contains(strings.ToUpper(exp.Headline), strings.ToUpper(post.JobOffer.Position.Name)) {
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

func (handler *JobOfferHandler) findRecommendations(loggedUser *domain.User, posts *domain.Posts) *jobOffer.Recommendations {
	jobOfferClient := services.NewJobOfferClient(handler.jobOfferAddress)

	user := &jobOffer.User{
		Id:       loggedUser.Id,
		Name:     loggedUser.Name,
		LastName: loggedUser.LastName,
		Email:    loggedUser.Email,
	}

	for _, exp := range loggedUser.Experience {
		newExp := &jobOffer.Experience{
			Id:       exp.Id,
			Headline: exp.Headline,
		}
		user.Experience = append(user.Experience, newExp)
	}

	for _, skill := range loggedUser.Skills {
		newSkill := &jobOffer.Skill{
			Id:   skill.Id,
			Name: skill.Name,
		}
		user.Skills = append(user.Skills, newSkill)
	}

	response := &jobOffer.JobOffers{
		JobOffers: []*jobOffer.Post{},
	}

	for _, post := range posts.AllPosts {
		newPost := &jobOffer.Post{
			Id:          post.Id,
			Text:        post.Text,
			Image:       post.Image,
			Links:       post.Links,
			DateCreated: timestamppb.New(post.DateCreated),
			UserId:      post.UserId,
			IsJobOffer:  post.IsJobOffer,
			Company: &jobOffer.Company{
				Id:          post.Company.Id,
				Name:        html.UnescapeString(post.Company.Name),
				Description: html.UnescapeString(post.Company.Description),
				PhoneNumber: post.Company.PhoneNumber,
				IsActive:    true,
			},
			JobOffer: &jobOffer.JobOffer{
				Id: post.JobOffer.Id,
				Position: &jobOffer.Position{
					Id:   post.JobOffer.Position.Id,
					Name: post.JobOffer.Position.Name,
					Pay:  post.JobOffer.Position.Pay,
				},
				Preconditions:   html.UnescapeString(post.JobOffer.Preconditions),
				DailyActivities: html.UnescapeString(post.JobOffer.DailyActivities),
				JobDescription:  html.UnescapeString(post.JobOffer.JobDescription),
			},
		}
		response.JobOffers = append(response.JobOffers, newPost)
	}

	req := &jobOffer.GetRequest{
		DTO: &jobOffer.Recommendation{
			User:      user,
			JobOffers: response,
		},
	}
	allJobs, err := jobOfferClient.GetRecommendations(context.TODO(), req)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error while finding job recommendations")

	}

	fmt.Println(allJobs)
	return allJobs
}
