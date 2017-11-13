package sample

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ryutah/gae-sample/sample/index"
	"github.com/ryutah/gae-sample/sample/model"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type respGetUserList struct {
	Users []*model.User `json:"users"`
	Count int           `json:"count"`
}

type (
	UserPost    struct{}
	UserGetList struct{}
)

func (u *UserPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	payload := new(model.User)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Errorf(ctx, "failed to parse body : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	key := model.NewUserKey(ctx)
	newUsr, err := model.PutUser(ctx, key, payload)
	if err != nil {
		log.Errorf(ctx, "failed to post user : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := addUpdateIndexTask(ctx); err != nil {
		log.Errorf(ctx, "failed to add task : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(newUsr); err != nil {
		log.Errorf(ctx, "failed to write response : %v", err)
		http.Error(w, err.Error(), 500)
	}
}

func (u *UserGetList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	q := r.FormValue("q")
	uIndexs, err := index.SearchUser(ctx, q)
	if err != nil {
		log.Errorf(ctx, "failed to search : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	keys := make([]*datastore.Key, 0, len(uIndexs))
	for _, u := range uIndexs {
		id, _ := strconv.ParseInt(u.ID, 10, 64)
		log.Infof(ctx, "User id : %v", id)
		keys = append(keys, model.UserKey(ctx, id))
	}
	users, err := model.GetMultiUser(ctx, keys)
	if err != nil {
		log.Errorf(ctx, "failed to get users : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	resp := respGetUserList{Users: users, Count: len(users)}
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "failed to write response : %v", err)
		http.Error(w, err.Error(), 500)
	}
}
