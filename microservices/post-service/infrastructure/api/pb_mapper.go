package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"html"
	"strconv"
	"strings"
	"time"
)

func mapPost(post *domain.Post) *pb.Post {
	if post.IsJobOffer {
		postPb := &pb.Post{
			Id:          post.Id.Hex(),
			Text:        html.UnescapeString(post.Text), /** UnescapeString **/
			DateCreated: timestamppb.New(post.DateCreated),
			UserId:      post.UserId,
			Images:      post.Images,
			Links:       post.Links,
			IsJobOffer:  post.IsJobOffer,
			Company: &pb.Company{
				//Id:          float64(companyIdInt),
				Name:        html.UnescapeString(post.Company.Name),
				Description: html.UnescapeString(post.Company.Description),
				PhoneNumber: post.Company.PhoneNumber,
				IsActive:    true,
			},
			JobOffer: &pb.JobOffer{
				//Id: post.JobOffer.Id.Hex().,
				Position: &pb.Position{
					//Id: post.JobOffer.Id.Hex().,
					Name: post.JobOffer.Position.Name,
					Pay:  post.JobOffer.Position.Pay,
				},
				Preconditions:   html.UnescapeString(post.JobOffer.Preconditions),
				DailyActivities: html.UnescapeString(post.JobOffer.DailyActivities),
				JobDescription:  html.UnescapeString(post.JobOffer.JobDescription),
			},
		}

		for _, like := range post.Likes {
			postPb.Likes = append(postPb.Likes, &pb.Like{
				Id:     like.Id.Hex(),
				UserId: like.UserId,
			})
		}

		for _, dislike := range post.Dislikes {
			postPb.Dislikes = append(postPb.Dislikes, &pb.Dislike{
				Id:     dislike.Id.Hex(),
				UserId: dislike.UserId,
			})
		}

		for _, comment := range post.Comments {
			postPb.Comments = append(postPb.Comments, &pb.Comment{
				Id:     comment.Id.Hex(),
				UserId: comment.UserId,
				Text:   html.UnescapeString(comment.Text), /** UnescapeString **/
			})
		}
		return postPb
	} else {
		postPb := &pb.Post{
			Id:          post.Id.Hex(),
			Text:        html.UnescapeString(post.Text), /** UnescapeString **/
			DateCreated: timestamppb.New(post.DateCreated),
			UserId:      post.UserId,
			Images:      post.Images,
			Links:       post.Links,
			IsJobOffer:  post.IsJobOffer,
		}

		for _, like := range post.Likes {
			postPb.Likes = append(postPb.Likes, &pb.Like{
				Id:     like.Id.Hex(),
				UserId: like.UserId,
			})
		}

		for _, dislike := range post.Dislikes {
			postPb.Dislikes = append(postPb.Dislikes, &pb.Dislike{
				Id:     dislike.Id.Hex(),
				UserId: dislike.UserId,
			})
		}

		for _, comment := range post.Comments {
			postPb.Comments = append(postPb.Comments, &pb.Comment{
				Id:     comment.Id.Hex(),
				UserId: comment.UserId,
				Text:   html.UnescapeString(comment.Text), /** UnescapeString **/
			})
		}
		return postPb
	}
}

func mapInsertPost(post *pb.InsertPost) (*domain.Post, error) {
	postPb := &domain.Post{
		Text:        strings.TrimSpace(post.Text), // function to remove leading and trailing whitespace
		Images:      post.Images,
		Links:       post.Links,
		DateCreated: time.Now(),
		IsJobOffer:  false,
	}

	return postPb, nil
}

func mapInsertJobOfferPost(post *pb.InsertJobOfferPost) (*domain.Post, error) {
	postPb := &domain.Post{
		Text:        strings.TrimSpace(post.Text),
		DateCreated: time.Now(),
		IsJobOffer:  true,
		JobOffer: domain.JobOffer{
			Id:              primitive.NewObjectID(),
			JobDescription:  post.JobOffer.JobDescription,
			DailyActivities: post.JobOffer.DailyActivities,
			Preconditions:   post.JobOffer.Preconditions,
			Position: domain.Position{
				Id:   primitive.NewObjectID(),
				Name: post.JobOffer.Position.Name,
				Pay:  post.JobOffer.Position.Pay,
			},
		},
		Company: domain.Company{
			Id:          primitive.NewObjectID(),
			Name:        post.Company.Name,
			Description: post.Company.Description,
			PhoneNumber: post.Company.PhoneNumber,
			IsActive:    true,
		},
	}

	return postPb, nil
}

func mapInsertPosition(position *pb.Position) (*domain.Position, error) {
	positionPb := &domain.Position{
		Id:   primitive.NewObjectID(),
		Name: position.Name,
		Pay:  position.Pay,
	}

	return positionPb, nil
}

func mapInsertCompany(company *pb.Company) (*domain.Company, error) {
	copmanyPb := &domain.Company{
		Id:          primitive.NewObjectID(),
		Name:        company.Name,
		Description: company.Description,
		PhoneNumber: company.PhoneNumber,
		IsActive:    company.IsActive,
	}

	return copmanyPb, nil
}

func mapInsertJobOffer(jobOffer *pb.JobOffer) (*domain.JobOffer, error) {
	var jobOfferPb = &domain.JobOffer{
		Id:              primitive.NewObjectID(),
		JobDescription:  jobOffer.JobDescription,
		DailyActivities: jobOffer.DailyActivities,
		Preconditions:   jobOffer.Preconditions,
	}

	return jobOfferPb, nil
}

func mapJobOffer(jobOffer *pb.JobOffer) (*domain.JobOffer, error) {
	v := strconv.FormatFloat(jobOffer.Id, 64, 5, 5)
	v1, _ := primitive.ObjectIDFromHex(v)

	var jobOfferPb = &domain.JobOffer{
		Id:              v1, // TODO !
		JobDescription:  jobOffer.JobDescription,
		DailyActivities: jobOffer.DailyActivities,
		Preconditions:   jobOffer.Preconditions,
		Position: domain.Position{
			Name: jobOffer.Position.Name,
			Pay:  jobOffer.Position.Pay,
		},
	}

	return jobOfferPb, nil
}

func mapCompany(company *pb.Company) (*domain.Company, error) {
	copmanyPb := &domain.Company{
		//Id:          company.Id,		// TODO !
		Name:        company.Name,
		Description: company.Description,
		PhoneNumber: company.PhoneNumber,
		IsActive:    company.IsActive,
	}

	return copmanyPb, nil
}

func mapPosition(position *pb.Position) (*domain.Position, error) {
	positionPb := &domain.Position{
		//Id:   primitive.NewObjectID(),	// TODO ?
		Name: position.Name,
		Pay:  position.Pay,
	}

	return positionPb, nil
}
