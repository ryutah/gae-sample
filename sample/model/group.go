package model

import (
	"context"

	"google.golang.org/appengine/datastore"
)

// Group group
type Group struct {
	Key  *datastore.Key `datastore:"-" json:"-"`
	ID   int64          `datastore:"-" json:"id"`
	Name string         `json:"name"`
}

func GetAllGroupKeys(ctx context.Context) ([]*datastore.Key, error) {
	return datastore.NewQuery("Group").KeysOnly().GetAll(ctx, nil)
}

func PutGroup(ctx context.Context, key *datastore.Key, g *Group) (*Group, error) {
	newKey, err := datastore.Put(ctx, key, g)
	if err != nil {
		return nil, err
	}
	newG := *g
	newG.Key, newG.ID = newKey, newKey.IntID()
	return &newG, nil
}

func DeleteGroup(ctx context.Context, key *datastore.Key) error {
	return datastore.Delete(ctx, key)
}

func GetGroup(ctx context.Context, key *datastore.Key) (*Group, error) {
	g := new(Group)
	if err := datastore.Get(ctx, key, g); err != nil {
		return nil, err
	}
	g.Key, g.ID = key, key.IntID()
	return g, nil
}

func GetMultiGroup(ctx context.Context, keys []*datastore.Key) ([]*Group, error) {
	groups := make([]*Group, len(keys))
	if err := datastore.GetMulti(ctx, keys, groups); err != nil {
		return nil, err
	}
	for i, g := range groups {
		g.Key, g.ID = keys[i], keys[i].IntID()
	}
	return groups, nil
}

func NewGroupKey(ctx context.Context) *datastore.Key {
	return datastore.NewIncompleteKey(ctx, "Group", nil)
}

func GroupKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, "Group", "", id, nil)
}
