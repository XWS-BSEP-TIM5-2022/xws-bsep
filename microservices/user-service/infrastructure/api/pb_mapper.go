package api

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
)

func mapUser(user *domain.User) *pb.User {
	userPb := &pb.User{
		Id:           user.Id.Hex(),
		Username:     user.Username,
		Name:         user.Name,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Gender:       mapGender(user.Gender),
		Birthday:     timestamppb.New(user.Birthday),
		Email:        user.Email,
		Biography:    user.Biography,
		Password:     user.Password,
	}

	for _, education := range user.Education {
		userPb.Education = append(userPb.Education, &pb.Education{
			Id:        education.Id.Hex(),
			Name:      education.Name,
			Level:     mapEducation(education.Level),
			Place:     education.Place,
			StartDate: timestamppb.New(education.StartDate),
			EndDate:   timestamppb.New(education.EndDate),
		})
	}

	for _, experience := range user.Experience {
		userPb.Experience = append(userPb.Experience, &pb.Experience{
			Id:        experience.Id.Hex(),
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: timestamppb.New(experience.StartDate),
			EndDate:   timestamppb.New(experience.EndDate),
		})
	}

	for _, skill := range user.Skills {
		userPb.Skills = append(userPb.Skills, &pb.Skill{
			Id:   skill.Id.Hex(),
			Name: skill.Name,
		})
	}

	for _, interest := range user.Interests {
		userPb.Interests = append(userPb.Interests, &pb.Interest{
			Id:          interest.Id.Hex(),
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	return userPb
}

func mapInsertUser(user *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(user.Id)

	userPb := &domain.User{
		Id:           id,
		Username:     user.Username,
		Name:         user.Name,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Gender:       mapInsertGender(user.Gender),
		Email:        user.Email,
		Biography:    user.Biography,
		Password:     user.Password,
	}

	if user.Birthday != nil {
		userPb.Birthday = user.Birthday.AsTime()
	}

	for _, education := range user.Education {

		ed_id, _ := primitive.ObjectIDFromHex(user.Id)

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     mapInsertEducation(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	for _, experience := range user.Experience {

		ex_id, _ := primitive.ObjectIDFromHex(user.Id)

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	for _, skill := range user.Skills {

		s_id, _ := primitive.ObjectIDFromHex(user.Id)

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	for _, interest := range user.Interests {

		in_id, _ := primitive.ObjectIDFromHex(user.Id)

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	return userPb
}

func mapInsertGender(gender pb.User_GenderEnum) domain.GenderEnum {
	switch gender {
	case pb.User_Male:
		return domain.Male
	}
	return domain.Female

}

func mapInsertEducation(education pb.Education_EducationEnum) domain.EducationEnum {

	switch education {
	case pb.Education_Primary:
		return domain.Primary
	case pb.Education_Secondary:
		return domain.Secondary
	case pb.Education_Bachelor:
		return domain.Bachelor
	case pb.Education_Master:
		return domain.Master
	}
	return domain.Doctorate

}

func mapGender(gender domain.GenderEnum) pb.User_GenderEnum {
	switch gender {
	case domain.Male:
		return pb.User_Male
	}
	return pb.User_Female

}

func mapEducation(education domain.EducationEnum) pb.Education_EducationEnum {
	switch education {
	case domain.Primary:
		return pb.Education_Primary
	case domain.Secondary:
		return pb.Education_Secondary
	case domain.Bachelor:
		return pb.Education_Bachelor
	case domain.Master:
		return pb.Education_Master
	}
	return pb.Education_Doctorate

}
