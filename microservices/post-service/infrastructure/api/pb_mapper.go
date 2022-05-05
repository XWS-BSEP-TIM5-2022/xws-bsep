package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			UserId: post.UserId,
		})
	}

	for _, dislike := range post.Dislikes {
		postPb.Dislikes = append(postPb.Dislikes, &pb.Dislike{
			Id:     dislike.Id.Hex(),
			UserId: post.UserId,
		})
	}

	for _, comment := range post.Comments {
		postPb.Comments = append(postPb.Comments, &pb.Comment{
			Id:     comment.Id.Hex(),
			UserId: post.UserId,
			Text:   post.Text,
		})
	}

	return postPb
}

func mapInsertPost(post *pb.Post) *domain.Post {
	id, _ := primitive.ObjectIDFromHex(post.Id)

	postPb := &domain.Post{
		Id:     id,
		Text:   post.Text,
		UserId: post.UserId,
		Images: post.Images,
		Links:  post.Links,
	}

	if post.DateCreated != nil {
		postPb.DateCreated = post.DateCreated.AsTime()
	}

	for _, like := range post.Likes {
		like_id, _ := primitive.ObjectIDFromHex(post.Id)
		postPb.Likes = append(postPb.Likes, domain.Like{
			Id:     like_id,
			UserId: like.UserId,
		})
	}

	for _, dislike := range post.Dislikes {
		dislike_id, _ := primitive.ObjectIDFromHex(post.Id)
		postPb.Dislikes = append(postPb.Dislikes, domain.Dislike{
			Id:     dislike_id,
			UserId: dislike.UserId,
		})
	}

	for _, comment := range post.Comments {
		comment_id, _ := primitive.ObjectIDFromHex(post.Id)
		postPb.Comments = append(postPb.Comments, domain.Comment{
			Id:     comment_id,
			UserId: comment.UserId,
			Text:   comment.Text,
		})
	}

	return postPb
}
