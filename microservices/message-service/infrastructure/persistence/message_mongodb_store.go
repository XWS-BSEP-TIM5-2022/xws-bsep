package persistence

import (
	"context"
	_ "context"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
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
	DATABASE   = "message_db"
	COLLECTION = "message"
)

type MessageMongoDBStore struct {
	messages *mongo.Collection
}

func (store MessageMongoDBStore) GetConversation(ctx context.Context, sender, receiver string) (*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConversation database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{
		"$or": []bson.M{
			{
				"user1": sender,
				"user2": receiver,
			},
			{
				"user1": receiver,
				"user2": sender,
			},
		}}

	return store.filterOne(filter)

}

func (store MessageMongoDBStore) GetAllConversationsForUser(ctx context.Context, user string) ([]*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConversationsForUser database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{
		"$or": []bson.M{
			{"user1": user},
			{"user2": user},
		}}

	return store.filter(filter)
}

func (store MessageMongoDBStore) NewMessage(ctx context.Context, message *domain.Message, sender string) (*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "NewMessage database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	conversation, err := store.GetConversation(ctx, sender, message.Receiver)
	if conversation == nil {

		//senderId, _ := primitive.ObjectIDFromHex(sender)
		newConversation := &domain.Conversation{
			Id:       primitive.NewObjectID(),
			User1:    sender,
			User2:    message.Receiver,
			Messages: nil,
		}

		newConversation.Messages = append(newConversation.Messages, *message)
		_, err = store.messages.InsertOne(context.TODO(), newConversation)

		if err != nil {
			return nil, err
		}

		conversation, err = store.GetConversationById(ctx, newConversation.Id)

	} else {

		messages := append(conversation.Messages, *message)

		oldData := bson.M{"_id": conversation.Id}

		newData := bson.M{"$set": bson.M{

			"messages": messages,
		}}

		_, err = store.messages.UpdateOne(context.TODO(), oldData, newData)

		if err != nil {
			return nil, err
		}

		conversation, err = store.GetConversationById(ctx, conversation.Id)

	}

	return conversation, err
}

func (store MessageMongoDBStore) GetConversationById(ctx context.Context, id primitive.ObjectID) (*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConversationById database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

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

func hexIdToId(hexId string) primitive.ObjectID {
	ret, _ := primitive.ObjectIDFromHex(hexId)
	return ret
}

func NewMessageMongoDBStore(client *mongo.Client) domain.MessageStore {
	validate = validator.New()

	messages := client.Database(DATABASE).Collection(COLLECTION)
	return &MessageMongoDBStore{
		messages: messages,
	}
}
