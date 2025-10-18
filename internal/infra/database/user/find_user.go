package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ivandersr/go-auction/config/logger"
	"github.com/ivandersr/go-auction/internal/entity/user"
	"github.com/ivandersr/go-auction/internal/internal_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (r *UserRepository) FindUserById(ctx context.Context, id string) (*user.User, *internal_errors.InternalError) {
	var userMongo UserEntityMongo
	var errorMessage string
	filter := bson.M{"_id": id}
	if err := r.Collection.FindOne(ctx, filter).Decode(&userMongo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			errorMessage = fmt.Sprintf("user not found with id %s", id)
			logger.Error(errorMessage, err)
			return nil, internal_errors.NewNotFoundError(errorMessage)
		}
		errorMessage = fmt.Sprintf("error when trying to find user with id %s", id)
		logger.Error(errorMessage, err)
		return nil, internal_errors.NewInternalServerError(errorMessage)
	}

	return &user.User{
		Id:   userMongo.Id,
		Name: userMongo.Name,
	}, nil
}
