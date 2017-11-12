package sample

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"

	"github.com/ryutah/gae-sample/sample/index"
	"github.com/ryutah/gae-sample/sample/model"
)

type UpdateIndex struct{}

func (u *UpdateIndex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Start update index")
	if err := updateUserIndexs(ctx); err != nil {
		log.Errorf(ctx, "failed to update index : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if err := updateGroupIndexs(ctx); err != nil {
		log.Errorf(ctx, "failed to update index : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func updateUserIndexs(ctx context.Context) error {
	keys, err := model.GetAllUserKeys(ctx)
	if err != nil {
		return err
	}
	users, err := model.GetMultiUser(ctx, keys)
	if err != nil {
		return err
	}

	for _, usr := range users {
		if err := updateUserIndex(ctx, usr); err != nil {
			return err
		}
	}
	return nil
}

func updateUserIndex(ctx context.Context, u *model.User) error {
	members, err := model.GetMultiMeberChildOf(ctx, u.Key)
	if err != nil {
		return err
	}

	gKeys := make([]*datastore.Key, 0, len(members))
	for _, member := range members {
		k := model.GroupKey(ctx, member.BelongTo)
		gKeys = append(gKeys, k)
	}
	groups, err := model.GetMultiGroup(ctx, gKeys)
	if err != nil {
		return err
	}

	idx := index.NewUserFromModel(u, groups...)
	if _, err := index.PutUser(ctx, u.ID, idx); err != nil {
		return err
	}
	return nil
}

func updateGroupIndexs(ctx context.Context) error {
	keys, err := model.GetAllGroupKeys(ctx)
	if err != nil {
		return err
	}
	groups, err := model.GetMultiGroup(ctx, keys)

	for _, g := range groups {
		if err := updateGroupIndex(ctx, g); err != nil {
			return err
		}
	}
	return nil
}

func updateGroupIndex(ctx context.Context, g *model.Group) error {
	belongs, err := model.GetMultiMeberChildOf(ctx, g.Key)
	if err != nil {
		return err
	}
	belongKeys := make([]*datastore.Key, 0, len(belongs))
	for _, b := range belongs {
		belongKeys = append(belongKeys, model.GroupKey(ctx, b.BelongTo))
	}
	belongGroups, err := model.GetMultiGroup(ctx, belongKeys)
	if err != nil {
		return err
	}

	belonged, err := model.GetMultiMemberBlongTo(ctx, g.ID)
	if err != nil {
		return err
	}
	userKeys := make([]*datastore.Key, 0, len(belonged))
	for _, b := range belonged {
		if b.Key.Parent().Kind() != "User" {
			continue
		}
		userKeys = append(userKeys, b.Key.Parent())
	}
	users, err := model.GetMultiUser(ctx, userKeys)

	gIdx := index.NewGroupFromModel(g, belongGroups, users)
	if _, err := index.PutGroup(ctx, g.ID, gIdx); err != nil {
		return err
	}
	return nil
}

func addUpdateIndexTask(ctx context.Context) error {
	t := taskqueue.Task{
		Path:   "/backend/index",
		Delay:  500 * time.Millisecond,
		Method: "POST",
	}
	_, err := taskqueue.Add(ctx, &t, "default")
	return err
}
