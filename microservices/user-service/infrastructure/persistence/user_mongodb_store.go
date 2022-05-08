package persistence

import (
	"context"
	_ "context"
	"errors"
	_ "errors"

	"strings"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "user"
	COLLECTION = "user"
)

type UserMongoDBStore struct {
	users *mongo.Collection
}

func NewUserMongoDBStore(client *mongo.Client) domain.UserStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserMongoDBStore{
		users: users,
	}
}

func (store *UserMongoDBStore) Get(id primitive.ObjectID) (*domain.User, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByUsername(username string) (*domain.User, error) {
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

	result, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (store *UserMongoDBStore) GetAll() ([]*domain.User, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) GetAllPublic() ([]*domain.User, error) {
	filter := bson.D{{"is_public", true}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) DeleteAll() {
	store.users.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *UserMongoDBStore) Update(user *domain.User) (string, error) {

	oldData := bson.M{"_id": user.Id}

	newData := bson.M{"$set": bson.M{
		"name":          user.Name,
		"last_name":     user.LastName,
		"mobile_number": user.MobileNumber,
		"gender":        user.Gender,
		"birthday":      user.Birthday,
		"email":         user.Email,
		"biography":     user.Biography,
		// "username":      user.Username,
		// "password":      user.Password,
		"is_public":  user.IsPublic,
		"education":  user.Education,
		"experience": user.Experience,
		"skills":     user.Skills,
		"interests":  user.Interests,
	}}

	oldUser, _ := store.filterOne(oldData)

	// if oldUser != nil && user.Username != "" && user.Username != oldUser.Username {

	// 	checkUsername, _ := store.GetByUsername(user.Username)

	// 	if checkUsername != nil {
	// 		return "username already exists", errors.New("username already exists")
	// 	}
	// }

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

func (store *UserMongoDBStore) Search(criteria string) ([]*domain.User, error) {
	criteria = strings.ToLower(criteria)
	criteria = strings.TrimSpace(criteria)
	words := strings.Split(criteria, " ")

	var ret []*domain.User

	users, err := store.GetAllPublic()

	if err != nil {
		return nil, err
	}

	for _, word := range words {

		for _, user := range users {

			name := strings.ToLower(user.Name)
			lastName := strings.ToLower(user.LastName)
			// username := strings.ToLower(user.Username)

			if strings.Contains(name, word) || strings.Contains(lastName, word) /*|| strings.Contains(username, word)*/ {
				ret = append(ret, user)
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

func (store *UserMongoDBStore) GetById(hexId string) (*domain.User, error) {
	id, err := hexIdToId(hexId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func hexIdToId(hexId string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(hexId)
}
