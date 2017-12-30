package main

import (
	graphql "github.com/neelance/graphql-go"
)

// Resolver struct
type Resolver struct {
	db *postgresDb
}

//     REALM
// =============

// RealmResolver struct
type RealmResolver struct {
	id    graphql.ID
	name  string
	title *string
	db    *postgresDb
}

// ID getter
func (r *RealmResolver) ID() graphql.ID {
	return r.id
}

// Name getter
func (r *RealmResolver) Name() string {
	return r.name
}

// Title getter
func (r *RealmResolver) Title() *string {
	return r.title
}

//     PLAYER
// ==============

// PlayerResolver struct
type PlayerResolver struct {
	id      graphql.ID
	realmID graphql.ID
	name    string
}

// ID getter
func (p *PlayerResolver) ID() graphql.ID {
	return p.id
}

// Name getter
func (p *PlayerResolver) Name() string {
	return p.name
}

// RealmID getter
func (p *PlayerResolver) RealmID() graphql.ID {
	return p.realmID
}

//     SESSION
// ===============

// SessionResolver struct
type SessionResolver struct {
	id      graphql.ID
	realmID graphql.ID
	name    *string
	time    string
}

// ID getter
func (s *SessionResolver) ID() graphql.ID {
	return s.id
}

// Name getter
func (s *SessionResolver) Name() *string {
	return s.name

}

// RealmID getter
func (s *SessionResolver) RealmID() graphql.ID {
	return s.realmID
}

// Time getter
func (s *SessionResolver) Time() string {
	return s.time
}

//     PLAYER_SESSION
// ======================

// PlayerSessionResolver struct
type PlayerSessionResolver struct {
	playerID  graphql.ID
	sessionID graphql.ID
	buyIn     int
	walkout   int
}

// PlayerID getter
func (ps *PlayerSessionResolver) PlayerID() graphql.ID {
	return ps.playerID
}

// SessionID getter
func (ps *PlayerSessionResolver) SessionID() graphql.ID {
	return ps.sessionID
}

// BuyIn getter
func (ps *PlayerSessionResolver) BuyIn() int {
	return ps.buyIn
}

// WalkOut getter
func (ps *PlayerSessionResolver) WalkOut() int {
	return ps.walkout
}
