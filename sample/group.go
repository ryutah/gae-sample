package sample

import (
	"encoding/json"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/ryutah/gae-sample/sample/index"
	"github.com/ryutah/gae-sample/sample/model"
)

type respGetGroupList struct {
	Groups []*model.Group `json:"groups"`
	Count  int            `json:"count"`
}

type (
	GroupPost    struct{}
	GroupGetList struct{}
)

func (g *GroupPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	payload := new(model.Group)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Errorf(ctx, "failed to parse body : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	key := model.NewGroupKey(ctx)
	newGroup, err := model.PutGroup(ctx, key, payload)
	if err != nil {
		log.Errorf(ctx, "failed to post group : %v", err)
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
	if err := json.NewEncoder(w).Encode(newGroup); err != nil {
		log.Errorf(ctx, "failed to write response : %v", err)
		http.Error(w, err.Error(), 500)
	}
}

func (g *GroupGetList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	q := r.FormValue("q")
	gIndexs, err := index.SearchGroup(ctx, q)
	if err != nil {
		log.Errorf(ctx, "failed to search : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	keys := make([]*datastore.Key, 0, len(gIndexs))
	for _, g := range gIndexs {
		id, _ := strconv.ParseInt(g.ID, 10, 64)
		log.Infof(ctx, "Group id : %v", id)
		keys = append(keys, model.GroupKey(ctx, id))
	}
	groups, err := model.GetMultiGroup(ctx, keys)
	if err != nil {
		log.Errorf(ctx, "failed to get groups : %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	resp := respGetGroupList{Groups: groups, Count: len(groups)}
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Errorf(ctx, "failed to write response : %v", err)
		http.Error(w, err.Error(), 500)
	}
}
