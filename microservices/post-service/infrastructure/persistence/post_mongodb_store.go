package persistence

import (
	"context"
	_ "context"
	"errors"
	_ "errors"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
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

func (store *PostMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "Get database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

func (store *PostMongoDBStore) GetAll(ctx context.Context) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}} // prazan filter, jer zelimo da dobavimo sve proizvode
	return store.filter(filter)
}

func (store *PostMongoDBStore) GetAllByUser(ctx context.Context, id string) ([]*domain.Post, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUser database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

func (store *PostMongoDBStore) Insert(ctx context.Context, post *domain.Post) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	post.Id = primitive.NewObjectID()

	/** EscapeString **/
	post.Text = html.EscapeString(post.Text)

	// validacija - da li su neka polja iz job offer-a nepopunjena
	if post.IsJobOffer {
		if strings.TrimSpace(post.JobOffer.JobDescription) == "" || strings.TrimSpace(post.JobOffer.DailyActivities) == "" ||
			strings.TrimSpace(post.JobOffer.Preconditions) == "" || strings.TrimSpace(post.Company.Name) == "" ||
			strings.TrimSpace(post.Company.Description) == "" {
			return "error", nil
		}

		post.JobOffer.JobDescription = html.EscapeString(post.JobOffer.JobDescription)
		post.JobOffer.DailyActivities = html.EscapeString(post.JobOffer.DailyActivities)
		post.JobOffer.Preconditions = html.EscapeString(post.JobOffer.Preconditions)
		post.Company.Name = html.EscapeString(post.Company.Name)
		post.Company.Description = html.EscapeString(post.Company.Description)
	}

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

//func decodeImage(path string) (primitive.Binary, error) {
//	image, err := base64.StdEncoding.DecodeString(path)
//	if err != nil {
//		return primitive.Binary{}, err
//	}
//	return primitive.Binary{Data: image}, nil
//}

func (store *PostMongoDBStore) Update(ctx context.Context, post *domain.Post) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Update database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	newData := bson.M{"$set": bson.M{
		"text":         post.Text,
		"date_created": post.DateCreated,
		"image":        post.Image,
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

func (store *PostMongoDBStore) LikePost(ctx context.Context, post *domain.Post, user_id string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "LikePost database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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
		"image":        post.Image,
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

func (store *PostMongoDBStore) DislikePost(ctx context.Context, post *domain.Post, user_id string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "DislikePost database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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
		"image":        post.Image,
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

func (store *PostMongoDBStore) CommentPost(ctx context.Context, post *domain.Post, user_id string, text string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "CommentPost database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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
		"image":        post.Image,
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

func (store *PostMongoDBStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DeleteAll database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

func (store *PostMongoDBStore) UpdateCompanyInfo(ctx context.Context, company *domain.Company, oldName string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateCompanyInfor database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

	posts, _ := store.GetAll(ctx)
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
