package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type MessageStore interface {
	GetConversation(sender, receiver string) (*Conversation, error)
	GetAllConversationsForUser(user string) ([]*Conversation, error)
	NewMessage(message *Message, sender string) (*Conversation, error)
	GetConversationById(id primitive.ObjectID) (*Conversation, error)
}
