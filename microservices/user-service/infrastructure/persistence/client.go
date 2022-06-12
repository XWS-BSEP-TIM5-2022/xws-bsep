package persistence

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient(host, port string) (*mongo.Client, error) {
	//dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	//// dsn := fmt.Sprintf("host=localhost user=postgres password=admin dbname=NOVA port=5432 sslmode=disable")
	//return gorm.Open(postgres.Open(dsn), &gorm.Config{})

	uri := fmt.Sprintf("mongodb://%s:%s/", host, port)
	options := options.Client().ApplyURI(uri)
	return mongo.Connect(context.TODO(), options)
}
