package persistence

import (
	"context"
	_ "context"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate *validator.Validate

const (
	DATABASE   = "notification"
	COLLECTION = "notification"
)

type NotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func (store NotificationMongoDBStore) GetById(id primitive.ObjectID) (*domain.Notification, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *NotificationMongoDBStore) filter(filter interface{}) ([]*domain.Notification, error) {
	cursor, err := store.notifications.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (notifications []*domain.Notification, err error) {
	for cursor.Next(context.TODO()) {
		var notification domain.Notification
		err = cursor.Decode(&notification)
		if err != nil {
			return
		}
		notifications = append(notifications, &notification)
	}
	err = cursor.Err()
	return
}

func (store *NotificationMongoDBStore) filterOne(filter interface{}) (notification *domain.Notification, err error) {
	result := store.notifications.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}

func hexIdToId(hexId string) primitive.ObjectID {
	ret, _ := primitive.ObjectIDFromHex(hexId)
	return ret
}

func NewNotificationMongoDBStore(client *mongo.Client) domain.NotificationStore {
	validate = validator.New()

	notifications := client.Database(DATABASE).Collection(COLLECTION)
	return &NotificationMongoDBStore{
		notifications: notifications,
	}
}
