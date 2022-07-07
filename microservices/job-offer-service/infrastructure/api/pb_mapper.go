package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func mapUser(user *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(user.Id)

	userPb := &domain.User{
		Id:       id,
		Name:     removeMalicious(user.Name),
		LastName: removeMalicious(user.LastName),
		Email:    user.Email,
		IsPublic: user.IsPublic,
	}

	for _, experience := range user.Experience {

		id, _ := primitive.ObjectIDFromHex(experience.Id)

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:       id,
			Headline: removeMalicious(experience.Headline),
		})
	}

	for _, skill := range user.Skills {

		id, _ := primitive.ObjectIDFromHex(skill.Id)

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   id,
			Name: removeMalicious(skill.Name),
		})
	}
	return userPb
}

func removeMalicious(value string) string {

	var lenId = len(value)
	var checkId = ""
	for i := 0; i < lenId; i++ {
		char := string(value[i])
		if char != "$" {
			checkId = checkId + char
		}
	}
	return checkId
}

func mapJobOffer(post *pb.Post) *domain.Post {
	id, _ := primitive.ObjectIDFromHex(post.Id)

	postPb := &domain.Post{
		Id:         id,
		Text:       strings.TrimSpace(post.Text),
		IsJobOffer: post.IsJobOffer,
		JobOffer: domain.JobOffer{
			Id:            primitive.NewObjectID(),
			Preconditions: strings.TrimSpace(post.JobOffer.Preconditions),
			Position: domain.Position{
				Id:   primitive.NewObjectID(),
				Name: strings.TrimSpace(post.JobOffer.Position.Name),
				Pay:  post.JobOffer.Position.Pay,
			},
		},
		Company: domain.Company{
			Id:          primitive.NewObjectID(), //TODO:ispravi
			Name:        strings.TrimSpace(post.Company.Name),
			Description: strings.TrimSpace(post.Company.Description),
			PhoneNumber: post.Company.PhoneNumber,
			IsActive:    true,
		},
	}

	return postPb
}
