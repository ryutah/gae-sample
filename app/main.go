package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryutah/gae-sample/sample"
)

func init() {
	r := mux.NewRouter()

	r.Handle("/groups", new(sample.GroupPost)).Methods("POST")
	r.Handle("/groups", new(sample.GroupGetList)).Methods("GET")

	r.Handle("/users", new(sample.UserPost)).Methods("POST")
	r.Handle("/users", new(sample.UserGetList)).Methods("GET")

	r.Handle("/groups/{groupId:[0-9]+}/members/{type:(?:group)|(?:user)}/{id:[0-9]+}", new(sample.MemberPost)).Methods("POST")
	r.Handle("/groups/{groupId:[0-9]+}/members", new(sample.MemberGetList)).Methods("GET")

	r.Handle("/backend/index", new(sample.UpdateIndex)).Methods("POST")

	http.Handle("/", r)
}
