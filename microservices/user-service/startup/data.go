package startup

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var users = []*domain.User{
	{
		Id:           getObjectId("623b0cc336a1d6fd8c1cf0f6"),
		Name:         "Ranko",
		LastName:     "Rankovic",
		MobileNumber: "0653829384",
		Gender:       domain.Male,
		Birthday:     time.Date(1997, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
		Email:        "ranko@gmail.com",
		Biography:    "Vredan od malih nogu",
		Username:     "rankoRankovic",
		Password:     "ranko",
		IsPublic:     true,
		Education:    []domain.Education{},
		Experience:   []domain.Experience{},
		Skills: []domain.Skill{
			{Id: getObjectId("623b0cc336a1d6fd8c4cf0f6"), Name: "Java"},
		},
		Interests: []domain.Interest{},
	},
	{
		Id:           getObjectId("623b4ac336a1d6fd8c1cf0f6"),
		Name:         "Marko",
		LastName:     "Markovic",
		MobileNumber: "06538293354",
		Gender:       domain.Male,
		Birthday:     time.Date(1967, time.Month(8), 21, 1, 10, 30, 0, time.UTC),
		Email:        "markic@gmail.com",
		Biography:    "Rodjen u Novom Sadu",
		Username:     "markooom",
		Password:     "marko",
		IsPublic:     true,
		Education: []domain.Education{
			{
				Id:        getObjectId("643b0cc136a1d6fd8c1cf0f6"),
				Name:      "OS ,,NikolaTesla''",
				Level:     domain.Primary,
				Place:     "Backa Topola",
				StartDate: time.Date(1974, time.Month(9), 1, 1, 10, 30, 0, time.UTC),
				EndDate:   time.Date(1982, time.Month(6), 15, 1, 10, 30, 0, time.UTC),
			},
			{
				Id:        getObjectId("642b2cc136a1d6fd8c1cf0f6"),
				Name:      "Gimnazija",
				Level:     domain.Secondary,
				Place:     "Subotica",
				StartDate: time.Date(1982, time.Month(9), 1, 1, 10, 30, 0, time.UTC),
				EndDate:   time.Date(1986, time.Month(6), 15, 1, 10, 30, 0, time.UTC),
			},
		},
		Experience: []domain.Experience{},
		Skills:     []domain.Skill{},
		Interests:  []domain.Interest{},
	},
	{
		Id:           getObjectId("623b4ab326a1d6fd8c1cf0f6"),
		Name:         "Jana",
		LastName:     "Markovic",
		MobileNumber: "06532293354",
		Gender:       domain.Male,
		Birthday:     time.Date(1969, time.Month(3), 11, 1, 10, 30, 0, time.UTC),
		Email:        "markovic@gmail.com",
		Biography:    "Brat Marko",
		Username:     "janamarkovic",
		Password:     "jana",
		IsPublic:     false,
		Education:    []domain.Education{},
		Experience: []domain.Experience{
			{
				Id:        getObjectId("623b4ab326a1d6fd8c3cf0f5"),
				Name:      "Synechron",
				Headline:  "backend developer",
				Place:     "Novi Sad",
				StartDate: time.Date(1982, time.Month(9), 1, 1, 10, 30, 0, time.UTC),
				EndDate:   time.Date(1986, time.Month(6), 15, 1, 10, 30, 0, time.UTC),
			},
		},
		Skills:    []domain.Skill{},
		Interests: []domain.Interest{},
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
