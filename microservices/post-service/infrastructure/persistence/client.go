package persistence

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient(host, port string) (*mongo.Client, error) { // pomocna metoda koja vraca mongo klijenta na osnovu neophodnih podataka (host, port)
	uri := fmt.Sprintf("mongodb://%s:%s/", host, port)
	options := options.Client().ApplyURI(uri)
	return mongo.Connect(context.TODO(), options) // konektuje se i vraca klijenta nazad
}
