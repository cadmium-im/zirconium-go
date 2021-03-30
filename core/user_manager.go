package core

import "go.mongodb.org/mongo-driver/mongo"

const (
	UsersCollectionName = "users"
)

type UserManager struct {
	usersCol *mongo.Collection
}

func NewUserManager(db *mongo.Database) *UserManager {
	col := db.Collection(UsersCollectionName)

	um := &UserManager{
		usersCol: col,
	}

	return um
}
