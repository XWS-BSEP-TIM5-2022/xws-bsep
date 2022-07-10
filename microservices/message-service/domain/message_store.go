package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageStore interface {
	GetConversation(ctx context.Context, sender, receiver string) (*Conversation, error)
	GetAllConversationsForUser(ctx context.Context, user string) ([]*Conversation, error)
	NewMessage(ctx context.Context, message *Message, sender string) (*Conversation, error)
	GetConversationById(ctx context.Context, id primitive.ObjectID) (*Conversation, error)
}
