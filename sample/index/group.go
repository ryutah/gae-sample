package index

import (
	"context"
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/ryutah/gae-sample/sample/model"

	"google.golang.org/appengine/search"
)

type groupIndex struct {
	Search search.HTML
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

func PutGroup(ctx context.Context, groupID int64, g *Group) (*Group, error) {
	index, err := search.Open("Group")
	if err != nil {
		return nil, err
	}
	id := strconv.FormatInt(groupID, 10)
	groupXml, err := xml.Marshal(g)
	if err != nil {
		return nil, err
	}
	groupIdx := &groupIndex{Search: search.HTML(groupXml)}
	if _, err := index.Put(ctx, id, groupIdx); err != nil {
		return nil, err
	}
	return g, nil
}

func SearchGroup(ctx context.Context, params []string) ([]*Group, error) {
	index, err := search.Open("Group")
	if err != nil {
		return nil, err
	}

	q := strings.Join(params, " AND ")
	ite := index.Search(ctx, q, &search.SearchOptions{
		Limit: 50,
	})

	var groups []*Group
	for {
		gIdx := new(groupIndex)
		if _, err := ite.Next(gIdx); err == search.Done {
			break
		} else if err != nil {
			return nil, err
		}

		g := new(Group)
		if err := xml.Unmarshal([]byte(gIdx.Search), g); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
