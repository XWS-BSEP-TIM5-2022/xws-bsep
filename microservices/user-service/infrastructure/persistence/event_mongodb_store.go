package persistence

import (
	"context"
	_ "context"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DatabaseEvent   = "event_db"
	CollectionEvent = "event"
)

type EventMongoDBStore struct {
	events *mongo.Collection
}

func (store EventMongoDBStore) NewEvent(event *domain.Event) (*domain.Event, error) {

	event.Id = primitive.NewObjectID()

	result, err := store.events.InsertOne(context.TODO(), event)
	if err != nil {
		return nil, err
	}
	event.Id = result.InsertedID.(primitive.ObjectID)
	return event, nil

}

func (store EventMongoDBStore) GetAllEvents() ([]*domain.Event, error) {

	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *EventMongoDBStore) filter(filter interface{}) ([]*domain.Event, error) {
	cursor, err := store.events.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decodeEvent(cursor)
}

func decodeEvent(cursor *mongo.Cursor) (events []*domain.Event, err error) {
	for cursor.Next(context.TODO()) {
		var event domain.Event
		err = cursor.Decode(&event)
		if err != nil {
			return
		}
		events = append(events, &event)
	}
	err = cursor.Err()
	return
}

func (store *EventMongoDBStore) filterOne(filter interface{}) (event *domain.Event, err error) {
	result := store.events.FindOne(context.TODO(), filter)
	err = result.Decode(&event)
	return
}

func NewEventMongoDBStore(client *mongo.Client) domain.EventStore {
	validate = validator.New()

	events := client.Database(DatabaseEvent).Collection(CollectionEvent)
	return &EventMongoDBStore{
		events: events,
	}
}
