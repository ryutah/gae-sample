package model

import (
	"context"

	"google.golang.org/appengine/datastore"
)

// Member member
type Member struct {
	Key      *datastore.Key `datastore:"-"`
	BelongTo int64
}

func PutMember(ctx context.Context, key *datastore.Key, m *Member) (*Member, error) {
	newMem := *m
	if _, err := datastore.Put(ctx, key, m); err != nil {
		return nil, err
	}
	newMem.Key = key
	return &newMem, nil
}

func GetMultiMeberChildOf(ctx context.Context, parent *datastore.Key) ([]*Member, error) {
	var members []*Member
	keys, err := datastore.NewQuery("Member").Ancestor(parent).GetAll(ctx, &members)
	if err != nil {
		return nil, err
	}
	for i, m := range members {
		m.Key = keys[i]
	}
	return members, nil
}

func GetMultiMemberBlongTo(ctx context.Context, belong int64) ([]*Member, error) {
	var members []*Member
	keys, err := datastore.NewQuery("Member").Filter("BelongTo=", belong).GetAll(ctx, &members)
	if err != nil {
		return nil, err
	}
	for i, m := range members {
		m.Key = keys[i]
	}
	return members, nil
}

func NewMemberKey(ctx context.Context, parent *datastore.Key) *datastore.Key {
	return datastore.NewIncompleteKey(ctx, "Member", parent)
}
