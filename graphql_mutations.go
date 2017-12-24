package main

import (
	"errors"
	"strconv"
	"time"

	graphql "github.com/neelance/graphql-go"
	null "gopkg.in/guregu/null.v3"
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
	}, r.db}, nil
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

// CreateSessionPlayerSession args struct
type CreateSessionPlayerSession struct {
	PlayerID int32
	Buyin    int32
	Walkout  int32
}

// CreateSession args struct
type CreateSession struct {
	Name           string
	RealmID        int32
	Time           string
	PlayerSessions []*CreateSessionPlayerSession
}

// CreateSession mutation
func (r *Resolver) CreateSession(args CreateSession) (*SessionResolver, error) {
	// TODO: Argument validation

	parsedTime, err := time.Parse(time.RFC3339, args.Time)
	if err != nil {
		return nil, err
	}

	var playerSessions = make([]PlayerSession, len(args.PlayerSessions))
	for i, csps := range args.PlayerSessions {
		playerSessions[i] = PlayerSession{
			PlayerID: int(csps.PlayerID),
			Buyin:    null.IntFrom(int64(csps.Buyin)),
			Walkout:  null.IntFrom(int64(csps.Walkout)),
		}
	}

	session, err := r.db.CreateSession(
		null.IntFrom(int64(args.RealmID)),
		null.StringFrom(args.Name),
		null.TimeFrom(parsedTime),
		playerSessions)

	if err != nil {
		return nil, err
	}
	return &SessionResolver{&GQLSession{
		ID:      graphql.ID(strconv.Itoa(session.ID)),
		Name:    session.Name.Ptr(),
		RealmID: graphql.ID(strconv.Itoa(session.RealmID)),
		Time:    session.Time.Format(time.RFC3339),
	}}, nil
}

// UpdateSessionPlayerSession args struct
type UpdateSessionPlayerSession struct {
	SessionID int32
	PlayerID  int32
	Buyin     int32
	Walkout   int32
}

// UpdateSession args struct
type UpdateSession struct {
	SessionID      int32
	Name           string
	RealmID        int32
	Time           string
	PlayerSessions []*UpdateSessionPlayerSession
}

// UpdateSession mutation
func (r *Resolver) UpdateSession(args UpdateSession) (*SessionResolver, error) {
	// TODO: Argument validation

	parsedTime, err := time.Parse(time.RFC3339, args.Time)
	if err != nil {
		return nil, err
	}

	var playerSessions = make([]PlayerSession, len(args.PlayerSessions))
	for i, usps := range args.PlayerSessions {
		playerSessions[i] = PlayerSession{
			PlayerID:  int(usps.PlayerID),
			SessionID: int(usps.SessionID),
			Buyin:     null.IntFrom(int64(usps.Buyin)),
			Walkout:   null.IntFrom(int64(usps.Walkout)),
		}
	}

	session, err := r.db.UpdateSession(
		null.IntFrom(int64(args.RealmID)),
		null.StringFrom(args.Name),
		null.TimeFrom(parsedTime),
		playerSessions)

	if err != nil {
		return nil, err
	}
	return &SessionResolver{&GQLSession{
		ID:      graphql.ID(strconv.Itoa(session.ID)),
		Name:    session.Name.Ptr(),
		RealmID: graphql.ID(strconv.Itoa(session.RealmID)),
		Time:    session.Time.Format(time.RFC3339),
	}}, nil
}
