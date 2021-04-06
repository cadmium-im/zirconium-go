package core

import (
	"context"
	"github.com/cadmium-im/zirconium-go/core/models"
	"github.com/cadmium-im/zirconium-go/core/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UsersCollectionName = "users"
)

type UserManager struct {
	usersCol *mongo.Collection
}

func NewUserManager(db *mongo.Database) (*UserManager, error) {
	col := db.Collection(UsersCollectionName)

	um := &UserManager{
		usersCol: col,
	}

	err := um.initMongo()
	if err != nil {
		return nil, err
	}

	return um, nil
}

func (um *UserManager) initMongo() error {
	usernameIndex := mongo.IndexModel{
		Keys: bson.M{
			"username": 1,
		}, Options: options.Index().SetUnique(true),
	}

	exists, err := utils.IsCollectionExists(context.Background(), um.usersCol.Database(), um.usersCol.Name())
	if err != nil {
		return err
	}
	if exists {
		_, err := um.usersCol.Indexes().DropAll(context.Background())
		if err != nil {
			return err
		}
	}

	_, err = um.usersCol.Indexes().CreateMany(context.Background(), []mongo.IndexModel{usernameIndex})
	return err
}

func (um *UserManager) SaveUser(user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.UUID = uuid.New().String()
	_, err := um.usersCol.InsertOne(context.Background(), user)
	return err
}

func (um *UserManager) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := um.usersCol.FindOne(context.Background(), bson.D{{"username", username}}).Decode(&user)
	return &user, err
}
