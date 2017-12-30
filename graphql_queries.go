package main

import (
	"errors"
	"strconv"
	"time"

	graphql "github.com/neelance/graphql-go"
)

// RealmByName resolver
func (r *Resolver) RealmByName(args struct{ Name string }) (*RealmResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply realm name")
	}

	realm, err := r.db.GetRealmByField("name", args.Name)
	if err != nil {
		return nil, err
	}

	return &RealmResolver{
		id:    graphql.ID(strconv.Itoa(realm.ID)),
		name:  realm.Name,
		title: realm.Title.Ptr(),
		db:    r.db}, nil
}

func (r *Resolver) RealmByID(args struct{ ID graphql.ID }) (*RealmResolver, error) {
	if args.ID == "" {
		return nil, errors.New("Must supply realm ID")
	}

	realm, err := r.db.GetRealmByField("id", args.ID)
	if err != nil {
		return nil, err
	}

	return &RealmResolver{
		id:    graphql.ID(strconv.Itoa(realm.ID)),
		name:  realm.Name,
		title: realm.Title.Ptr(),
		db:    r.db}, nil
}

func (r *RealmResolver) Players() ([]*PlayerResolver, error) {
	id, err := strconv.Atoi(string(r.ID()))
	if err != nil {
		return nil, errors.New("realm id must be numerical")
	}

	players, err := r.db.GetPlayersByRealmID(id)
	if err != nil {
		return nil, err
	}

	var pr []*PlayerResolver

	for _, p := range players {
		pr = append(pr,
			&PlayerResolver{
				id:      graphql.ID(strconv.Itoa(p.ID)),
				realmID: graphql.ID(strconv.Itoa(p.RealmID)),
				name:    p.Name,
			},
		)
	}
	return pr, nil

}

func (r *RealmResolver) Sessions() ([]*SessionResolver, error) {
	id, err := strconv.Atoi(string(r.ID()))
	if err != nil {
		return nil, errors.New("realm id must be numerical")
	}

	sessions, err := r.db.GetSessionsByRealmID(id)
	if err != nil {
		return nil, err
	}

	var sr []*SessionResolver

	for _, s := range sessions {
		sr = append(sr,
			&SessionResolver{
				id:      graphql.ID(strconv.Itoa(s.ID)),
				realmID: graphql.ID(strconv.Itoa(s.RealmID)),
				name:    s.Name.Ptr(),
				time:    s.Time.UTC().Format(time.RFC3339),
			},
		)
	}
	return sr, nil

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

	return &PlayerResolver{
		id:      graphql.ID(strconv.Itoa(player.ID)),
		realmID: graphql.ID(strconv.Itoa(player.RealmID)),
		name:    player.Name,
	}, nil
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

	return &SessionResolver{
		id:      graphql.ID(strconv.Itoa(session.ID)),
		realmID: graphql.ID(strconv.Itoa(session.RealmID)),
		name:    session.Name.Ptr(),
		time:    session.Time.UTC().Format(time.RFC3339),
	}, nil
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

	sessions, err := r.db.GetSessionsByRealmID(id)
	if err != nil {
		return nil, err
	}

	var sr []*SessionResolver

	for _, s := range sessions {
		sr = append(sr,
			&SessionResolver{
				id:      graphql.ID(strconv.Itoa(s.ID)),
				realmID: graphql.ID(strconv.Itoa(s.RealmID)),
				name:    s.Name.Ptr(),
				time:    s.Time.UTC().Format(time.RFC3339),
			},
		)
	}

	return &sr, nil
}
