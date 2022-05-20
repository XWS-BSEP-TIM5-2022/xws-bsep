package api

import (
	"fmt"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func mapPost(post *domain.Post) *pb.Post {
	postPb := &pb.Post{
		Id:          post.Id.Hex(),
		Text:        post.Text,
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
			Text:   comment.Text,
		})
	}

	return postPb
}

func mapInsertPost(post *pb.Post) (*domain.Post, error) {
	id, _ := primitive.ObjectIDFromHex(post.Id)

	postPb := &domain.Post{
		Id:          id,
		Text:        post.Text,
		UserId:      post.UserId,
		Images:      post.Images,
		Links:       post.Links,
		DateCreated: time.Now(),
	}

	err := validate.Struct(postPb)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("---------------- pocetak greske ----------------")
			fmt.Println(err.Field())
			fmt.Println(err.Tag())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------- kraj greske ----------------")
		}

		return nil, err
	}

	return postPb, nil
}

func mapLike(like *pb.Like) *domain.Like { // TODO : InsertLike ???
	id := primitive.NewObjectID()

	likePb := &domain.Like{
		Id:     id,
		UserId: like.UserId,
	}

	return likePb
}

func mapUpdatePost(oldData *pb.Post, newData *pb.Post) *domain.Post {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	postPb := &domain.Post{
		Id:          id,
		Text:        newData.Text,
		UserId:      oldData.UserId, // ne moze se kreator post-a
		Images:      newData.Images,
		Links:       newData.Links,
		DateCreated: oldData.DateCreated.AsTime(),
	}

	for _, like := range newData.Likes {
		if like.Id == "" {
			if likedPostByUser(oldData, like.UserId) == false {
				like_id := primitive.NewObjectID()
				postPb.Likes = append(postPb.Likes, domain.Like{
					Id:     like_id,
					UserId: like.UserId,
				})
			}
		} else {
			like_id, _ := primitive.ObjectIDFromHex(like.Id)
			postPb.Likes = append(postPb.Likes, domain.Like{
				Id:     like_id,
				UserId: like.UserId,
			})
		}
	}

	for _, dislike := range newData.Dislikes {
		if dislike.Id == "" {
			if dislikedPostByUser(oldData, dislike.UserId) == false {
				dislike_id := primitive.NewObjectID()
				postPb.Dislikes = append(postPb.Dislikes, domain.Dislike{
					Id:     dislike_id,
					UserId: dislike.UserId,
				})
			}
		} else {
			dislike_id, _ := primitive.ObjectIDFromHex(dislike.Id)
			postPb.Dislikes = append(postPb.Dislikes, domain.Dislike{
				Id:     dislike_id,
				UserId: dislike.UserId,
			})
		}
	}

	for _, comment := range newData.Comments {
		if comment.Id == "" {
			comment_id := primitive.NewObjectID()
			postPb.Comments = append(postPb.Comments, domain.Comment{
				Id:     comment_id,
				UserId: comment.UserId,
				Text:   comment.Text,
			})
		} else {
			comment_id, _ := primitive.ObjectIDFromHex(comment.Id)
			postPb.Comments = append(postPb.Comments, domain.Comment{
				Id:     comment_id,
				UserId: comment.UserId,
				Text:   comment.Text,
			})
		}
	}

	return postPb
}

func likedPostByUser(post *pb.Post, userId string) bool {
	for _, like := range post.Likes {
		if like.UserId == userId {
			fmt.Println("Postoji duplikat - like")
			fmt.Println("ISPIS:", like.UserId, userId)
			return true
		}
	}
	return false
}

func dislikedPostByUser(post *pb.Post, userId string) bool {
	for _, dislike := range post.Dislikes {
		if dislike.UserId == userId {
			fmt.Println("Postoji duplikat - dislike")
			fmt.Println("ISPIS:", dislike.UserId, userId)
			return true
		}
	}
	return false
}
