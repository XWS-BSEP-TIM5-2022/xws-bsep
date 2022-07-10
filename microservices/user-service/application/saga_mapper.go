package application

import (
	events "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/create_user"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
)

func mapNewUser(user *domain.User, username, password string) *events.CreateUserCommand {
	event := &events.CreateUserCommand{
		Type: events.CreateUser,
		User: events.UserDetails{
			// Id:           userId,
			Name:                user.Name,
			LastName:            user.LastName,
			Email:               user.Email,
			MobileNumber:        user.MobileNumber,
			Gender:              events.GenderEnum(user.Gender),
			Birthday:            user.Birthday,
			Username:            username,
			Biography:           user.Biography,
			Role:                user.Role,
			Password:            password,
			PostNotification:    true,
			MessageNotification: true,
			FollowNotification:  true,
		},
	}

	for _, edu := range user.Education {
		eventItem := events.Education{
			Id:        edu.Id.Hex(),
			Name:      edu.Name,
			Level:     events.EducationEnum(edu.Level),
			Place:     edu.Place,
			StartDate: edu.StartDate,
			EndDate:   edu.EndDate,
		}
		event.User.Education = append(event.User.Education, eventItem)
	}

	for _, edu := range user.Experience {
		eventItem2 := events.Experience{
			Id:        edu.Id.Hex(),
			Name:      edu.Name,
			Headline:  edu.Headline,
			Place:     edu.Place,
			StartDate: edu.StartDate,
			EndDate:   edu.EndDate,
		}
		event.User.Experience = append(event.User.Experience, eventItem2)
	}
	for _, edu := range user.Interests {
		eventItem3 := events.Interest{
			Id:          edu.Id.Hex(),
			Name:        edu.Name,
			Description: edu.Description,
		}
		event.User.Interests = append(event.User.Interests, eventItem3)
	}
	for _, edu := range user.Skills {
		eventItem3 := events.Skill{
			Id:   edu.Id.Hex(),
			Name: edu.Name,
		}
		event.User.Skills = append(event.User.Skills, eventItem3)
	}

	return event
}
