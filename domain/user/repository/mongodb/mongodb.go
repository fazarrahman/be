package mongodb

import (
	"context"
	ue "be/domain/user/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongodb struct {
	usersCollection *mongo.Collection
}

func New(_mdb *mongo.Database) *Mongodb {
	return &Mongodb{usersCollection: _mdb.Collection("users")}
}

func (m *Mongodb) InsertUser(ctx context.Context, user *ue.User) error {
	_, err := m.usersCollection.InsertOne(ctx, user)
	return err
}

func (m *Mongodb) GetUsers(ctx context.Context, filter bson.M) ([]*ue.User, error) {
	var users []*ue.User
	res, err := m.usersCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for res.Next(ctx) {
		var user ue.User
		err = res.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (m *Mongodb) GetUser(ctx context.Context, username string) (*ue.User, error) {
	res := m.usersCollection.FindOne(ctx, bson.M{"username": username})
	var user ue.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
