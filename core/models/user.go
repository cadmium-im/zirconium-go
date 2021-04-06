package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"-" bson:"_id"`
	UUID         string             `json:"uuid"`
	Username     string             `json:"username"`
	PasswordHash string             `json:"-" bson:"passwordHash"`
}
