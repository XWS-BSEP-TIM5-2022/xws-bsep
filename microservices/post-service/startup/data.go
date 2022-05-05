package startup

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var posts = []*domain.Post{
	{
		Id:          getObjectId("623b0cc3a34d25d8567f9f82"),
		Title:       "post 1",
		DateCreated: "22.02.2022.",
	},
	{
		Id:          getObjectId("623b0cc3a34d25d8567f9f83"),
		Title:       "post 2",
		DateCreated: "06.60.2021.",
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
