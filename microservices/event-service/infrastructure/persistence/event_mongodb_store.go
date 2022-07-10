package persistence

import (
	"context"
	_ "context"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate *validator.Validate

const (
	DATABASE   = "event_db"
	COLLECTION = "event"
)

type EventMongoDBStore struct {
	events *mongo.Collection
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
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (events []*domain.Event, err error) {
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

	events := client.Database(DATABASE).Collection(COLLECTION)
	return &EventMongoDBStore{
		events: events,
	}
}
