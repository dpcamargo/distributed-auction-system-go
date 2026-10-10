package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/dpcamargo/distributed-auction-system-go/configuration/logger"
	"github.com/dpcamargo/distributed-auction-system-go/internal/entity/user_entity"
	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserEntityMongo struct {
	ID   string `bson:"_id"`
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

func (u *UserRepository) FindUserById(ctx context.Context, userID string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": userID}

	var userEntityMongo UserEntityMongo
	err := u.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error("User not found", err, zap.String("id", userID))
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("User not found with this id = %s", userID))
		}
		logger.Error("Error trying to find user by userID", err, zap.String("id", userID))
		return nil, internal_error.NewInternalServerError(
			fmt.Sprintf("Error trying to find user by userID = %s", userID))
	}
	userEntity := &user_entity.User{
		ID:   userEntityMongo.ID,
		Name: userEntityMongo.Name,
	}

	return userEntity, nil
}
