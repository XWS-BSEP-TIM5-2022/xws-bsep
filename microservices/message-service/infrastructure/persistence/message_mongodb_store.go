package persistence

import (
	"context"
	_ "context"
	"errors"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate *validator.Validate

const (
	DATABASE   = "message"
	COLLECTION = "message"
)

type MessageMongoDBStore struct {
	messages *mongo.Collection
}

func (store MessageMongoDBStore) GetConversation(sender, receiver string) (*domain.Conversation, error) {

	filter := bson.M{
		"$or": []bson.M{
			{"$and": []bson.M{
				{"user_1": sender},
				{"user_2": receiver},
			}},
			{"$and": []bson.M{
				{"user_1": receiver},
				{"user_2": sender},
			}},
		},
	}

	return store.filterOne(filter)

}

func (store MessageMongoDBStore) GetAllConversationsForUser(user string) ([]*domain.Conversation, error) {

	filter := bson.M{
		"$or": []bson.M{
			{"user_1": user},
			{"user_2": user},
		}}

	return store.filter(filter)
}

func (store MessageMongoDBStore) NewMessage(message *domain.Message, sender string) (*domain.Conversation, error) {

	conversation, err := store.GetConversation(sender, message.Receiver.Hex())

	if conversation == nil {

		senderId, _ := primitive.ObjectIDFromHex(sender)
		newConversation := &domain.Conversation{
			Id:       primitive.NewObjectID(),
			User1:    senderId,
			User2:    message.Receiver,
			Messages: nil,
		}

		added, err := store.messages.InsertOne(context.TODO(), newConversation)

		if err != nil {
			return nil, err
		}

		filter := bson.M{"_id": added.InsertedID.(primitive.ObjectID)}
		conversation, err = store.filterOne(filter)

	} else {

		messages := append(conversation.Messages, *message)

		oldData := bson.M{"_id": conversation.Id}

		newData := bson.M{"$set": bson.M{

			"messages": messages,
		}}
		//bson.D{
		//	{"$set", bson.D{{"messages", messages}}

		result, err := store.messages.UpdateOne(context.TODO(), oldData, newData)

		if err != nil {
			return nil, err
		}

		if result.MatchedCount != 1 {
			return nil, errors.New("one document should've been updated")
		}
	}

	conversation.Messages = append(conversation.Messages, *message)
	return conversation, err
}

func (store MessageMongoDBStore) GetConversationById(id primitive.ObjectID) (*domain.Conversation, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *MessageMongoDBStore) filter(filter interface{}) ([]*domain.Conversation, error) {
	cursor, err := store.messages.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (conversations []*domain.Conversation, err error) {
	for cursor.Next(context.TODO()) {
		var conversation domain.Conversation
		err = cursor.Decode(&conversation)
		if err != nil {
			return
		}
		conversations = append(conversations, &conversation)
	}
	err = cursor.Err()
	return
}

func (store *MessageMongoDBStore) filterOne(filter interface{}) (conversation *domain.Conversation, err error) {
	result := store.messages.FindOne(context.TODO(), filter)
	err = result.Decode(&conversation)
	return
}

func NewMessageMongoDBStore(client *mongo.Client) domain.MessageStore {
	validate = validator.New()

	messages := client.Database(DATABASE).Collection(COLLECTION)
	return &MessageMongoDBStore{
		messages: messages,
	}
}
