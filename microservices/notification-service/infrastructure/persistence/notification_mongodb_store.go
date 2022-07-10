package persistence

import (
	"context"
	_ "context"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
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
	DATABASE   = "notification_db"
	COLLECTION = "notification"
)

type NotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func (store *NotificationMongoDBStore) GetAllByUser(ctx context.Context, id string) ([]*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUser database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"user_id": id}
	return store.filter(filter)
}

func (store NotificationMongoDBStore) GetAll(ctx context.Context) ([]*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filter(filter)
}

func (store NotificationMongoDBStore) Insert(ctx context.Context, post *domain.Notification) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.notifications.InsertOne(context.TODO(), post)
	if err != nil {
		return "error", err
	}
	post.Id = result.InsertedID.(primitive.ObjectID)
	return "success", nil
}

func (store NotificationMongoDBStore) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetById database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

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
