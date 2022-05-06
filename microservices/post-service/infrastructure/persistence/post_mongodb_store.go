package persistence

import (
	"context"
	_ "context"
	"errors"
	_ "errors"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "post_db"
	COLLECTION = "post"
)

// ovde su implementirane sve metode iz post_store interfejsa

type PostMongoDBStore struct {
	posts *mongo.Collection
}

func NewPostMongoDBStore(client *mongo.Client) domain.PostStore {
	posts := client.Database(DATABASE).Collection(COLLECTION) // preuzimamo kolekciju proizvoda nad kojima radimo sve ostale operacije
	return &PostMongoDBStore{
		posts: posts,
	}
}

func (store *PostMongoDBStore) Get(id primitive.ObjectID) (*domain.Post, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *PostMongoDBStore) GetAll() ([]*domain.Post, error) {
	filter := bson.D{{}} // prazan filter, jer zelimo da dobavimo sve proizvode
	return store.filter(filter)
}

func (store *PostMongoDBStore) GetAllByUser(id string) ([]*domain.Post, error) {
	filter := bson.M{"user_id": id}
	return store.filter(filter)
}

func (store *PostMongoDBStore) Insert(post *domain.Post) (string, error) {
	post.Id = primitive.NewObjectID()
	result, err := store.posts.InsertOne(context.TODO(), post)
	if err != nil {
		return "error", err
	}
	post.Id = result.InsertedID.(primitive.ObjectID)
	return "success", nil
}

func (store *PostMongoDBStore) Update(post *domain.Post) (string, error) {
	newData := bson.M{"$set": bson.M{
		"text":         post.Text,
		"date_created": post.DateCreated,
		"images":       post.Images,
		"links":        post.Links,
		"likes":        post.Likes,    // TODO
		"dislikes":     post.Dislikes, // TODO
		"comments":     post.Comments, // TODO
		"user_id":      post.UserId,
	}}

	opts := options.Update().SetUpsert(true)
	result, err := store.posts.UpdateOne(context.TODO(), bson.M{"_id": post.Id}, newData, opts)

	if err != nil {
		return "error", err
	}
	if result.MatchedCount != 1 {
		return "one document should've been updated", errors.New("one document should've been updated")
	}
	return "success", nil
}

func (store *PostMongoDBStore) DeleteAll() {
	store.posts.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *PostMongoDBStore) filter(filter interface{}) ([]*domain.Post, error) {
	cursor, err := store.posts.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *PostMongoDBStore) filterOne(filter interface{}) (post *domain.Post, err error) {
	result := store.posts.FindOne(context.TODO(), filter)
	err = result.Decode(&post)
	return
}

func decode(cursor *mongo.Cursor) (posts []*domain.Post, err error) {
	for cursor.Next(context.TODO()) {
		var post domain.Post
		err = cursor.Decode(&post)
		if err != nil {
			return
		}
		posts = append(posts, &post)
	}
	err = cursor.Err()
	return
}
