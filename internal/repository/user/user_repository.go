package user

import (
	"backend/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) UserExist(email string) (bool, error) {
	err := u.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(models.User{})
	if err != nil {
		return false, err
	}

	return true, nil
}
