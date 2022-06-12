package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"html"
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
			Image:       post.Image,
			Links:       post.Links,
			IsJobOffer:  post.IsJobOffer,
			Company: &pb.Company{
				Id:          post.Company.Id.Hex(),
				Name:        html.UnescapeString(post.Company.Name),
				Description: html.UnescapeString(post.Company.Description),
				PhoneNumber: post.Company.PhoneNumber,
				IsActive:    true,
			},
			JobOffer: &pb.JobOffer{
				Id: post.JobOffer.Id.Hex(),
				Position: &pb.Position{
					Id:   post.JobOffer.Position.Id.Hex(),
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
			Image:       post.Image,
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
		Image:       post.Image,
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
			JobDescription:  strings.TrimSpace(post.JobOffer.JobDescription),
			DailyActivities: strings.TrimSpace(post.JobOffer.DailyActivities),
			Preconditions:   strings.TrimSpace(post.JobOffer.Preconditions),
			Position: domain.Position{
				Id:   primitive.NewObjectID(),
				Name: strings.TrimSpace(post.JobOffer.Position.Name),
				Pay:  post.JobOffer.Position.Pay,
			},
		},
		Company: domain.Company{
			Id:          primitive.NewObjectID(),
			Name:        strings.TrimSpace(post.Company.Name),
			Description: strings.TrimSpace(post.Company.Description),
			PhoneNumber: post.Company.PhoneNumber,
			IsActive:    true,
		},
	}

	return postPb, nil
}

func mapCompanyInfo(company *pb.CompanyInfoDTO) (*domain.Company, error) {
	companyPb := &domain.Company{
		Name:        company.Name,
		Description: company.Description,
		PhoneNumber: company.PhoneNumber,
		IsActive:    company.IsActive,
	}

	return companyPb, nil
}

//func encodeImage(image primitive.Binary) string {
//	return base64.StdEncoding.EncodeToString(image.Data)
//}
//
//func decodeImage(path string) (primitive.Binary, error) {
//
//	fmt.Println("usao sam")
//	image, err := base64.StdEncoding.DecodeString(path)
//	if err != nil {
//		return primitive.Binary{}, err
//	}
//	return primitive.Binary{Data: image}, nil
// }
