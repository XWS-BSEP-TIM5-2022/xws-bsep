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
		Name:         user.Name,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Gender:       mapGender(user.Gender),
		Birthday:     timestamppb.New(user.Birthday),
		Email:        user.Email,
		Biography:    user.Biography,
		IsPublic:     user.IsPublic,
		IsActive:     user.IsActive,
		Role:         user.Role,
		Username:     user.Username,
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
		Name:         user.Name,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Gender:       mapInsertGender(user.Gender),
		Email:        user.Email,
		Biography:    user.Biography,
		IsPublic:     user.IsPublic,
		IsActive:     user.IsActive,
		Role:         user.Role,
		Username:     user.Username,
	}

	if user.Birthday != nil {
		userPb.Birthday = user.Birthday.AsTime()
	}

	for _, education := range user.Education {

		ed_id := primitive.NewObjectID()

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

		ex_id := primitive.NewObjectID()

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

		s_id := primitive.NewObjectID()

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	for _, interest := range user.Interests {

		in_id := primitive.NewObjectID()

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	return userPb
}

func mapUpdateUser(oldData *pb.User, newData *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	userPb := &domain.User{
		Id:           id,
		Name:         newData.Name,
		LastName:     newData.LastName,
		MobileNumber: newData.MobileNumber,
		Gender:       mapInsertGender(newData.Gender),
		Email:        newData.Email,
		Biography:    newData.Biography,
		IsPublic:     oldData.IsPublic,
		IsActive:     oldData.IsActive,
		Role:         oldData.Role,
		Username:     oldData.Username,
	}

	if mapInsertGender(newData.Gender) == -1 {
		userPb.Gender = mapInsertGender(oldData.Gender)
	}

	if newData.Birthday != nil {
		userPb.Birthday = newData.Birthday.AsTime()
	}

	if newData.Username == "" {
		userPb.Username = oldData.Username
	}

	if newData.Name == "" {
		userPb.Name = oldData.Name
	}

	if newData.LastName == "" {
		userPb.LastName = oldData.LastName
	}

	educations := newData.Education

	for _, education := range educations {

		ed_id := primitive.NewObjectID()

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     mapInsertEducation(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	experiences := newData.Experience

	for _, experience := range experiences {

		ex_id := primitive.NewObjectID()

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	skills := newData.Skills

	for _, skill := range skills {

		s_id := primitive.NewObjectID()

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	interests := newData.Interests

	for _, interest := range interests {

		in_id := primitive.NewObjectID()

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
	case pb.User_Female:
		return domain.Female
	case pb.User_Male:
		return domain.Male
	}
	return -1

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
	case pb.Education_Doctorate:
		return domain.Doctorate
	}
	return -1

}

func mapGender(gender domain.GenderEnum) pb.User_GenderEnum {
	switch gender {
	case domain.Male:
		return pb.User_Male
	case domain.Female:
		return pb.User_Female
	}
	return -1

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
	case domain.Doctorate:
		return pb.Education_Doctorate
	}
	return -1

}

func mapBasicInfo(oldData *pb.User, newData *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	userPb := &domain.User{
		Id:           id,
		Name:         newData.Name,
		LastName:     newData.LastName,
		MobileNumber: newData.MobileNumber,
		Gender:       mapInsertGender(newData.Gender),
		Email:        newData.Email,
		Biography:    newData.Biography,
		IsPublic:     oldData.IsPublic,
		IsActive:     oldData.IsActive,
		Role:         oldData.Role,
		Username:     oldData.Username,
	}

	if mapInsertGender(newData.Gender) == -1 {
		userPb.Gender = mapInsertGender(oldData.Gender)
	}

	if newData.Birthday != nil {
		userPb.Birthday = newData.Birthday.AsTime()
	}

	if newData.Name == "" {
		userPb.Name = oldData.Name
	}

	if newData.LastName == "" {
		userPb.LastName = oldData.LastName
	}

	educations := oldData.Education

	for _, education := range educations {

		ed_id, _ := primitive.ObjectIDFromHex(education.Id)

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     mapInsertEducation(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	experiences := oldData.Experience

	for _, experience := range experiences {

		ex_id, _ := primitive.ObjectIDFromHex(experience.Id)

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	skills := oldData.Skills

	for _, skill := range skills {

		s_id, _ := primitive.ObjectIDFromHex(skill.Id)

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	interests := oldData.Interests

	for _, interest := range interests {

		in_id, _ := primitive.ObjectIDFromHex(interest.Id)

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	return userPb
}

func mapExperienceAndEducation(oldData *pb.User, newData *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	userPb := &domain.User{
		Id:           id,
		Name:         oldData.Name,
		LastName:     oldData.LastName,
		MobileNumber: oldData.MobileNumber,
		Gender:       mapInsertGender(oldData.Gender),
		Birthday:     oldData.Birthday.AsTime(),
		Email:        oldData.Email,
		Biography:    oldData.Biography,
		IsPublic:     oldData.IsPublic,
		IsActive:     oldData.IsActive,
		Role:         oldData.Role,
		Username:     oldData.Username,
	}

	educations := newData.Education

	for _, education := range educations {

		ed_id := primitive.NewObjectID()

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     mapInsertEducation(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	experiences := newData.Experience

	for _, experience := range experiences {

		ex_id := primitive.NewObjectID()

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	skills := oldData.Skills

	for _, skill := range skills {

		s_id, _ := primitive.ObjectIDFromHex(skill.Id)

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	interests := oldData.Interests

	for _, interest := range interests {

		in_id, _ := primitive.ObjectIDFromHex(interest.Id)

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	return userPb
}

func mapSkillsAndInterests(oldData *pb.User, newData *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	userPb := &domain.User{
		Id:           id,
		Name:         oldData.Name,
		LastName:     oldData.LastName,
		MobileNumber: oldData.MobileNumber,
		Gender:       mapInsertGender(oldData.Gender),
		Birthday:     oldData.Birthday.AsTime(),
		Email:        oldData.Email,
		Biography:    oldData.Biography,
		IsPublic:     oldData.IsPublic,
		IsActive:     oldData.IsActive,
		Role:         oldData.Role,
		Username:     oldData.Username,
	}

	educations := oldData.Education

	for _, education := range educations {

		ed_id, _ := primitive.ObjectIDFromHex(education.Id)

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     mapInsertEducation(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	experiences := oldData.Experience

	for _, experience := range experiences {

		ex_id, _ := primitive.ObjectIDFromHex(experience.Id)

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	skills := newData.Skills

	for _, skill := range skills {

		s_id := primitive.NewObjectID()

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	interests := newData.Interests

	for _, interest := range interests {

		in_id := primitive.NewObjectID()

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	return userPb
}
