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
	postPb := &pb.Post{
		Id:          post.Id.Hex(),
		Text:        html.UnescapeString(post.Text), /** UnescapeString **/
		DateCreated: timestamppb.New(post.DateCreated),
		UserId:      post.UserId,
		Images:      post.Images,
		Links:       post.Links,
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

func mapInsertPost(post *pb.InsertPost) (*domain.Post, error) {
	postPb := &domain.Post{
		Text:        strings.TrimSpace(post.Text), // function to remove leading and trailing whitespace
		Images:      post.Images,
		Links:       post.Links,
		DateCreated: time.Now(),
		//IsJobOffer: false,		// TODO ??
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

func mapInsertJobOfferPost(post *pb.InsertJobOfferPost) (*domain.Post, error) {
	postPb := &domain.Post{
		Text:        strings.TrimSpace(post.Text), // function to remove leading and trailing whitespace
		DateCreated: time.Now(),
		IsJobOffer:  true,
	}

	return postPb, nil
}

//func mapUpdatePost(oldData *pb.InsertPost, newData *pb.InsertPost) *domain.Post {
//	id, _ := primitive.ObjectIDFromHex(oldData.Id)
//
//	postPb := &domain.Post{
//		Id:     id,
//		Text:   newData.Text,
//		UserId: oldData.UserId, // ne moze se kreator post-a
//		Images: newData.Images,
//		Links:  newData.Links,
//	}
//
//	for _, like := range newData.Likes {
//		if like.Id == "" {
//			if likedPostByUser(oldData, like.UserId) == false {
//				like_id := primitive.NewObjectID()
//				postPb.Likes = append(postPb.Likes, domain.Like{
//					Id:     like_id,
//					UserId: like.UserId,
//				})
//			}
//		} else {
//			like_id, _ := primitive.ObjectIDFromHex(like.Id)
//			postPb.Likes = append(postPb.Likes, domain.Like{
//				Id:     like_id,
//				UserId: like.UserId,
//			})
//		}
//	}
//
//	for _, dislike := range newData.Dislikes {
//		if dislike.Id == "" {
//			if dislikedPostByUser(oldData, dislike.UserId) == false {
//				dislike_id := primitive.NewObjectID()
//				postPb.Dislikes = append(postPb.Dislikes, domain.Dislike{
//					Id:     dislike_id,
//					UserId: dislike.UserId,
//				})
//			}
//		} else {
//			dislike_id, _ := primitive.ObjectIDFromHex(dislike.Id)
//			postPb.Dislikes = append(postPb.Dislikes, domain.Dislike{
//				Id:     dislike_id,
//				UserId: dislike.UserId,
//			})
//		}
//	}
//
//	for _, comment := range newData.Comments {
//		if comment.Id == "" {
//			comment_id := primitive.NewObjectID()
//			postPb.Comments = append(postPb.Comments, domain.Comment{
//				Id:     comment_id,
//				UserId: comment.UserId,
//				Text:   comment.Text,
//			})
//		} else {
//			comment_id, _ := primitive.ObjectIDFromHex(comment.Id)
//			postPb.Comments = append(postPb.Comments, domain.Comment{
//				Id:     comment_id,
//				UserId: comment.UserId,
//				Text:   comment.Text,
//			})
//		}
//	}
//
//	return postPb
//}

//func likedPostByUser(post *pb.InsertPost, userId string) bool {
//	for _, like := range post.Likes {
//		if like.UserId == userId {
//			fmt.Println("Postoji duplikat - like")
//			fmt.Println("ISPIS:", like.UserId, userId)
//			return true
//		}
//	}
//	return false
//}
//
//func dislikedPostByUser(post *pb.InsertPost, userId string) bool {
//	for _, dislike := range post.Dislikes {
//		if dislike.UserId == userId {
//			fmt.Println("Postoji duplikat - dislike")
//			fmt.Println("ISPIS:", dislike.UserId, userId)
//			return true
//		}
//	}
//	return false
//}
