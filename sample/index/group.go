package index

import (
	"context"
	"encoding/xml"
	"strconv"

	"github.com/ryutah/gae-sample/sample/model"

	"google.golang.org/appengine/search"
)

type groupIndex struct {
	Name    string
	Belongs search.HTML
	Users   search.HTML
}

func NewGroupFromModel(g *model.Group, belongs []*model.Group, users []*model.User) *Group {
	group := &Group{
		ID:   strconv.FormatInt(g.ID, 10),
		Name: g.Name,
	}
	gBelongs := make([]GroupBelong, 0, len(belongs))
	for _, b := range belongs {
		gb := GroupBelong{ID: strconv.FormatInt(b.ID, 10), Name: b.Name}
		gBelongs = append(gBelongs, gb)
	}
	gUsers := make([]GroupUser, 0, len(users))
	for _, u := range users {
		gu := GroupUser{ID: strconv.FormatInt(u.ID, 10), Name: u.Name}
		gUsers = append(gUsers, gu)
	}
	group.Belongs, group.Users = gBelongs, gUsers
	return group
}

type (
	Group struct {
		ID      string        `xml:"id"`
		Name    string        `xml:"name"`
		Belongs []GroupBelong `xml:"belong"`
		Users   []GroupUser   `xml:"users"`
	}

	GroupBelong struct {
		ID   string `xml:"id"`
		Name string `xml:"name"`
	}

	GroupUser struct {
		ID   string `xml:"id"`
		Name string `xml:"name"`
	}
)

func (g *Group) Encode() (*groupIndex, error) {
	belongs, err := xml.Marshal(g.Belongs)
	if err != nil {
		return nil, err
	}
	users, err := xml.Marshal(g.Users)
	if err != nil {
		return nil, err
	}
	return &groupIndex{
		Name:    g.Name,
		Belongs: search.HTML(belongs),
		Users:   search.HTML(users),
	}, nil
}

func (g *Group) Decode(gi *groupIndex) error {
	g.Name = gi.Name
	if len(gi.Belongs) != 0 {
		if err := xml.Unmarshal([]byte(gi.Belongs), &g.Belongs); err != nil {
			return err
		}
	}
	if len(gi.Users) != 0 {
		if err := xml.Unmarshal([]byte(gi.Users), &g.Users); err != nil {
			return err
		}
	}
	return nil
}

func PutGroup(ctx context.Context, groupID int64, g *Group) (*Group, error) {
	index, err := search.Open("Group")
	if err != nil {
		return nil, err
	}
	id := strconv.FormatInt(groupID, 10)
	groupIdx, err := g.Encode()
	if _, err := index.Put(ctx, id, groupIdx); err != nil {
		return nil, err
	}
	return g, nil
}

func SearchGroup(ctx context.Context, q string) ([]*Group, error) {
	index, err := search.Open("Group")
	if err != nil {
		return nil, err
	}

	ite := index.Search(ctx, q, &search.SearchOptions{
		Limit: 50,
	})

	var groups []*Group
	for {
		gIdx := new(groupIndex)
		id, err := ite.Next(gIdx)
		if err == search.Done {
			break
		} else if err != nil {
			return nil, err
		}

		g := new(Group)
		if err := g.Decode(gIdx); err != nil {
			return nil, err
		}
		g.ID = id
		groups = append(groups, g)
	}
	return groups, nil
}
