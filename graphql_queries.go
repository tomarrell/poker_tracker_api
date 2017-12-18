package main

import (
	"errors"
	"fmt"
	"strconv"

	graphql "github.com/neelance/graphql-go"
)

// RealmByName resolver
func (r *Resolver) RealmByName(args struct{ Name string }) (*RealmResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply realm name")
	}

	realm, err := r.db.GetRealmByName(args.Name)
	if err != nil {
		return nil, err
	}

	return &RealmResolver{&GQLRealm{
		ID:    graphql.ID(strconv.Itoa(realm.ID)),
		Name:  realm.Name,
		Title: realm.Title.Ptr(),
	}}, nil
}

// PlayerByID resolver
func (r *Resolver) PlayerByID(args struct{ ID graphql.ID }) (*PlayerResolver, error) {
	if args.ID == "" {
		return nil, errors.New("Must supply player id")
	}

	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, errors.New("player id must be numerical")
	}

	player, err := r.db.GetPlayerByID(id)
	if err != nil {
		return nil, err
	}

	return &PlayerResolver{&GQLPlayer{
		ID:      graphql.ID(strconv.Itoa(player.ID)),
		RealmID: graphql.ID(strconv.Itoa(player.RealmID)),
		Name:    player.Name,
	}}, nil
}

// SessionByID resolver
func (r *Resolver) SessionByID(args struct{ ID graphql.ID }) (*SessionResolver, error) {
	if args.ID == "" {
		return nil, errors.New("Must supply session id")
	}

	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, errors.New("Session ID must be numerical")
	}

	session, err := r.db.GetSessionByID(id)
	if err != nil {
		return nil, err
	}

	return &SessionResolver{&GQLSession{
		ID:      graphql.ID(strconv.Itoa(session.ID)),
		RealmID: graphql.ID(strconv.Itoa(session.RealmID)),
		Name:    session.Name.Ptr(),
		Time:    fmt.Sprintf("%s", session.Time),
	}}, nil
}

// SessionsByRealmID resolver
func (r *Resolver) SessionsByRealmID(args struct{ RealmID graphql.ID }) (*[]*SessionResolver, error) {
	if args.RealmID == "" {
		return nil, errors.New("Must supply RealmID")
	}

	id, err := strconv.Atoi(string(args.RealmID))
	if err != nil {
		return nil, errors.New("RealmID must be numerical")
	}

	sessions, err := r.db.GetSessions(id)
	if err != nil {
		return nil, err
	}

	var sr []*SessionResolver

	for _, s := range sessions {
		sr = append(sr,
			&SessionResolver{&GQLSession{
				ID:      graphql.ID(strconv.Itoa(s.ID)),
				RealmID: graphql.ID(strconv.Itoa(s.RealmID)),
				Name:    s.Name.Ptr(),
				Time:    fmt.Sprintf("%s", s.Time),
			}},
		)
	}

	return &sr, nil
}
