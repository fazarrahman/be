package radix

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	ue "be/domain/user/entity"
	user_mongo_repo "be/domain/user/repository"

	"github.com/mediocregopher/radix/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Radix struct {
	rdx            *radix.Pool
	usr_mongo_repo user_mongo_repo.Repository
}

func New(rdx *radix.Pool, _usr user_mongo_repo.Repository) *Radix {
	return &Radix{rdx: rdx, usr_mongo_repo: _usr}
}

func (r *Radix) InsertUser(ctx context.Context, user *ue.User) error {
	return r.usr_mongo_repo.InsertUser(ctx, user)
}

func (r *Radix) GetUsers(ctx context.Context, filter bson.M) ([]*ue.User, error) {
	return r.usr_mongo_repo.GetUsers(ctx, filter)
}

func (r *Radix) GetUser(ctx context.Context, username string) (*ue.User, error) {
	return r.usr_mongo_repo.GetUser(ctx, username)
}

func (r *Radix) setUserRadix(user *ue.User) error {
	jUser, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		return err
	}

	err = r.rdx.Do(radix.FlatCmd(nil, "SET", (fmt.Sprintf("user:%d", user.Username)), jUser))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
