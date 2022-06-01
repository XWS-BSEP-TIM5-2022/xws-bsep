package persistence

import (
	"context"
	_ "context"
	"errors"
	_ "errors"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html"
	"net/url"
	"strings"
)

var validate *validator.Validate

const (
	DATABASE   = "post_db"
	COLLECTION = "post"
)

// ovde su implementirane sve metode iz post_store interfejsa

type PostMongoDBStore struct {
	posts *mongo.Collection
}

func NewPostMongoDBStore(client *mongo.Client) domain.PostStore {
	validate = validator.New()

	posts := client.Database(DATABASE).Collection(COLLECTION) // preuzimamo kolekciju proizvoda nad kojima radimo sve ostale operacije
	return &PostMongoDBStore{
		posts: posts,
	}
}

func (store *PostMongoDBStore) Get(id primitive.ObjectID) (*domain.Post, error) {
	/** Escape '$' - Prevent NoSQL Injection **/
	var convertedID = id.Hex()
	var lenId = len(convertedID)
	var checkId = ""
	for i := 0; i < lenId; i++ {
		char := string(convertedID[i])
		if char != "$" {
			checkId = checkId + char
		}
	}

	/** EscapeString **/
	newId, _ := primitive.ObjectIDFromHex(html.EscapeString(strings.TrimSpace(checkId)))
	filter := bson.M{"_id": newId}
	return store.filterOne(filter)
}

func (store *PostMongoDBStore) GetAll() ([]*domain.Post, error) {
	filter := bson.D{{}} // prazan filter, jer zelimo da dobavimo sve proizvode
	return store.filter(filter)
}

func (store *PostMongoDBStore) GetAllByUser(id string) ([]*domain.Post, error) {
	if len(id) != 24 {
		err := errors.New("id not valid")
		return nil, err
	}

	/** Escape '$' - Prevent NoSQL Injection **/
	var checkId = ""
	for i := 0; i < len(id); i++ {
		char := string(id[i])
		if char != "$" {
			checkId = checkId + char
		}
	}

	/** EscapeString - escapes special characters  <, >, &, ', " **/
	filter := bson.M{"user_id": html.EscapeString(checkId)}
	return store.filter(filter)
}

func (store *PostMongoDBStore) Insert(post *domain.Post) (string, error) {
	post.Id = primitive.NewObjectID()

	/** EscapeString **/
	post.Text = html.EscapeString(post.Text)

	// validate links
	for _, link := range post.Links {
		u, err := url.ParseRequestURI(link)
		if err != nil {
			fmt.Println("URL: ", u)
			return "error", err
		}
	}

	// validate post
	err := validate.Struct(post)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return "error", err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("---------------- pocetak greske ----------------")
			fmt.Println(err.Field())
			fmt.Println(err.Tag())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------- kraj greske ----------------")
		}
		return "error", err
	}

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
		"likes":        post.Likes,
		"dislikes":     post.Dislikes,
		"comments":     post.Comments,
		"user_id":      post.UserId,
	}}

	// validate likes
	for _, like := range post.Likes {
		err := validate.Struct(like)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err)
				return "error", err
			}

			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println("---------------- pocetak greske ----------------")
				fmt.Println(err.Field())
				fmt.Println(err.Tag())
				fmt.Println(err.Type())
				fmt.Println(err.Value())
				fmt.Println(err.Param())
				fmt.Println("---------------- kraj greske ----------------")
			}
			return "error", err
		}
	}

	// validate dislikes
	for _, dislike := range post.Dislikes {
		err := validate.Struct(dislike)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err)
				return "error", err
			}

			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println("---------------- pocetak greske ----------------")
				fmt.Println(err.Field())
				fmt.Println(err.Tag())
				fmt.Println(err.Type())
				fmt.Println(err.Value())
				fmt.Println(err.Param())
				fmt.Println("---------------- kraj greske ----------------")
			}
			return "error", err
		}
	}

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

func (store *PostMongoDBStore) LikePost(post *domain.Post, user_id string) (string, error) {
	like := domain.Like{}
	like.Id = primitive.NewObjectID()
	like.UserId = user_id

	// validate like
	err := validate.Struct(like)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return "error", err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("---------------- pocetak greske ----------------")
			fmt.Println(err.Field())
			fmt.Println(err.Tag())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------- kraj greske ----------------")
		}
		return "error", err
	}

	post.Likes = append(post.Likes, like)
	newData := bson.M{"$set": bson.M{
		"text":         post.Text,
		"date_created": post.DateCreated,
		"images":       post.Images,
		"links":        post.Links,
		"likes":        post.Likes,
		"dislikes":     post.Dislikes,
		"comments":     post.Comments,
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

func (store *PostMongoDBStore) DislikePost(post *domain.Post, user_id string) (string, error) {
	dislike := domain.Dislike{}
	dislike.Id = primitive.NewObjectID()
	dislike.UserId = user_id

	// validate dislike
	err := validate.Struct(dislike)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return "error", err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("---------------- pocetak greske ----------------")
			fmt.Println(err.Field())
			fmt.Println(err.Tag())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------- kraj greske ----------------")
		}
		return "error", err
	}

	post.Dislikes = append(post.Dislikes, dislike)
	newData := bson.M{"$set": bson.M{
		"text":         post.Text,
		"date_created": post.DateCreated,
		"images":       post.Images,
		"links":        post.Links,
		"likes":        post.Likes,
		"dislikes":     post.Dislikes,
		"comments":     post.Comments,
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

func (store *PostMongoDBStore) CommentPost(post *domain.Post, user_id string, text string) (string, error) {
	comment := domain.Comment{}
	comment.Id = primitive.NewObjectID()
	comment.UserId = user_id
	/** EscapeString **/
	comment.Text = html.EscapeString(text)

	// validate comment
	err := validate.Struct(comment)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return "error", err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("---------------- pocetak greske ----------------")
			fmt.Println(err.Field())
			fmt.Println(err.Tag())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println("---------------- kraj greske ----------------")
		}
		return "error", err
	}

	post.Comments = append(post.Comments, comment)
	newData := bson.M{"$set": bson.M{
		"text":         post.Text,
		"date_created": post.DateCreated,
		"images":       post.Images,
		"links":        post.Links,
		"likes":        post.Likes,
		"dislikes":     post.Dislikes,
		"comments":     post.Comments,
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

func (store *PostMongoDBStore) UpdateCompanyInfo(company *domain.Company, oldName string) (string, error) {
	newData := bson.M{
		"id":          primitive.NewObjectID(),
		"name":        company.Name,
		"description": company.Description,
		"phoneNumber": company.PhoneNumber,
		"is_active":   company.IsActive,
	}

	newCompany := bson.M{"$set": bson.M{
		"company": newData,
	}}

	posts, _ := store.GetAll()
	opts := options.Update().SetUpsert(true)

	for _, post := range posts {
		_, err := store.posts.UpdateOne(context.TODO(), bson.M{"_id": post.Id}, newCompany, opts)
		fmt.Println("azurirao")
		if err != nil {
			return "error", err
		}
	}

	return "success", nil
}
