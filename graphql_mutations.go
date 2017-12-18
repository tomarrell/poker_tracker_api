package main

import (
	"errors"
	"strconv"

	graphql "github.com/neelance/graphql-go"
	"gopkg.in/guregu/null.v3"
)

// CreateRealm args struct
type CreateRealm struct {
	Name  string
	Title *string
}

// CreateRealm mutation
func (r *Resolver) CreateRealm(args CreateRealm) (*RealmResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply realm name")
	}

	realm, err := r.db.CreateRealm(null.StringFrom(args.Name), null.StringFromPtr(args.Title))
	if err != nil {
		return nil, err
	}

	return &RealmResolver{&GQLRealm{
		ID:    graphql.ID(strconv.Itoa(realm.ID)),
		Name:  realm.Name,
		Title: realm.Title.Ptr(),
	}}, nil
}

// CreatePlayer args struct
type CreatePlayer struct {
	Name    string
	RealmID int32
}

// CreatePlayer mutation
func (r *Resolver) CreatePlayer(args CreatePlayer) (*PlayerResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply player name")
	}
	if args.RealmID == 0 {
		return nil, errors.New("Must supply Realm ID")
	}

	// TODO: Check if realm exists first to provide nicer error?
	player, err := r.db.CreatePlayer(null.StringFrom(args.Name), null.IntFrom(int64(args.RealmID)))

	if err != nil {
		return nil, err
	}

	return &PlayerResolver{&GQLPlayer{
		ID:      graphql.ID(strconv.Itoa(player.ID)),
		Name:    player.Name,
		RealmID: graphql.ID(strconv.Itoa(player.RealmID)),
	}}, nil
}
