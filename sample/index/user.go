package index

import (
	"context"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/ryutah/gae-sample/sample/model"

	"google.golang.org/appengine/search"
)

type userIndex struct {
	Search search.HTML
}

func NewUserFromModel(u *model.User, belongs ...*model.Group) *User {
	usr := &User{
		ID:    strconv.FormatInt(u.ID, 10),
		Name:  u.Name,
		Email: u.Email,
	}
	groups := make([]UserBelong, 0, len(belongs))
	for _, b := range belongs {
		g := UserBelong{
			ID:   strconv.FormatInt(b.ID, 10),
			Name: b.Name,
		}
		groups = append(groups, g)
	}
	usr.Belongs = groups
	return usr
}

type (
	User struct {
		ID      string       `xml:"id"`
		Name    string       `xml:"name"`
		Email   string       `xml:"email"`
		Belongs []UserBelong `xml:"belongs"`
	}

	UserBelong struct {
		ID   string `xml:"id"`
		Name string `xml:name`
	}
)

func PutUser(ctx context.Context, userID int64, u *User) (*User, error) {
	index, err := search.Open("User")
	if err != nil {
		return nil, err
	}
	id := strconv.FormatInt(userID, 10)
	userXml, err := xml.Marshal(u)
	if err != nil {
		return nil, err
	}
	userIdx := &userIndex{Search: search.HTML(userXml)}
	if _, err := index.Put(ctx, id, userIdx); err != nil {
		return nil, err
	}
	return u, nil
}

func SearchUser(ctx context.Context, params []string) ([]*User, error) {
	index, err := search.Open("User")
	if err != nil {
		return nil, err
	}

	q := strings.Join(params, " AND ")
	ite := index.Search(ctx, q, &search.SearchOptions{
		Limit: 50,
	})

	var users []*User
	for {
		uIdx := new(userIndex)
		if _, err := ite.Next(uIdx); err == search.Done {
			break
		} else if err != nil {
			return nil, err
		}

		u := new(User)
		if err := xml.Unmarshal([]byte(uIdx.Search), u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
