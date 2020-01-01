package service

import (
	"context"
	ue "be/domain/user/entity"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Username  string `validate:"required,min=1"`
	FirstName string `validate:"required,min=1"`
	LastName  string `validate:"required,min=1"`
	Password  string `validate:"required,min=1"`
	Status    int8
}

type UserPasswordCheckRequest struct {
	Username string `validate:"required,min=1"`
	Password string `validate:"required,min=1"`
}

func (s *svc) GetUser(ctx context.Context, username string) (*ue.User, error) {
	user, err := s.UserRepository.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *svc) GetUsers(ctx context.Context, filter bson.M) ([]*ue.User, error) {
	users, err := s.UserRepository.GetUsers(ctx, filter)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *svc) InsertUser(ctx context.Context, r *UserRequest) error {
	pwdhash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.UserRepository.InsertUser(ctx, &ue.User{
		Username:  r.Username,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Password:  pwdhash,
		Status:    1,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *svc) CheckUsernamePassword(ctx context.Context, r *UserPasswordCheckRequest) (bool, error) {
	user, err := s.GetUser(ctx, r.Username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(r.Password))
	if err != nil {
		return false, err
	}
	return true, nil
}
