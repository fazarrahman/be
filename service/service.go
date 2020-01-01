package service

import (
	"context"
	ue "be/domain/user/entity"
	user "be/domain/user/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetUser(ctx context.Context, username string) (*ue.User, error)
	GetUsers(ctx context.Context, filter bson.M) ([]*ue.User, error)
	InsertUser(ctx context.Context, r *UserRequest) error
	CheckUsernamePassword(ctx context.Context, r *UserPasswordCheckRequest) (bool, error)
}

type svc struct {
	UserRepository user.Repository
}

func New(_usrRepo user.Repository) Service {
	return &svc{
		UserRepository: _usrRepo,
	}
}
