package api

//
//import (
//	"google.golang.org/protobuf/types/known/timestamppb"
//
//	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
//	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
//)
//
//func mapUser(user *domain.User) *pb.User {
//	userPb := &pb.User{
//		Id:           user.Id,
//		Username:     user.Username,
//		Name:         user.Name,
//		LastName:     user.LastName,
//		MobileNumber: user.MobileNumber,
//		Gender:       mapGender(user.Gender),
//		Birthday:     timestamppb.New(user.Birthday),
//		Email:        user.Email,
//		Biography:    user.Biography,
//		Password:     user.Password,
//	}
//
//	return userPb
//}
//
//func mapInsertUser(user *pb.User) *domain.User {
//	userPb := &domain.User{
//		Id:           user.Id,
//		Username:     user.Username,
//		Name:         user.Name,
//		LastName:     user.LastName,
//		MobileNumber: user.MobileNumber,
//		Gender:       mapInsertGender(user.Gender),
//		Birthday:     user.Birthday.AsTime(),
//		Email:        user.Email,
//		Biography:    user.Biography,
//		Password:     user.Password,
//	}
//
//	return userPb
//}
//
//func mapInsertGender(status pb.User_GenderEnum) domain.GenderEnum {
//	switch status {
//	case pb.User_Male:
//		return domain.Male
//	}
//	return domain.Female
//
//}
//
//func mapGender(status domain.GenderEnum) pb.User_GenderEnum {
//	switch status {
//	case domain.Male:
//		return pb.User_Male
//	}
//	return pb.User_Female
//
//}
