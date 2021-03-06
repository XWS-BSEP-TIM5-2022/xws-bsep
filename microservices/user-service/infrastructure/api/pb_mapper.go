package api

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	events "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/create_user"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
)

func mapUser(user *domain.User) *pb.User {
	userPb := &pb.User{
		Id:               user.Id.Hex(),
		Name:             user.Name,
		LastName:         user.LastName,
		MobileNumber:     user.MobileNumber,
		Gender:           mapGender(user.Gender),
		Birthday:         timestamppb.New(user.Birthday),
		Email:            user.Email,
		Biography:        user.Biography,
		IsPublic:         user.IsPublic,
		IsActive:         user.IsActive,
		Role:             user.Role,
		Username:         user.Username,
		PostNotification: user.PostNotification,
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
		Id:               id,
		Name:             removeMalicious(user.Name),
		LastName:         removeMalicious(user.LastName),
		MobileNumber:     removeMalicious(user.MobileNumber),
		Gender:           mapInsertGender(user.Gender),
		Email:            user.Email,
		Biography:        removeMalicious(user.Biography),
		IsPublic:         user.IsPublic,
		IsActive:         user.IsActive,
		Role:             user.Role,
		Username:         user.Username,
		PostNotification: true,
	}

	if user.Birthday != nil {
		userPb.Birthday = user.Birthday.AsTime()
	}

	for _, education := range user.Education {

		ed_id := primitive.NewObjectID()

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      removeMalicious(education.Name),
			Level:     mapInsertEducation(education.Level),
			Place:     removeMalicious(education.Place),
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	for _, experience := range user.Experience {

		ex_id := primitive.NewObjectID()

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      removeMalicious(experience.Name),
			Headline:  removeMalicious(experience.Headline),
			Place:     removeMalicious(experience.Place),
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	for _, skill := range user.Skills {

		s_id := primitive.NewObjectID()

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: removeMalicious(skill.Name),
		})
	}

	for _, interest := range user.Interests {

		in_id := primitive.NewObjectID()

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        removeMalicious(interest.Name),
			Description: removeMalicious(interest.Description),
		})
	}

	return userPb
}

func mapUpdateUser(oldData *pb.User, newData *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	userPb := &domain.User{
		Id:               id,
		Name:             removeMalicious(newData.Name),
		LastName:         removeMalicious(newData.LastName),
		MobileNumber:     removeMalicious(newData.MobileNumber),
		Gender:           mapInsertGender(newData.Gender),
		Email:            newData.Email,
		Biography:        removeMalicious(newData.Biography),
		IsPublic:         oldData.IsPublic,
		IsActive:         oldData.IsActive,
		Role:             oldData.Role,
		Username:         newData.Username,
		PostNotification: oldData.PostNotification,
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
			Name:      removeMalicious(education.Name),
			Level:     mapInsertEducation(education.Level),
			Place:     removeMalicious(education.Place),
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	experiences := newData.Experience

	for _, experience := range experiences {

		ex_id := primitive.NewObjectID()

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      removeMalicious(experience.Name),
			Headline:  removeMalicious(experience.Headline),
			Place:     removeMalicious(experience.Place),
			StartDate: experience.StartDate.AsTime(),
			EndDate:   experience.EndDate.AsTime(),
		})
	}

	skills := newData.Skills

	for _, skill := range skills {

		s_id := primitive.NewObjectID()

		userPb.Skills = append(userPb.Skills, domain.Skill{
			Id:   s_id,
			Name: removeMalicious(skill.Name),
		})
	}

	interests := newData.Interests

	for _, interest := range interests {

		in_id := primitive.NewObjectID()

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        removeMalicious(interest.Name),
			Description: removeMalicious(interest.Description),
		})
	}

	return userPb
}

func mapUpdateNotificationUser(oldData *pb.User, newData *pb.User) *domain.User {
	id, _ := primitive.ObjectIDFromHex(oldData.Id)

	userPb := &domain.User{
		Id:               id,
		Name:             removeMalicious(oldData.Name),
		LastName:         removeMalicious(oldData.LastName),
		MobileNumber:     removeMalicious(oldData.MobileNumber),
		Gender:           mapInsertGender(oldData.Gender),
		Email:            oldData.Email,
		Biography:        removeMalicious(oldData.Biography),
		IsPublic:         oldData.IsPublic,
		IsActive:         oldData.IsActive,
		Role:             oldData.Role,
		Username:         oldData.Username,
		PostNotification: newData.PostNotification,
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
		Id:               id,
		Name:             removeMalicious(newData.Name),
		LastName:         removeMalicious(newData.LastName),
		MobileNumber:     removeMalicious(newData.MobileNumber),
		Gender:           mapInsertGender(newData.Gender),
		Email:            newData.Email,
		Biography:        removeMalicious(newData.Biography),
		IsPublic:         oldData.IsPublic,
		IsActive:         oldData.IsActive,
		Role:             oldData.Role,
		Username:         newData.Username,
		PostNotification: oldData.PostNotification,
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
		Id:               id,
		Name:             oldData.Name,
		LastName:         oldData.LastName,
		MobileNumber:     oldData.MobileNumber,
		Gender:           mapInsertGender(oldData.Gender),
		Birthday:         oldData.Birthday.AsTime(),
		Email:            oldData.Email,
		Biography:        oldData.Biography,
		IsPublic:         oldData.IsPublic,
		IsActive:         oldData.IsActive,
		Role:             oldData.Role,
		Username:         oldData.Username,
		PostNotification: oldData.PostNotification,
	}

	educations := newData.Education

	for _, education := range educations {

		ed_id := primitive.NewObjectID()

		userPb.Education = append(userPb.Education, domain.Education{
			Id:        ed_id,
			Name:      removeMalicious(education.Name),
			Level:     mapInsertEducation(education.Level),
			Place:     removeMalicious(education.Place),
			StartDate: education.StartDate.AsTime(),
			EndDate:   education.EndDate.AsTime(),
		})
	}

	experiences := newData.Experience

	for _, experience := range experiences {

		ex_id := primitive.NewObjectID()

		userPb.Experience = append(userPb.Experience, domain.Experience{
			Id:        ex_id,
			Name:      removeMalicious(experience.Name),
			Headline:  removeMalicious(experience.Headline),
			Place:     removeMalicious(experience.Place),
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
		Id:               id,
		Name:             oldData.Name,
		LastName:         oldData.LastName,
		MobileNumber:     oldData.MobileNumber,
		Gender:           mapInsertGender(oldData.Gender),
		Birthday:         oldData.Birthday.AsTime(),
		Email:            oldData.Email,
		Biography:        oldData.Biography,
		IsPublic:         oldData.IsPublic,
		IsActive:         oldData.IsActive,
		Role:             oldData.Role,
		Username:         oldData.Username,
		PostNotification: oldData.PostNotification,
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
			Name: removeMalicious(skill.Name),
		})
	}

	interests := newData.Interests

	for _, interest := range interests {

		in_id := primitive.NewObjectID()

		userPb.Interests = append(userPb.Interests, domain.Interest{
			Id:          in_id,
			Name:        removeMalicious(interest.Name),
			Description: interest.Description,
		})
	}

	return userPb
}

// SD - saga ----------------------------------------------------
func mapInsertUserSaga(user *pb.RegisterRequest) *domain.User {
	userPb := &domain.User{
		// Id:           id,
		Name:         user.Name,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Gender:       domain.GenderEnum(user.Gender),
		Email:        user.Email,
		Biography:    user.Biography,
		IsActive:     false,
		Role:         user.Role,
		// IsPublic:     user.IsPublic,
		PostNotification: true,
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

func mapCommandUserToDomainUser(command *events.CreateUserCommand) *domain.User {
	user := domain.User{
		// Id: User.Id,
		Name:             command.User.Name,
		LastName:         command.User.LastName,
		MobileNumber:     command.User.MobileNumber,
		Username:         command.User.Username,
		Gender:           domain.GenderEnum(command.User.Gender),
		Birthday:         command.User.Birthday,
		Email:            command.User.Email,
		Biography:        command.User.Biography,
		IsPublic:         command.User.IsPublic,
		IsActive:         false, // nalog nije aktivan
		PostNotification: true,
	}
	for _, education := range command.User.Education {
		ed_id := primitive.NewObjectID()
		user.Education = append(user.Education, domain.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     domain.EducationEnum(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate,
			EndDate:   education.EndDate,
		})
	}

	for _, experience := range command.User.Experience {
		ex_id := primitive.NewObjectID()
		user.Experience = append(user.Experience, domain.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate,
			EndDate:   experience.EndDate,
		})
	}

	for _, skill := range command.User.Skills {
		s_id := primitive.NewObjectID()
		user.Skills = append(user.Skills, domain.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	for _, interest := range command.User.Interests {
		in_id := primitive.NewObjectID()
		user.Interests = append(user.Interests, domain.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}
	return &user
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
