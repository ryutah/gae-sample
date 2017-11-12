package sample

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	gcxt "golang.org/x/net/context"

	"github.com/gorilla/mux"

	"github.com/ryutah/gae-sample/sample/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type respMemberGetList struct {
	Groups []*model.Group `json:"groups"`
	Users  []*model.User  `json:"users"`
}

type (
	MemberPost    struct{}
	MemberGetList struct{}
)

func (m *MemberPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	params := mux.Vars(r)
	var (
		groupID, _ = strconv.ParseInt(params["groupId"], 10, 64)
		typ        = params["type"]
		id, _      = strconv.ParseInt(params["id"], 10, 64)
	)

	if err := addGroupMember(ctx, typ, groupID, id); err != nil {
		log.Errorf(ctx, "failed to put member : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := addUpdateIndexTask(ctx); err != nil {
		log.Errorf(ctx, "failed to add task : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(201)
}

func (m *MemberGetList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	groupID, _ := strconv.ParseInt(mux.Vars(r)["groupId"], 10, 64)
	groups, users, err := getMembers(ctx, groupID)
	if err != nil {
		log.Errorf(ctx, "failed to get members : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	resp := &respMemberGetList{
		Groups: groups,
		Users:  users,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "failed to write body : %v", err)
		http.Error(w, err.Error(), 500)
	}
}

func addGroupMember(ctx context.Context, typ string, groupID, id int64) error {
	belongKey := model.GroupKey(ctx, groupID)
	belong, err := model.GetGroup(ctx, belongKey)
	if err != nil {
		return err
	}

	var key *datastore.Key

	switch typ {
	case "group":
		key, err = memberKeyForGroups(ctx, belong, id)
	case "user":
		key, err = memberKeyForUser(ctx, belong, id)
	default:
		err = fmt.Errorf("unkonw member type %v", typ)
	}
	if err != nil {
		return err
	}

	m := &model.Member{BelongTo: belong.ID}
	if _, err := model.PutMember(ctx, key, m); err != nil {
		return err
	}
	return nil
}

func runInTransaction(ctx context.Context, f func(tc context.Context) error) error {
	return datastore.RunInTransaction(ctx, func(tc gcxt.Context) error {
		return f(context.Context(tc))
	}, nil)
}

func memberKeyForUser(ctx context.Context, belong *model.Group, id int64) (*datastore.Key, error) {
	usrKey := model.UserKey(ctx, id)
	if _, err := model.GetUser(ctx, usrKey); err != nil {
		return nil, err
	}

	return model.NewMemberKey(ctx, usrKey), nil
}

func memberKeyForGroups(ctx context.Context, belong *model.Group, id int64) (*datastore.Key, error) {
	gKey := model.GroupKey(ctx, id)
	if _, err := model.GetGroup(ctx, gKey); err != nil {
		return nil, err
	}

	return model.NewMemberKey(ctx, gKey), nil
}

func getMembers(ctx context.Context, belong int64) ([]*model.Group, []*model.User, error) {
	members, err := model.GetMultiMemberBlongTo(ctx, belong)
	if err != nil {
		return nil, nil, err
	}
	var gKeys, uKeys []*datastore.Key
	for _, m := range members {
		if m.Key.Parent().Kind() == "Group" {
			gKeys = append(gKeys, m.Key.Parent())
		} else {
			uKeys = append(uKeys, m.Key.Parent())
		}
	}
	groups, err := model.GetMultiGroup(ctx, gKeys)
	if err != nil {
		return nil, nil, err
	}
	users, err := model.GetMultiUser(ctx, uKeys)
	if err != nil {
		return nil, nil, err
	}
	return groups, users, nil
}
