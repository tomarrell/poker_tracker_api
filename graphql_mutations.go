package main

import (
	"errors"

	graphql "github.com/neelance/graphql-go"
	"gopkg.in/guregu/null.v3"
)

// CreateRealm args struct
type CreateRealm struct {
	Name  string
	title *string
}

// CreateRealm mutation
func (r *Resolver) CreateRealm(args CreateRealm) (*RealmResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply realm name")
	}

	realm, err := r.db.CreateRealm(null.StringFrom(args.Name), null.StringFromPtr(args.title))
	if err != nil {
		return nil, err
	}

	return &RealmResolver{&GQLRealm{
		ID:    graphql.ID(realm.ID),
		Name:  realm.Name,
		Title: realm.Title.Ptr(),
	}}, nil
}
