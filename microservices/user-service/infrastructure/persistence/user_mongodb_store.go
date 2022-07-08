package persistence

import (
	"context"
	_ "context"
	"errors"
	_ "errors"
	"fmt"
	"log"
	"strings"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate *validator.Validate

const (
	DATABASE   = "user"
	COLLECTION = "user"
)

type UserMongoDBStore struct {
	users *mongo.Collection
}

func NewUserMongoDBStore(client *mongo.Client) domain.UserStore {
	validate = validator.New()

	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserMongoDBStore{
		users: users,
	}
}

func (store *UserMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "Get database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByUsername database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByEmail(email string) (*domain.User, error) {
	filter := bson.M{"email": email}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) Insert(user *domain.User) (*domain.User, error) {

	user.Id = primitive.NewObjectID()

	// checkUsername, _ := store.GetByUsername(user.Username)

	// if checkUsername != nil {
	// 	return nil, errors.New("username already exists")
	// }

	checkEmail, _ := store.GetByEmail(user.Email)

	if checkEmail != nil {
		return nil, errors.New("email already exists")
	}
	fmt.Printf("var1 = %T\n", user.IsPublic)

	err := validate.Struct(user)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return nil, err
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
		return nil, err
	}

	result, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (store *UserMongoDBStore) GetAll(ctx context.Context) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	filter := bson.D{{"is_active", true}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) GetAllPublic(ctx context.Context) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllPublic database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	filter := bson.D{{"is_public", true}, {"is_active", true}, {"role", "User"}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) DeleteAll() {
	store.users.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *UserMongoDBStore) Update(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Update database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	oldData := bson.M{"_id": user.Id}

	newData := bson.M{"$set": bson.M{
		"name":          user.Name,
		"last_name":     user.LastName,
		"mobile_number": user.MobileNumber,
		"gender":        user.Gender,
		"birthday":      user.Birthday,
		"email":         user.Email,
		"biography":     user.Biography,
		"is_public":     user.IsPublic,
		"education":     user.Education,
		"experience":    user.Experience,
		"skills":        user.Skills,
		"interests":     user.Interests,
		"username":      user.Username,
	}}

	oldUser, _ := store.filterOne(oldData)

	if oldUser != nil && user.Username != "" && user.Username != oldUser.Username {

		checkUsername, _ := store.GetByUsername(ctx, user.Username)

		if checkUsername != nil {
			return "username already exists", errors.New("username already exists")
		}
	}

	if oldUser != nil && user.Email != "" && user.Email != oldUser.Email {

		checkEmail, _ := store.GetByEmail(user.Email)

		if checkEmail != nil {
			return "email already exists", errors.New("email already exists")
		}
	}

	opts := options.Update().SetUpsert(true)

	result, err := store.users.UpdateOne(context.TODO(), oldData, newData, opts)

	if err != nil {
		return "error", err
	}
	if result.MatchedCount != 1 {
		return "one document should've been updated", errors.New("one document should've been updated")
	}
	return "success", nil

}

func (store *UserMongoDBStore) Search(ctx context.Context, criteria string) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "Search database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	criteria = strings.ToLower(criteria)
	criteria = strings.TrimSpace(criteria)
	words := strings.Split(criteria, " ")

	var ret []*domain.User

	users, err := store.GetAllPublic(ctx)

	if err != nil {
		return nil, err
	}

	for _, word := range words {

		for _, user := range users {
			if user.IsActive {
				name := strings.ToLower(user.Name)
				lastName := strings.ToLower(user.LastName)
				email := strings.ToLower(user.Email)

				if strings.Contains(name, word) || strings.Contains(lastName, word) || strings.Contains(email, word) {
					ret = append(ret, user)
				}
			}
		}
	}

	return ret, nil

}

func (store *UserMongoDBStore) filterOne(filter interface{}) (user *domain.User, err error) {
	result := store.users.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (store *UserMongoDBStore) filter(filter interface{}) ([]*domain.User, error) {
	cursor, err := store.users.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (users []*domain.User, err error) {
	for cursor.Next(context.TODO()) {
		var user domain.User
		err = cursor.Decode(&user)
		if err != nil {
			return
		}
		users = append(users, &user)
	}
	err = cursor.Err()
	return
}

func (store *UserMongoDBStore) GetById(ctx context.Context, hexId string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetById database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	id, err := hexIdToId(hexId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	user, err := store.filterOne(filter)
	if !user.IsActive {
		return nil, errors.New("User is not actived")
	}
	return user, nil
}

func hexIdToId(hexId string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(hexId)
}

func (store *UserMongoDBStore) UpdateBasicInfo(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateBasicInfo database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	oldData := bson.M{"_id": user.Id}

	newData := bson.M{"$set": bson.M{
		"name":          user.Name,
		"last_name":     user.LastName,
		"mobile_number": user.MobileNumber,
		"gender":        user.Gender,
		"birthday":      user.Birthday,
		"email":         user.Email,
		"biography":     user.Biography,
		"is_public":     user.IsPublic,
	}}

	oldUser, _ := store.filterOne(oldData)

	if oldUser != nil && user.Email != "" && user.Email != oldUser.Email {

		checkEmail, _ := store.GetByEmail(user.Email)

		if checkEmail != nil {
			return "email already exists", errors.New("email already exists")
		}
	}

	opts := options.Update().SetUpsert(true)

	result, err := store.users.UpdateOne(context.TODO(), oldData, newData, opts)

	if err != nil {
		return "error", err
	}
	if result.MatchedCount != 1 {
		return "one document should've been updated", errors.New("one document should've been updated")
	}
	return "success", nil

}

func (store *UserMongoDBStore) UpdateExperienceAndEducation(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateExperienceAndEducation database store")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	oldData := bson.M{"_id": user.Id}

	newData := bson.M{"$set": bson.M{

		"education":  user.Education,
		"experience": user.Experience,
	}}

	result, err := store.users.UpdateOne(context.TODO(), oldData, newData)

	if err != nil {
		return "error", err
	}
	if result.MatchedCount != 1 {
		return "one document should've been updated", errors.New("one document should've been updated")
	}
	return "success", nil

}

func (store *UserMongoDBStore) UpdateSkillsAndInterests(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateSkillsAndInterests database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	oldData := bson.M{"_id": user.Id}

	newData := bson.M{"$set": bson.M{

		"skills":    user.Skills,
		"interests": user.Interests,
	}}

	result, err := store.users.UpdateOne(context.TODO(), oldData, newData)

	if err != nil {
		return "error", err
	}
	if result.MatchedCount != 1 {
		return "one document should've been updated", errors.New("one document should've been updated")
	}
	return "success", nil

}

func (store *UserMongoDBStore) UpdateIsActiveById(ctx context.Context, userId string) error {
	span := tracer.StartSpanFromContext(ctx, "UpdateIsActiveById database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		panic(err)
	}

	oldData := bson.M{"_id": objID}
	newData := bson.M{"$set": bson.M{
		"is_active": true,
	}}

	opts := options.Update().SetUpsert(true)
	result, err := store.users.UpdateOne(context.TODO(), oldData, newData, opts)
	if err != nil || result.ModifiedCount == 0 {
		return err
	}
	return nil
}

func (store *UserMongoDBStore) GetIdByEmail(ctx context.Context, email string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIdByEmail database store")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"email": email}
	user, err := store.filterOne(filter)
	if err != nil {
		return "", err
	}
	if !user.IsActive {
		return "", errors.New("User is not actived")
	}
	return user.Id.Hex(), nil
}

func (store *UserMongoDBStore) DeleteUser(userId, email string) error {
	filter := bson.M{"email": email}
	res, err := store.users.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println("User is not deleted, err: ", err)
		return err
	}
	log.Println("Deleted count: ", res.DeletedCount)
	return nil
}
