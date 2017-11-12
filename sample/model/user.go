package model

import (
	"context"

	"google.golang.org/appengine/datastore"
)

// User user
type User struct {
	Key   *datastore.Key `datastore:"-" json:"-"`
	ID    int64          `datastore:"-" json:"id"`
	Name  string         `json:"name"`
	Email string         `json:"email"`
}

func GetAllUserKeys(ctx context.Context) ([]*datastore.Key, error) {
	return datastore.NewQuery("User").KeysOnly().GetAll(ctx, nil)
}

func PutUser(ctx context.Context, key *datastore.Key, u *User) (*User, error) {
	newKey, err := datastore.Put(ctx, key, u)
	if err != nil {
		return nil, err
	}
	newUsr := *u
	newUsr.Key, newUsr.ID = newKey, newKey.IntID()
	return &newUsr, nil
}

func DeleteUser(ctx context.Context, key *datastore.Key) error {
	return datastore.Delete(ctx, key)
}

func GetUser(ctx context.Context, key *datastore.Key) (*User, error) {
	u := new(User)
	if err := datastore.Get(ctx, key, u); err != nil {
		return nil, err
	}
	u.Key, u.ID = key, key.IntID()
	return u, nil
}

func GetMultiUser(ctx context.Context, keys []*datastore.Key) ([]*User, error) {
	users := make([]*User, len(keys))
	if err := datastore.GetMulti(ctx, keys, users); err != nil {
		return nil, err
	}
	for i, u := range users {
		u.Key, u.ID = keys[i], keys[i].IntID()
	}
	return users, nil
}

func NewUserKey(ctx context.Context) *datastore.Key {
	return datastore.NewIncompleteKey(ctx, "User", nil)
}

func UserKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "User", "", id, nil)
}
