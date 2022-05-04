package startup

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var users = []*domain.User{
	{
		Id:       getObjectId("623b0cc336a1d6fd8c1cf0f6"),
		Name:     "Ranko",
		LastName: "Rankovic",
	},
	{
		Id:       getObjectId("623b4ac336a1d6fd8c1cf0f6"),
		Name:     "Marko",
		LastName: "Markovic",
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
