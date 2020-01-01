package repository

import (
	ue "be/domain/user/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type Repository interface {
	InsertUser(ctx context.Context, user *ue.User) error
	GetUsers(ctx context.Context, filter bson.M) ([]*ue.User, error)
	GetUser(ctx context.Context, username string) (*ue.User, error)
}
