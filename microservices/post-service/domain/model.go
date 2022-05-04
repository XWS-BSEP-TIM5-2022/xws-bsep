package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	//Code string `bson:"code"`
	UserId string `bson:"user_id"`
}

type Post struct {
	Id          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	DateCreated string             `bson:"date_created"`
	//Likes       []Like             `bson:"likes"`

	//idUser
	//tekst, slika i linkovi
	//post se trajno nalazi na korisnikovom profilu
	//like, dislike, pisanje komentara
	//datum i vreme kreiranja ?
}
