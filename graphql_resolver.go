package main

import (
	"errors"
	"fmt"
	"strconv"

	graphql "github.com/neelance/graphql-go"
)

type GQLRealm struct {
	ID    graphql.ID
	Name  string
	Title *string
}
type GQLPlayer struct {
	ID      graphql.ID
	RealmID graphql.ID
	Name    string
}
type GQLSession struct {
	ID      graphql.ID
	RealmID graphql.ID
	Name    *string
	Time    string
}
type GQLPlayerSession struct {
	PlayerID  graphql.ID
	SessionID graphql.ID
	BuyIn     int
	Walkout   int
}
type Resolver struct {
	db *postgresDb
}
type realmResolver struct {
	r *GQLRealm
}

func (r *realmResolver) ID() graphql.ID {
	return r.r.ID
}
func (r *realmResolver) Name() string {
	return r.r.Name

}
func (r *realmResolver) Title() *string {
	return r.r.Title

}

type playerResolver struct {
	p *GQLPlayer
}

func (p *playerResolver) ID() graphql.ID {
	return p.p.ID
}
func (p *playerResolver) Name() string {
	return p.p.Name

}
func (p *playerResolver) RealmID() graphql.ID {
	return p.p.RealmID
}

type sessionResolver struct {
	s *GQLSession
}

func (s *sessionResolver) ID() graphql.ID {
	return s.s.ID
}
func (s *sessionResolver) Name() *string {
	return s.s.Name

}
func (s *sessionResolver) RealmID() graphql.ID {
	return s.s.RealmID
}
func (s *sessionResolver) Time() string {
	return s.s.Time
}

type playerSessionResolver struct {
	ps *GQLPlayerSession
}

func (ps *playerSessionResolver) PlayerID() graphql.ID {
	return ps.ps.PlayerID
}
func (ps *playerSessionResolver) SessionID() graphql.ID {
	return ps.ps.SessionID
}
func (ps *playerSessionResolver) BuyIn() int {
	return ps.ps.BuyIn

}
func (ps *playerSessionResolver) WalkOut() int {
	return ps.ps.Walkout
}

func (r *Resolver) RealmByName(args struct{ Name string }) (*realmResolver, error) {
	if args.Name == "" {
		return nil, errors.New("Must supply realm name")
	}
	realm, err := r.db.GetRealmByName(args.Name)
	if err != nil {
		return nil, err
	}
	return &realmResolver{&GQLRealm{
		ID:    graphql.ID(realm.ID),
		Name:  realm.Name,
		Title: realm.Title.Ptr(),
	}}, nil
}

func (r *Resolver) PlayerById(args struct{ ID graphql.ID }) (*playerResolver, error) {
	if args.ID == "" {
		return nil, errors.New("Must supply player id")
	}
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, errors.New("player id must be numerical")
	}
	player, err := r.db.GetPlayerById(id)
	if err != nil {
		return nil, err
	}
	return &playerResolver{&GQLPlayer{
		ID:      graphql.ID(strconv.Itoa(player.ID)),
		RealmID: graphql.ID(strconv.Itoa(player.RealmID)),
		Name:    player.Name,
	}}, nil
}

func (r *Resolver) SessionById(args struct{ ID graphql.ID }) (*sessionResolver, error) {
	if args.ID == "" {
		return nil, errors.New("Must supply session id")
	}
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, errors.New("Session id must be numerical")
	}
	session, err := r.db.GetSessionById(id)
	if err != nil {
		return nil, err
	}
	return &sessionResolver{&GQLSession{
		ID:      graphql.ID(strconv.Itoa(session.ID)),
		RealmID: graphql.ID(strconv.Itoa(session.RealmID)),
		Name:    session.Name.Ptr(),
		Time:    fmt.Sprintf("%s", session.Time),
	}}, nil
}

func (r *Resolver) SessionsByRealmId(args struct{ RealmID graphql.ID }) (*[]*sessionResolver, error) {
	if args.RealmID == "" {
		return nil, errors.New("Must supply realm id")
	}
	id, err := strconv.Atoi(string(args.RealmID))
	if err != nil {
		return nil, errors.New("realm id must be numerical")
	}
	sessions, err := r.db.GetSessions(id)
	if err != nil {
		return nil, err
	}
	var sr []*sessionResolver
	for _, s := range sessions {
		sr = append(sr,
			&sessionResolver{&GQLSession{
				ID:      graphql.ID(strconv.Itoa(s.ID)),
				RealmID: graphql.ID(strconv.Itoa(s.RealmID)),
				Name:    s.Name.Ptr(),
				Time:    fmt.Sprintf("%s", s.Time),
			}},
		)
	}
	return &sr, nil
}
