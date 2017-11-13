package index

import (
	"context"
	"encoding/xml"
	"strconv"

	"github.com/ryutah/gae-sample/sample/model"

	"google.golang.org/appengine/search"
)

type userIndex struct {
	Name    string
	Email   string
	Belongs search.HTML
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

func (u *User) encode() (*userIndex, error) {
	belongs, err := xml.Marshal(u.Belongs)
	if err != nil {
		return nil, err
	}
	return &userIndex{
		Name:    u.Name,
		Email:   u.Email,
		Belongs: search.HTML(belongs),
	}, nil
}

func (u *User) decode(ui *userIndex) error {
	u.Name, u.Email = ui.Name, ui.Email
	if len(ui.Belongs) != 0 {
		if err := xml.Unmarshal([]byte(ui.Belongs), &u.Belongs); err != nil {
			return err
		}
	}
	return nil
}

func PutUser(ctx context.Context, userID int64, u *User) (*User, error) {
	index, err := search.Open("User")
	if err != nil {
		return nil, err
	}
	id := strconv.FormatInt(userID, 10)
	userIdx, err := u.encode()
	if _, err := index.Put(ctx, id, userIdx); err != nil {
		return nil, err
	}
	return u, nil
}

func SearchUser(ctx context.Context, q string) ([]*User, error) {
	index, err := search.Open("User")
	if err != nil {
		return nil, err
	}

	ite := index.Search(ctx, q, &search.SearchOptions{Limit: 50})
	var users []*User
	for {
		uIdx := new(userIndex)
		id, err := ite.Next(uIdx)
		if err == search.Done {
			break
		} else if err != nil {
			return nil, err
		}

		u := new(User)
		if err := u.decode(uIdx); err != nil {
			return nil, err
		}
		u.ID = id
		users = append(users, u)
	}
	return users, nil
}
