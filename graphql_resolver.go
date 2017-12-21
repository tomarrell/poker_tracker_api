package main

import (
	graphql "github.com/neelance/graphql-go"
)

// GQLRealm struct
type GQLRealm struct {
	ID    graphql.ID
	Name  string
	Title *string
}

// GQLPlayer struct
type GQLPlayer struct {
	ID      graphql.ID
	RealmID graphql.ID
	Name    string
}

// GQLSession struct
type GQLSession struct {
	ID      graphql.ID
	RealmID graphql.ID
	Name    *string
	Time    string
}

// GQLPlayerSession struct
type GQLPlayerSession struct {
	PlayerID  graphql.ID
	SessionID graphql.ID
	BuyIn     int
	Walkout   int
}

// Resolver struct
type Resolver struct {
	db *postgresDb
}

//     REALM
// =============

// RealmResolver struct
type RealmResolver struct {
	r *GQLRealm
}

// ID getter
func (r *RealmResolver) ID() graphql.ID {
	return r.r.ID
}

// Name getter
func (r *RealmResolver) Name() string {
	return r.r.Name
}

// Title getter
func (r *RealmResolver) Title() *string {
	return r.r.Title
}

//     PLAYER
// ==============

// PlayerResolver struct
type PlayerResolver struct {
	p *GQLPlayer
}

// ID getter
func (p *PlayerResolver) ID() graphql.ID {
	return p.p.ID
}

// Name getter
func (p *PlayerResolver) Name() string {
	return p.p.Name
}

// RealmID getter
func (p *PlayerResolver) RealmID() graphql.ID {
	return p.p.RealmID
}

//     SESSION
// ===============

// SessionResolver struct
type SessionResolver struct {
	s *GQLSession
}

// ID getter
func (s *SessionResolver) ID() graphql.ID {
	return s.s.ID
}

// Name getter
func (s *SessionResolver) Name() *string {
	return s.s.Name

}

// RealmID getter
func (s *SessionResolver) RealmID() graphql.ID {
	return s.s.RealmID
}

// Time getter
func (s *SessionResolver) Time() string {
	return s.s.Time
}

//     PLAYER_SESSION
// ======================

// PlayerSessionResolver struct
type PlayerSessionResolver struct {
	ps *GQLPlayerSession
}

// PlayerID getter
func (ps *PlayerSessionResolver) PlayerID() graphql.ID {
	return ps.ps.PlayerID
}

// SessionID getter
func (ps *PlayerSessionResolver) SessionID() graphql.ID {
	return ps.ps.SessionID
}

// BuyIn getter
func (ps *PlayerSessionResolver) BuyIn() int {
	return ps.ps.BuyIn
}

// WalkOut getter
func (ps *PlayerSessionResolver) WalkOut() int {
	return ps.ps.Walkout
}
