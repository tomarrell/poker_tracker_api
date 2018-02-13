package main

import (
	"errors"
	"strconv"
	"time"

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

// Players getter
func (r *RealmResolver) Players() ([]*PlayerResolver, error) {
	id, err := strconv.Atoi(string(r.ID()))
	if err != nil {
		return nil, errors.New("Realm id must be numerical")
	}

	players, err := r.db.GetPlayersByRealmID(id)
	if err != nil {
		return nil, err
	}

	var pr []*PlayerResolver

	for _, p := range players {
		pr = append(pr,
			&PlayerResolver{
				id:      toGQL(p.ID),
				realmID: graphql.ID(strconv.Itoa(p.RealmID)),
				name:    p.Name,
				db:      r.db,
			},
		)
	}

	return pr, nil
}

// Sessions getter
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
				id:      toGQL(s.ID),
				realmID: toGQL(s.RealmID),
				name:    s.Name.Ptr(),
				time:    s.Time.UTC().Format(time.RFC3339),
				db:      r.db,
			},
		)
	}

	return sr, nil
}

//     PLAYER
// ==============

// PlayerResolver struct
type PlayerResolver struct {
	id                graphql.ID
	realmID           graphql.ID
	name              string
	db                *postgresDb
	realBalance       int32
	historicalBalance int32
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

func (p *PlayerResolver) RealBalance() (int32, error) {
	id, err := strconv.Atoi(string(p.ID()))
	if err != nil {
		return 0, errors.New("player id must be numerical")
	}

	balance, err := p.db.GetRealBalanceByPlayerID(id)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (p *PlayerResolver) HistoricalBalance() (int32, error) {
	id, err := strconv.Atoi(string(p.ID()))
	if err != nil {
		return 0, errors.New("player id must be numerical")
	}

	balance, err := p.db.GetHistoricalBalanceByPlayerID(id)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// PlayerSessions getter
func (p *PlayerResolver) PlayerSessions() ([]*PlayerSessionResolver, error) {
	id, err := strconv.Atoi(string(p.ID()))
	if err != nil {
		return nil, errors.New("player id must be numerical")
	}

	pSessions, err := p.db.GetPlayerSessionsByField("player_id", id)
	if err != nil {
		return nil, err
	}

	sr := make([]*PlayerSessionResolver, len(pSessions))

	for i, s := range pSessions {
		sr[i] = &PlayerSessionResolver{
			playerID:  toGQL(s.PlayerID),
			sessionID: toGQL(s.SessionID),
			buyIn:     int32(s.Buyin.Int64),
			walkout:   int32(s.Walkout.Int64),
			db:        p.db,
		}
	}

	return sr, nil
}

//     SESSION
// ===============

// SessionResolver struct
type SessionResolver struct {
	id      graphql.ID
	realmID graphql.ID
	name    *string
	time    string
	db      *postgresDb
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

// PlayerSessions getter
func (s *SessionResolver) PlayerSessions() ([]*PlayerSessionResolver, error) {
	id, err := strconv.Atoi(string(s.ID()))
	if err != nil {
		return nil, errors.New("session id must be numerical")
	}

	pSessions, err := s.db.GetPlayerSessionsByField("session_id", id)
	if err != nil {
		return nil, err
	}

	sr := make([]*PlayerSessionResolver, len(pSessions))

	for i, ps := range pSessions {
		sr[i] = &PlayerSessionResolver{
			playerID:  toGQL(ps.PlayerID),
			sessionID: toGQL(ps.SessionID),
			buyIn:     int32(ps.Buyin.Int64),
			walkout:   int32(ps.Walkout.Int64),
			db:        s.db,
		}
	}

	return sr, nil
}

//     PLAYER_SESSION
// ======================

// PlayerSessionResolver struct
type PlayerSessionResolver struct {
	player    PlayerResolver
	playerID  graphql.ID
	sessionID graphql.ID
	buyIn     int32
	walkout   int32
	db        *postgresDb
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
func (ps *PlayerSessionResolver) BuyIn() int32 {
	return ps.buyIn
}

// WalkOut getter
func (ps *PlayerSessionResolver) WalkOut() int32 {
	return ps.walkout
}

// Player getter
func (ps *PlayerSessionResolver) Player() (*PlayerResolver, error) {
	playerID, err := strconv.Atoi(string(ps.PlayerID()))
	if err != nil {
		return nil, errors.New("PlayerID must be numerical")
	}

	player, err := ps.db.GetPlayerByID(playerID)
	if err != nil {
		return nil, err
	}

	return &PlayerResolver{
		id:      toGQL(player.ID),
		realmID: graphql.ID(strconv.Itoa(player.RealmID)),
		name:    player.Name,
		db:      ps.db,
	}, nil
}
